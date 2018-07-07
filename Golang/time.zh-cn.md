# time

* [time.Newtimer](#timenewtimer)
* [time.After](#timeafter)
* [time.Afterfunc](#timeafterfunc)
* [time.Ticker](#timeticker)
* []()



## time.Newtimer

初始化一个到期时间据此时的间隔为3小时30分的定时器

t := time.Newtimer(3*time.Hour + 30*time.Minute)

注意，这里的变量t是*time.NewTimer类型的，这个指针类型的方法集合包含两个方法

Rest
用于重置定时器,该方法返回一个bool类型的值

Stop
用来停止定时器,该方法返回一个bool类型的值，如果返回false，说明该定时器在之前已经到期或者已经被停止了,反之返回true。

通过定时器的字段C,我们可以及时得知定时器到期的这个事件来临，C是一个chan time.Time类型的缓冲通道，一旦触及到期时间，定时器就会向自己的C字段发送一个time.Time类型的元素值

示例一：一个简单定时器

```go
package main

import (
    "fmt"
    "time"
)

func main(){
    //初始化定时器
    t := time.NewTimer(2 * time.Second)
    //当前时间
    now := time.Now()
    fmt.Printf("Now time : %v.\n", now)

    expire := <- t.C
    fmt.Printf("Expiration time: %v.\n", expire)
}
```

```
Now time : 2015-10-31 01:19:07.210771347 +0800 CST.
Expiration time: 2015-10-31 01:19:09.215489592 +0800 CST.
```

示例二：我们在改造下之前的那个简单超时操作

```go
package main

import (
    "fmt"
    "time"
)
func main(){
    //初始化通道
    ch11 := make(chan int, 1000)
    sign := make(chan byte, 1)

    //给ch11通道写入数据
    for i := 0; i < 1000; i++ {
        ch11 <- i
    }

    //单独起一个Goroutine执行select
    go func(){
        var e int
        ok := true
        //首先声明一个*time.Timer类型的值，然后在相关case之后声明的匿名函数中尽可能的复用它
        var timer *time.Timer

        for{
            select {
                case e = <- ch11:
                    fmt.Printf("ch11 -> %d\n",e)
                case <- func() <-chan time.Time {
                    if timer == nil{
                        //初始化到期时间据此间隔1ms的定时器
                        timer = time.NewTimer(time.Millisecond)
                    }else {
                        //复用，通过Reset方法重置定时器
                        timer.Reset(time.Millisecond)
                    }
                    //得知定时器到期事件来临时，返回结果
                    return timer.C
                }():
                    fmt.Println("Timeout.")
                    ok = false
                    break
            }
            //终止for循环
            if !ok {
                sign <- 0
                break
            }
        }

    }()

    //惯用手法，读取sign通道数据，为了等待select的Goroutine执行。
    <- sign
}
```

## time.After

* time.After函数， 表示多少时间之后，但是在取出channel内容之前不阻塞，后续程序可以继续执行
* 鉴于After特性，其通常用来处理程序超时问题

```go
package main

import (
    "fmt"
    "time"
)

func main(){
    ch1 := make(chan int, 1)
    ch2 := make(chan int, 1)

    select {
        case e1 := <-ch1:
        //如果ch1通道成功读取数据，则执行该case处理语句
            fmt.Printf("1th case is selected. e1=%v",e1)
        case e2 := <-ch2:
        //如果ch2通道成功读取数据，则执行该case处理语句
            fmt.Printf("2th case is selected. e2=%v",e2)
        case <- time.After(2 * time.Second):
            fmt.Println("Timed out")
    }
}
```

## time.Afterfunc

```go
package main

import (
    "fmt"
    "time"
)
func main(){
    var t *time.Timer

    f := func(){
        fmt.Printf("Expiration time : %v.\n", time.Now())
        fmt.Printf("C`s len: %d\n", len(t.C))
    }

    t = time.AfterFunc(1*time.Second, f)
    //让当前Goroutine 睡眠2s，确保大于内容的完整
    //这样做原因是，time.AfterFunc的调用不会被阻塞。它会以一部的方式在到期事件来临执行我们自定义函数f。
    time.Sleep(2 * time.Second)
}
```

```
Expiration time : 2015-10-31 01:04:42.579988801 +0800 CST.
C`s len: 0
```

第二行打印内容说明：定时器的字段C并没有缓冲任何元素值。这也说明了，在给定了自定义函数后，默认的处理方法(向C发送代表绝对到期时间的元素值)就不会被执行了。


## time.Ticker

结构体类型time.Ticker表示了断续器的静态结构。
就是周期性的传达到期时间的装置。这种装置的行为方式与仅有秒针的钟表有些类似，只不过间隔时间可以不是1s。
初始化一个断续器

### 示例一：使用时间控制停止ticker

```go
package main

import (
    "fmt"
    "time"
)

func main(){
    //初始化断续器,间隔2s
    var ticker *time.Ticker = time.NewTicker(1 * time.Second)

    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()

    time.Sleep(time.Second * 5)   //阻塞，则执行次数为sleep的休眠时间/ticker的时间
    ticker.Stop()     
    fmt.Println("Ticker stopped")
}
```

### 示例二：使用channel控制停止ticker

```go
package main

import (
    "fmt"
    "time"
)

func main(){
    //初始化断续器,间隔2s
    var ticker *time.Ticker = time.NewTicker(100 * time.Millisecond)

    //num为指定的执行次数
    num := 2
    c := make(chan int, num) 
    go func() {
        for t := range ticker.C {
            c <- 1
            fmt.Println("Tick at", t)
        }
    }()

    time.Sleep(time.Millisecond * 1500)
    ticker.Stop()     
    fmt.Println("Ticker stopped")
}
```

参考链接:

1. [Golang time包的定时器/断续器](https://www.kancloud.cn/digest/batu-go/153534)
1. []()
1. 
1. 
1. 
1. 
1. 
1. 
1. 
