# 一致性哈希 + 多 Redis 集群请求分离

基于一致性哈希（Consistent Hashing）将 Redis 请求按 **key** 路由到多个 Redis 集群，实现请求分离与平滑扩缩容。

## 特性

- **请求分离**：不同 key 固定落在不同集群，减轻单集群压力。
- **扩缩容友好**：增删集群时，仅环上相邻区间的 key 会迁移，大部分请求仍命中原集群。
- **虚拟节点**：每个物理集群对应多个虚拟节点，分布更均匀，减少热点。
- **动态集群**：支持 `AddCluster` / `RemoveCluster` 运行时变更。

## 目录结构

```
consistent-hash-redis/
├── consistent_hash.go   # 一致性哈希环
├── redis_router.go      # 按 key 路由到多集群的 Redis 封装
├── demo.go              # 演示：key -> 集群 与可选 Set/Get
├── cmd/main.go          # 可执行示例
└── readme.md
```

## 使用方式

### 1. 仅看路由分布（无需 Redis）

```bash
cd Golang/interview
go run ./consistent-hash-redis/cmd
```

会打印若干 key 被路由到哪个集群（cluster-a / cluster-b / cluster-c）。

### 2. 在代码中创建 Router

```go
import "interview/consistent-hash-redis"

addrs := map[string]string{
    "cluster-a": "127.0.0.1:6379",
    "cluster-b": "127.0.0.1:6380",
    "cluster-c": "127.0.0.1:6381",
}
router, err := chredis.NewRouter(150, addrs)
if err != nil {
    log.Fatal(err)
}
defer func() {
    // 若有需要，可在这里关闭各集群连接
}()

ctx := context.Background()
// 按 key 自动路由到对应集群
router.Set(ctx, "user:1001", "alice", 0)
val, err := router.Get(ctx, "user:1001").Result()
cluster := router.ClusterForKey("user:1001") // 查看该 key 落在哪一集群
```

### 3. 使用已有 Redis 客户端（如 ClusterClient）

```go
clients := map[string]chredis.RedisClient{
    "cluster-a": redis.NewClusterClient(&redis.ClusterOptions{...}),
    "cluster-b": redis.NewClusterClient(&redis.ClusterOptions{...}),
}
router := chredis.NewRouterWithClients(150, clients)
```

### 4. 动态扩缩容

```go
router.AddCluster("cluster-d", redis.NewClient(&redis.Options{Addr: "127.0.0.1:6382"}))
// ...
router.RemoveCluster("cluster-d")
```

## 依赖

- Go 1.25+
- `github.com/redis/go-redis/v9`

在项目根目录执行：

```bash
go get github.com/redis/go-redis/v9
go mod tidy
```

## 注意事项

- **同一 key 始终落在同一集群**，因此跨 key 的事务或多 key 命令需业务层保证落在同一集群，或使用 Hash Tag 等策略将相关 key 收敛到同一集群。
- **Del/Exists** 等多 key 接口当前按第一个 key 所在集群执行，跨集群多 key 需自行拆分调用。
- 虚拟节点数 `replicas` 建议 100～200，越大分布越均匀，内存略增。
