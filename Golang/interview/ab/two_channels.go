package ab

import (
	"fmt"
	"sync"
)

// RunTwoChannels 使用两个无缓冲 channel 实现两个协程交替顺序打印 1..n。
// 协程 A 打印奇数，协程 B 打印偶数，通过互相传递“令牌”严格交替。
func RunTwoChannels(n int) {
	if n < 1 {
		return
	}
	chA := make(chan struct{})
	chB := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= n; i += 2 {
			<-chA
			fmt.Println(i)
			if i+1 <= n {
				chB <- struct{}{}
			}
		}
		if n >= 2 && n%2 == 0 {
			<-chA // 收掉偶数方最后一轮的回传，避免 even 阻塞
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= n; i += 2 {
			<-chB
			fmt.Println(i)
			chA <- struct{}{} // 每轮都回传令牌，让奇数方能继续或退出
		}
	}()

	chA <- struct{}{}
	wg.Wait()
}
