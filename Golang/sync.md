
## 并发模型

一共列举了3种

1. channel
1. WaitGroup
1. context


```
package main

import (
	"context"
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"net/http"
	"sync"
	"time"
)

func main() {
	// useChannel()
	// useWaitGroup()
	useContext()
}

func useChannel() {
	// 通过无缓冲通道来实现多 goroutine 并发控制

	// create channel to synchronize
	done := make(chan bool) // 无缓冲通道
	defer close(done)

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("one done")
		done <- true
	}()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("two done")
		done <- true
	}()

	// wait until both are done
	// 当主 goroutine 运行到 <-done 接受 channel 的值的时候，如果该  channel 中没有数据，就会一直阻塞等待，直到有值。
	for c := 0; c < 2; c++ {
		<-done
	}
	fmt.Println("handle1 done")
}

func useWaitGroup() {
	// 通过sync包中的WaitGroup 实现并发控制

	var wg sync.WaitGroup

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(5 * time.Second)
		fmt.Println("1 done")
		wg.Done()
	}(&wg)

	wg.Add(1)
	go func(wg *sync.WaitGroup) {
		time.Sleep(9 * time.Second)
		fmt.Println("2 done")
		wg.Done()
	}(&wg)
	wg.Wait()
	fmt.Println("handle2 done")

	// 在 sync 包中，提供了 WaitGroup ，它会等待它收集的所有 goroutine 任务全部完成，在主 goroutine 中 Add(delta int) 索要等待goroutine 的数量。在每一个 goroutine 完成后 Done() 表示这一个goroutine 已经完成，当所有的 goroutine 都完成后，在主 goroutine 中 WaitGroup 返回。
}

func useContext() {

	wg := &sync.WaitGroup{}
	values := []string{"https://www.baidu.com/", "https://www.zhihu.com/"}
	ctx, cancel := context.WithCancel(context.Background())

	for _, url := range values {
		wg.Add(1)
		subCtx := context.WithValue(ctx, favContextKey("url"), url)
		go reqURL(subCtx, wg)
	}

	go func() {
		//3秒后取消所有任务
		time.Sleep(time.Second * 3)
		cancel()
	}()

	wg.Wait()
	fmt.Println("exit main goroutine")
}

type favContextKey string

func reqURL(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	url, _ := ctx.Value(favContextKey("url")).(string)
	for {
		select {
		case <-ctx.Done():
			fmt.Printf("stop getting url:%s\n", url)
			return
		default:
			r, err := http.Get(url)
			if r.StatusCode == http.StatusOK && err == nil {
				body, _ := ioutil.ReadAll(r.Body)
				subCtx := context.WithValue(ctx, favContextKey("resp"), fmt.Sprintf("%s%x", url, md5.Sum(body)))
				wg.Add(1)
				go showResp(subCtx, wg)
			}
			r.Body.Close()
			//启动子goroutine是为了不阻塞当前goroutine，这里在实际场景中可以去执行其他逻辑，这里为了方便直接sleep一秒
			// doSometing()
			time.Sleep(time.Second * 1)
		}
	}
}

func showResp(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()
	for {
		select {
		case <-ctx.Done():
			fmt.Println("stop showing resp")
			return
		default:
			//子goroutine里一般会处理一些IO任务，如读写数据库或者rpc调用，这里为了方便直接把数据打印
			fmt.Println("printing ", ctx.Value(favContextKey("resp")))
			time.Sleep(time.Second * 1)
		}
	}
}

```

参考:
1. [Golang并发模型](http://xuchongfeng.github.io/2016/03/24/Golang%E5%B9%B6%E5%8F%91%E6%A8%A1%E5%9E%8B/)
2. [Golang CSP并发模型](https://www.jianshu.com/p/36e246c6153d)
1. [Golang 并发控制的两种模式](https://www.golangnote.com/topic/184.html)
1. [深入golang之---goroutine并发控制与通信](https://juejin.im/entry/5b32f89151882574d3249aba)

