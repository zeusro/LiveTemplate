package main

import (
	"fmt"
	"time"
)

// 错误示例：Channel 泄漏 - 发送者阻塞
// 问题：发送者向无缓冲 channel 发送数据，但没有接收者，导致 goroutine 永远阻塞

func senderChannelLeak() {
	ch := make(chan int) // 无缓冲 channel
	
	go func() {
		// 发送数据，但没有接收者
		ch <- 42
		fmt.Println("This will never be printed")
	}()
	
	// 主 goroutine 退出，发送者 goroutine 永远阻塞
	// 导致内存泄漏
}

// 错误示例：Channel 泄漏 - 接收者阻塞
func receiverChannelLeak() {
	ch := make(chan int)
	
	go func() {
		// 等待接收数据，但永远不会收到
		val := <-ch
		fmt.Println("Received:", val)
	}()
	
	// 没有发送者，接收者永远阻塞
}

// 错误示例：Channel 泄漏 - 多个发送者
func multipleSenderLeak() {
	ch := make(chan int)
	
	// 启动多个发送者
	for i := 0; i < 100; i++ {
		go func(id int) {
			ch <- id // 所有发送者都会阻塞
		}(i)
	}
	
	// 没有接收者，所有发送者 goroutine 都阻塞
	// 导致严重的内存泄漏
}

// 错误示例：Channel 泄漏 - 忘记关闭 channel
func unclosedChannelLeak() {
	ch := make(chan int)
	
	go func() {
		for i := 0; i < 10; i++ {
			ch <- i
		}
		// 忘记关闭 channel
		// close(ch) // 应该在这里关闭
	}()
	
	go func() {
		// 使用 range 等待 channel 关闭
		for val := range ch {
			fmt.Println("Received:", val)
		}
		// 由于 channel 未关闭，这个 goroutine 永远阻塞
	}()
}

// 正确示例：使用 context 或 done channel 控制
func correctChannelExample() {
	ch := make(chan int)
	done := make(chan bool)
	
	// 发送者
	go func() {
		defer close(ch)
		for i := 0; i < 10; i++ {
			select {
			case ch <- i:
			case <-done:
				return
			}
		}
	}()
	
	// 接收者
	go func() {
		for val := range ch {
			fmt.Println("Received:", val)
		}
	}()
	
	time.Sleep(1 * time.Second)
	close(done) // 通知发送者退出
}

// 正确示例：使用缓冲 channel 或超时
func correctChannelWithTimeout() {
	ch := make(chan int, 10) // 使用缓冲 channel
	
	go func() {
		for i := 0; i < 100; i++ {
			select {
			case ch <- i:
			case <-time.After(1 * time.Second):
				// 超时处理，避免永久阻塞
				fmt.Println("Send timeout")
				return
			}
		}
		close(ch)
	}()
	
	go func() {
		for val := range ch {
			fmt.Println("Received:", val)
		}
	}()
	
	time.Sleep(2 * time.Second)
}

func main() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// senderChannelLeak()
	// receiverChannelLeak()
	// multipleSenderLeak()
	// unclosedChannelLeak()
	
	// 正确示例
	correctChannelExample()
	correctChannelWithTimeout()
}
