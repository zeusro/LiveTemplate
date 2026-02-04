package ab

import (
	"fmt"
	"sync"
)

// RunCond 使用 sync.Mutex + sync.Cond 实现两个协程交替顺序打印 1..n。
// 共享 turn：0 表示 A 的回合，1 表示 B 的回合；不是自己回合就 Wait，打完改 turn 并 Signal。
func RunCond(n int) {
	if n < 1 {
		return
	}
	var mu sync.Mutex
	cond := sync.NewCond(&mu)
	turn := 0 // 0 = A(奇数), 1 = B(偶数)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= n; i += 2 {
			mu.Lock()
			for turn != 0 {
				cond.Wait()
			}
			fmt.Println(i)
			turn = 1
			cond.Signal()
			mu.Unlock()
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= n; i += 2 {
			mu.Lock()
			for turn != 1 {
				cond.Wait()
			}
			fmt.Println(i)
			turn = 0
			cond.Signal()
			mu.Unlock()
		}
	}()

	mu.Lock()
	turn = 0
	cond.Broadcast()
	mu.Unlock()

	wg.Wait()
}
