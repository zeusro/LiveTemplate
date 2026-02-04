// 可执行示例：一致性哈希将请求路由到多个 Redis 集群，实现请求分离。
package main

import (
	"fmt"
	chredis "interview/consistent-hash-redis"
)

func main() {
	fmt.Println("一致性哈希 + 多 Redis 集群请求分离示例")
	// 传 nil 时仅打印 key -> 集群 映射，不依赖 Redis 进程
	chredis.Demo(nil)
	fmt.Println()
	// 若本机有 6379/6380/6381 三个 Redis，可改为：
	// chredis.Demo(map[string]string{
	//     "cluster-a": "127.0.0.1:6379",
	//     "cluster-b": "127.0.0.1:6380",
	//     "cluster-c": "127.0.0.1:6381",
	// })
}
