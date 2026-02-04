package chredis

import (
	"context"
	"fmt"
	"time"
)

// Demo 演示一致性哈希路由：打印若干 key 被路由到的集群，并可选执行真实 Redis 请求。
// 若不传 addrs 或传 nil，仅打印 key -> 集群 的映射（无需启动 Redis）。
func Demo(addrs map[string]string) {
	namesOnly := len(addrs) == 0
	if namesOnly {
		// 仅演示路由分布，不连接 Redis
		addrs = map[string]string{
			"cluster-a": "127.0.0.1:6379",
			"cluster-b": "127.0.0.1:6380",
			"cluster-c": "127.0.0.1:6381",
		}
	}
	ring := New(150, nil)
	for name := range addrs {
		ring.Add(name)
	}

	keys := []string{"user:1001", "user:1002", "order:2001", "order:2002", "session:abc", "cache:foo", "cache:bar"}
	fmt.Println("=== 一致性哈希：Key -> 集群 ===")
	for _, key := range keys {
		cluster := ring.Get(key)
		fmt.Printf("  %-15s -> %s\n", key, cluster)
	}

	if namesOnly {
		return
	}
	router, err := NewRouter(150, addrs)
	if err != nil {
		fmt.Printf("创建 Router 失败（可忽略，仅需 Redis 未启动）: %v\n", err)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	pings := router.Ping(ctx)
	allOk := true
	for name, err := range pings {
		if err != nil {
			fmt.Printf("  Ping %s: %v\n", name, err)
			allOk = false
		}
	}
	if !allOk {
		fmt.Println("部分 Redis 未就绪，跳过 Set/Get 演示。请启动多个 Redis 实例（如 6379/6380/6381）后再运行。")
		return
	}

	fmt.Println("\n=== 实际 Set/Get 演示 ===")
	for _, key := range keys[:3] {
		val := "value-" + key
		if err := router.Set(ctx, key, val, 0).Err(); err != nil {
			fmt.Printf("  Set %s: %v\n", key, err)
			continue
		}
		got, err := router.Get(ctx, key).Result()
		if err != nil {
			fmt.Printf("  Get %s: %v\n", key, err)
			continue
		}
		fmt.Printf("  %s -> 集群 %s, Get = %s\n", key, router.ClusterForKey(key), got)
	}
}
