package ab

import (
	"fmt"
	"sync"
)

// msg 携带要打印的数以及轮到谁（odd=true 表示奇数方处理）
type msg struct {
	v   int
	odd bool
}

// RunSingleChannel 使用单一 channel 实现两个协程交替顺序打印 1..n。
// channel 传 (数值, 轮到谁)；每个协程只处理自己的消息，否则放回 channel，保证严格交替。
func RunSingleChannel(n int) {
	if n < 1 {
		return
	}
	ch := make(chan msg, 1)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 1; i <= n; i += 2 {
			for {
				m := <-ch
				if !m.odd {
					ch <- m
					continue
				}
				fmt.Println(m.v)
				if i+1 <= n {
					ch <- msg{v: i + 1, odd: false}
				}
				break
			}
		}
	}()

	go func() {
		defer wg.Done()
		for i := 2; i <= n; i += 2 {
			for {
				m := <-ch
				if m.odd {
					ch <- m
					continue
				}
				fmt.Println(m.v)
				if i+1 <= n {
					ch <- msg{v: i + 1, odd: true}
				}
				break
			}
		}
	}()

	ch <- msg{v: 1, odd: true}
	wg.Wait()
}
