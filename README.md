# WeVote

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

## 项目架构

```
.
├── cmd                  // 项目执行入口，包括 server 和 cronjob
│   ├── cronjob
│   │   └── cronjob.go
│   └── server.go
├── configs              // 项目配置文件
├── internal
│   ├── common           // 通用组件，定义了规范化的错误处理与日志处理
│   │   ├── errors
│   │   └── log
│   ├── config           // 配置读取
│   ├── models           // model 层
│   │   ├── bitset.go    // 使用 Redis 实现布隆过滤器底层的 BitSet 结构
│   │   ├── conn.go      // 初始化数据库连接
│   │   ├── ticket.go    // 使用 Redis Cluster 同步 Global Ticket 与 Usage Counter
│   │   └── vote.go      // 投票数据使用 Redis 进行缓存，使用 Write-behind 策略异步回写数据库
│   ├── pkg
│   │   ├── cache
│   │   │   └── write_behind.go  // Write-Behind data model 的抽象实现
│   │   └── ticket
│   │       └── ticket.go        // 签发与检验 Ticket
│   ├── schema                   // GraphQL Schema 层，接收 GraphQL 请求并依次调用 Service 层接口
│   └── services
│       ├── captcha.go   // ReCAPTCHA 人机验证
│       ├── filter.go    // 布隆过滤器
│       ├── init.go
│       ├── ticket.go    // Ticket Service，获取 Ticket & 消耗 Ticket 次数并在超过限制时返回错误
│       └── vote.go      // Vote Service，投票 & 获取投票数据
├── pkg
│   ├── bloomfilter      // 布隆过滤器抽象实现
│   ├── captcha          // ReCAPTCHA Client 操作
│   └── cronjob          // 定时任务
├── README.md
└── scripts
```



![image-20240406231550938](https://cdn.just-plain.fun/img/image-20240406231550938.png)



## 设计

### Ticket

使用无状态 JWT 实现，有效期为 2s

![image-20240406224802752](https://cdn.just-plain.fun/img/image-20240406224802752.png)

JWT Payload 包含签发时间(`iat`)、过期时间(`exp`)、唯一标识符 TicketID(`sub`)，使用 HS256 加密签名来确保 Ticket 的合法性

在 Redis 中记录 TID 对应的投票次数，TTL 比 Ticket 稍长，设为 3s

### 投票

投票数据使用 Redis 进行缓存，使用 Write-behind 策略异步回写数据库

使用 reCAPTCHA v3 对接口进行人机验证，防御恶意刷票

使用布隆过滤器对投票请求参数进行过滤，避免缓存穿透攻击

### 性能

- 在 Write-Behind 过程中，使用乐观锁来确保数据的一致性

- 对于全局数据，使用互斥锁来确保并发安全

- 每分钟进行一次缓存 Write Back，同时使用高可用的 Redis 分布式集群避免宕机，并利用 Redis RDB 和 AOF 混合机制进行持久化
- 多实例情况下能够保证数据一致性，使用 Redis Cluster 并支持数据库水平分表，具有良好的水平拓展性

