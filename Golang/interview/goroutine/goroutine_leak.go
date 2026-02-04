package goroutine

import (
	"fmt"
	"time"
)

// 错误示例：Goroutine 泄漏
// 问题：启动的 goroutine 永远不会退出，导致内存泄漏

func goroutineLeakExample() {
	// 启动一个 goroutine，但没有提供退出机制
	go func() {
		for {
			// 无限循环，goroutine 永远不会退出
			fmt.Println("Running...")
			time.Sleep(1 * time.Second)
		}
	}()
	
	// 主函数返回后，这个 goroutine 仍然在运行
	// 导致内存泄漏
}

// 错误示例：多个 goroutine 泄漏
func multipleGoroutineLeak() {
	for i := 0; i < 1000; i++ {
		go func(id int) {
			// 没有退出条件的 goroutine
			for {
				time.Sleep(1 * time.Second)
				// 做一些工作，但永远不会退出
			}
		}(i)
	}
	// 1000 个 goroutine 永远不会退出，造成严重的内存泄漏
}

// 错误示例：goroutine 等待永远不会到来的信号
func waitingGoroutineLeak() {
	ch := make(chan bool)
	
	go func() {
		// 等待一个永远不会发送的信号
		<-ch
		fmt.Println("This will never be printed")
	}()
	
	// ch 永远不会被关闭或发送数据
	// goroutine 永远阻塞，导致泄漏
}

// 正确示例：使用 context 控制 goroutine 退出
func correctGoroutineExample() {
	done := make(chan bool)
	
	go func() {
		defer close(done)
		for {
			select {
			case <-done:
				return // 正确退出
			default:
				fmt.Println("Working...")
				time.Sleep(1 * time.Second)
			}
		}
	}()
	
	// 在适当的时候关闭 done channel
	time.Sleep(5 * time.Second)
	close(done)
}

// RunGoroutineLeakExamples 演示本文件中的正确用法示例。
func RunGoroutineLeakExamples() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// goroutineLeakExample()
	// multipleGoroutineLeak()
	// waitingGoroutineLeak()
	
	// 正确示例
	correctGoroutineExample()
}
