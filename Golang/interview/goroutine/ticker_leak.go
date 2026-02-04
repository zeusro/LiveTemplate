package goroutine

import (
	"fmt"
	"time"
)

// 错误示例：Ticker 泄漏
// 问题：time.Ticker 创建后没有调用 Stop()，导致底层 goroutine 和资源无法释放

func tickerLeakExample() {
	// 创建 ticker，每秒触发一次
	ticker := time.NewTicker(1 * time.Second)
	
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	
	// 忘记调用 ticker.Stop()
	// ticker 会一直运行，导致内存泄漏
	// 即使 ticker 不再使用，底层的 goroutine 和 channel 也不会被释放
}

// 错误示例：多个 Ticker 泄漏
func multipleTickerLeak() {
	for i := 0; i < 100; i++ {
		ticker := time.NewTicker(1 * time.Second)
		
		go func(id int, t *time.Ticker) {
			for range t.C {
				fmt.Printf("Ticker %d ticked\n", id)
			}
		}(i, ticker)
		
		// 每个 ticker 都没有被停止
		// 导致 100 个 ticker 和对应的 goroutine 泄漏
	}
}

// 错误示例：Timer 泄漏（类似问题）
func timerLeakExample() {
	// 创建 timer
	timer := time.NewTimer(5 * time.Second)
	
	go func() {
		<-timer.C
		fmt.Println("Timer fired")
	}()
	
	// 如果 timer 被取消或不再需要，但没有调用 Stop()
	// 可能导致资源泄漏（虽然 Timer 比 Ticker 好一些）
}

// 正确示例：使用 defer 确保 Ticker 被停止
func correctTickerExample() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop() // 确保 ticker 被停止
	
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	
	// 运行一段时间后退出
	time.Sleep(5 * time.Second)
	// defer 会确保 ticker.Stop() 被调用
}

// 正确示例：使用 context 控制 Ticker 生命周期
func correctTickerWithContext() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	done := make(chan bool)
	
	go func() {
		for {
			select {
			case t := <-ticker.C:
				fmt.Println("Tick at", t)
			case <-done:
				return
			}
		}
	}()
	
	time.Sleep(5 * time.Second)
	close(done) // 通知 goroutine 退出
}

// 正确示例：使用 time.After 代替 Ticker（如果只需要一次）
func correctTimerExample() {
	// 如果只需要延迟执行，使用 time.After 更简单
	select {
	case <-time.After(5 * time.Second):
		fmt.Println("5 seconds passed")
	}
	// time.After 会自动清理，不需要手动停止
}

// RunTickerLeakExamples 演示本文件中的正确用法示例。
func RunTickerLeakExamples() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// tickerLeakExample()
	// multipleTickerLeak()
	// timerLeakExample()
	
	// 正确示例
	correctTickerExample()
	correctTickerWithContext()
	correctTimerExample()
}
