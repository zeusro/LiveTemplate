// Part 2 审计题的修复实现
// 修复：并发安全 map、timer 正确释放、goroutine 内不引用 request、优雅停机与超时配置

package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	stats   = make(map[string]int)
	statsMu sync.Mutex
)

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/track", func(w http.ResponseWriter, r *http.Request) {
		deviceID := r.URL.Query().Get("id")
		// 在启动 goroutine 前拷贝所需数据，避免在 goroutine 内引用 r
		ctx := r.Context()

		go func(id string, ctx context.Context) {
			// 修复 1：写 map 时加锁
			statsMu.Lock()
			stats[id]++
			count := stats[id]
			statsMu.Unlock()

			// 修复 2：使用 NewTimer + Stop，避免 time.After 泄漏
			timer := time.NewTimer(2 * time.Second)
			defer timer.Stop()

			select {
			case <-timer.C:
				statsMu.Lock()
				count = stats[id]
				statsMu.Unlock()
				fmt.Printf("Device %s processed. Total: %d\n", id, count)
			case <-ctx.Done():
				fmt.Println("Client disconnected")
			}
		}(deviceID, ctx)

		w.Write([]byte("Tracking started"))
	})

	// 修复 4：生产级超时与优雅停机
	server := &http.Server{
		Addr:              ":8080",
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("server error: %v\n", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		fmt.Printf("shutdown error: %v\n", err)
	}
	fmt.Println("server stopped")
}
