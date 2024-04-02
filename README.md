# WeVote

## 技术栈

- https://github.com/graphql-go/graphql

## 接口

```graphql
type User {
  username: String!
  votes: Int!
}

type Query {
  getUserVotes(username: String!): User
  getTicket: String
}

type Mutation {
  vote(users: [String!]!, ticket: String!): String
}
```

## 设计

### Ticket

使用无状态JWT实现，有效期为2s

JWT Claims 需包含签发时间、过期时间、唯一标识符 TID

Redis 记录 TID 对应的投票次数，TTL为2s

### 投票

投票数据使用 Redis 进行缓存，使用 Write-behind 策略异步回写数据库

使用 reCAPTCHA v3 对接口进行人机验证，防御恶意刷票

### 避免缓存穿透

- 加个布隆过滤器

## 条件竞争问题

投票过程即 Redis 执行 INCR，需要考虑数据竞争情况

Redis 集群本身能够通过自身机制来保证原子性（分片、读写分离）

## 数据缓存一致性

如何保证数据一致性？

- 每分钟写入一次数据库做持久化存储
- 使用 Redis RDB和AOF混合持久化机制来保证服务宕机或重启时数据的一致性

## 水平扩展性

1. 多节点负载均衡
2. Redis 分布式集群
3. 数据库表水平拆分

## 性能结果