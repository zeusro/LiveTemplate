// 题目中的原始代码（存在多处隐患，仅用于审计对照）
// 编译可通过，但在生产环境中存在严重隐患。
// 参见 recruitment-evaluation-solutions.md Part 2 的详细说明。

package main

import (
	"fmt"
	"net/http"
	"time"
)

// 全局统计缓存
// 隐患 1：多 goroutine 并发读写，无同步 → concurrent map writes panic
var stats = make(map[string]int)

func main() {
	http.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("id")

		// 模拟异步处理埋点逻辑
		go func() {
			// [隐患点 1] 并发写 map，会 panic
			stats[deviceID]++

			// 模拟耗时较长的分析操作
			// [隐患点 2] time.After 在 select 中：若 ctx.Done() 先触发，timer 不会释放 → 高并发下 timer 泄漏
			// [隐患点 3] goroutine 内使用 r：handler 返回后 r 可能被复用，只读 r.Context() 尚可，但不宜再引用 r
			ctx := r.Context()
			select {
			case <-time.After(2 * time.Second):
				// 这里再次读 map，与其它 goroutine 的写并发
				fmt.Printf("Device %s processed. Total: %d\n", deviceID, stats[deviceID])
			case <-ctx.Done():
				fmt.Println("Client disconnected")
			}
		}()

		w.Write([]byte("Tracking started"))
	})

	// [隐患点 4] 无优雅停机、无 ReadTimeout/WriteTimeout 等生产级配置
	http.ListenAndServe(":8080", nil)
}
