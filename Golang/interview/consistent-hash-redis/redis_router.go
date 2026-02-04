package chredis

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 抽象 Redis 操作，便于测试或替换为 ClusterClient。
type RedisClient interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	Ping(ctx context.Context) *redis.StatusCmd
}

// Router 根据 key 通过一致性哈希将请求路由到多个 Redis 集群之一。
type Router struct {
	ring    *Ring
	clients map[string]RedisClient
}

// NewRouter 创建路由器。replicas 为每个集群的虚拟节点数，建议 100～200。
// addrs 为「节点名 -> Redis 地址」映射，节点名用于环上标识，如 "cluster-a", "cluster-b"。
func NewRouter(replicas int, addrs map[string]string) (*Router, error) {
	ring := New(replicas, nil)
	clients := make(map[string]RedisClient, len(addrs))
	for name, addr := range addrs {
		clients[name] = redis.NewClient(&redis.Options{Addr: addr})
		ring.Add(name)
	}
	return &Router{ring: ring, clients: clients}, nil
}

// NewRouterWithClients 使用已有 Redis 客户端构造路由器（便于接入 ClusterClient 或自定义实现）。
func NewRouterWithClients(replicas int, clients map[string]RedisClient) *Router {
	ring := New(replicas, nil)
	for name := range clients {
		ring.Add(name)
	}
	return &Router{ring: ring, clients: clients}
}

// clientByKey 根据 key 返回对应的 Redis 客户端。
func (r *Router) clientByKey(key string) RedisClient {
	name := r.ring.Get(key)
	return r.clients[name]
}

// Get 根据 key 路由到对应集群并执行 Get。
func (r *Router) Get(ctx context.Context, key string) *redis.StringCmd {
	c := r.clientByKey(key)
	if c == nil {
		return redis.NewStringResult("", redis.Nil)
	}
	return c.Get(ctx, key)
}

// Set 根据 key 路由到对应集群并执行 Set。
func (r *Router) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd {
	c := r.clientByKey(key)
	if c == nil {
		return redis.NewStatusResult("", redis.ErrClosed)
	}
	return c.Set(ctx, key, value, expiration)
}

// Del 根据 key 路由到对应集群并执行 Del。
func (r *Router) Del(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) == 0 {
		return redis.NewIntResult(0, nil)
	}
	c := r.clientByKey(keys[0])
	if c == nil {
		return redis.NewIntResult(0, redis.ErrClosed)
	}
	return c.Del(ctx, keys...)
}

// Exists 根据 key 路由到对应集群并执行 Exists。
func (r *Router) Exists(ctx context.Context, keys ...string) *redis.IntCmd {
	if len(keys) == 0 {
		return redis.NewIntResult(0, nil)
	}
	c := r.clientByKey(keys[0])
	if c == nil {
		return redis.NewIntResult(0, redis.ErrClosed)
	}
	return c.Exists(ctx, keys...)
}

// ClusterForKey 返回 key 所路由到的集群名，便于监控或调试。
func (r *Router) ClusterForKey(key string) string {
	return r.ring.Get(key)
}

// Clusters 返回当前所有集群名。
func (r *Router) Clusters() []string {
	return r.ring.Nodes()
}

// AddCluster 动态添加新集群（扩缩容）。
func (r *Router) AddCluster(name string, client RedisClient) {
	r.ring.Add(name)
	r.clients[name] = client
}

// RemoveCluster 从环上移除集群（client 需由调用方自行关闭）。
func (r *Router) RemoveCluster(name string) {
	r.ring.Remove(name)
	delete(r.clients, name)
}

// Ping 对所有集群执行 Ping，返回 集群名 -> error。
func (r *Router) Ping(ctx context.Context) map[string]error {
	out := make(map[string]error)
	for name, c := range r.clients {
		out[name] = c.Ping(ctx).Err()
	}
	return out
}
