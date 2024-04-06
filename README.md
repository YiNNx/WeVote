# WeVote

## 技术栈

- 使用 gqlgen 作为 GraphQL 框架，数据库使用 Redis 集群与 PostgreSQL
- 使用 JWT 进行 Ticket 签发
- 使用 Write-Behind 缓存策略对投票数据进行缓存处理
- 使用布隆过滤器防御缓存穿透攻击，并可选择开启第三方人机验证服务防御恶意刷票

## 接口

```graphql
type Query {
  getUserVotes(username: String!): Int
  getTicket: String!
}

type Mutation {
  vote(users: [String!]!, ticket: String!,recaptchaToken: String): Boolean!
}
```

## 架构

代码架构

```
.
├── cmd                  // 项目执行入口
│   ├── cronjob
│   └── server.go
├── configs              // 配置文件
├── internal
│   ├── common           // 通用组件，定义了规范化的错误返回信息与全局 Logger
│   ├── config           // 配置读取
│   ├── models           // model 层
│   │   ├── bitset.go    // 使用 Redis 实现布隆过滤器底层的 BitSet 结构
│   │   ├── conn.go      // 初始化数据库连接
│   │   ├── ticket.go    // Redis Cluster 同步 Global Ticket 与 Usage Counter
│   │   └── vote.go      // 使用 Write-behind 策略异步回写投票数据
│   ├── pkg
│   │   ├── cache
│   │   │   └── write_behind.go  // Write-Behind data model 的抽象实现
│   │   └── ticket
│   │       └── ticket.go        // 签发与检验 Ticket
│   ├── schema                   // GraphQL Schema 层，接收 GraphQL 请求并依次调用 Service 层接口
│   └── services         // services 业务逻辑层
│       ├── captcha.go   // ReCAPTCHA 人机验证
│       ├── filter.go    // 布隆过滤器
│       ├── ticket.go    // Ticket Service，获取 Ticket & 消耗 Ticket 次数并在超过限制时返回错误
│       └── vote.go      // Vote Service，投票 & 获取投票数据
├── pkg
│   ├── bloomfilter      // 布隆过滤器抽象实现
│   ├── captcha          // ReCAPTCHA Client 操作
│   └── cronjob          // 定时任务
├── README.md
└── scripts
```

项目架构

![image-20240407023420330](https://cdn.just-plain.fun/img/image-20240407023420330.png)

横向拓展服务架构

![image-20240407023643588](https://cdn.just-plain.fun/img/image-20240407023643588.png)

## 设计

### Ticket

使用无状态 JWT 实现，有效期为 2s

![image-20240406224802752](https://cdn.just-plain.fun/img/image-20240406224802752.png)

JWT Payload 包含签发时间 (`iat`)、过期时间 (`exp`)、唯一标识符 TicketID (`sub`)，使用 HS256 加密签名用于检验 Ticket 的合法性

通过 CronJob 每 2 s 签发一次全局的 Ticket，并存储在 Redis 中实现多实例共享，TTL 为 2s。同时记录 TID 对应的投票次数，TTL 比 Ticket 稍长，设为 3s。

消耗 Ticket 使用次数使用 Redis INCR 实现。在多实例情况下，由于 Redis 命令本身具有原子性且串行执行，能够避免多实例可能存在的数据竞争问题，无需额外使用分布式锁。

### 投票

当人机验证开关打开时，会使用 reCAPTCHA v3 对接口进行人机验证，防御恶意刷票

同时使用布隆过滤器对投票的用户名列表进行检验，避免缓存穿透攻击

#### 投票逻辑：

1. 进行人机验证，通过布隆过滤器检验投票用户列表参数
2. 检验 Ticket 有效性，更新 Ticket 使用次数
3. 更新 Vote 数据

#### 性能优化：

投票数据为读写均非常频繁的数据，因此使用 Write-Behind 策略进行缓存处理，即同步更新缓存，异步回写数据库，从而大幅提升程序性能。

获取投票数据：

1. 首先在缓存中获取，若命中缓存则直接返回

2. 未命中缓存，查询数据库

3. 若存在数据，则使用乐观锁更新缓存：

   若更新时缓存不存在则直接写入，若缓存此时已存在则放弃写入，直接查询并返回

进行投票：

1. 首先执行上面的命令逻辑来获取或同步缓存
2. 对缓存数据进行原子性的 IncrBy 操作，在操作前将该数据标记为 dirty 状态

异步回写数据库：

使用 CronJob 异步将缓存回写数据库，每 2s 将缓存中标记为 dirty 的数据更新至数据库。

同时利用 Redis RDB 和 AOF 混合机制进行 Redis 数据持久化，避免 Redis 集群宕机导致 dirty 数据丢失

