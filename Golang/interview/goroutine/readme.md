# Go 语言内存泄漏示例

本文档介绍 Go 语言中常见的几种内存泄漏场景，每种场景都包含错误示例和正确示例，帮助开发者识别和避免内存泄漏问题。

## 目录

1. [Goroutine 泄漏](#1-goroutine-泄漏)
2. [Channel 泄漏](#2-channel-泄漏)
3. [定时器泄漏](#3-定时器泄漏)
4. [全局变量引用泄漏](#4-全局变量引用泄漏)
5. [切片泄漏](#5-切片泄漏)
6. [接口泄漏](#6-接口泄漏)

---

## 1. Goroutine 泄漏

### 问题描述

Goroutine 泄漏是最常见的内存泄漏类型之一。当 goroutine 启动后无法正常退出时，会导致：
- Goroutine 占用的栈内存无法释放（每个 goroutine 默认 2KB，可增长到 1GB）
- Goroutine 中引用的对象无法被 GC 回收
- 系统资源耗尽

### 错误示例

#### 示例 1：无限循环的 Goroutine

```go
func goroutineLeakExample() {
    go func() {
        for {
            // 无限循环，goroutine 永远不会退出
            fmt.Println("Running...")
            time.Sleep(1 * time.Second)
        }
    }()
    // 主函数返回后，这个 goroutine 仍然在运行
}
```

**问题**：Goroutine 中的无限循环没有退出条件，导致 goroutine 永远运行。

#### 示例 2：等待永远不会到来的信号

```go
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
```

**问题**：Goroutine 在 channel 上阻塞等待，但 channel 永远不会收到数据或关闭。

### 正确做法

#### 使用 context 或 done channel 控制退出

```go
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
```

#### 使用 context.Context（推荐）

```go
func correctWithContext(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            default:
                // 执行工作
            }
        }
    }()
}
```

### 检测方法

- 使用 `runtime.NumGoroutine()` 监控 goroutine 数量
- 使用 `pprof` 工具分析 goroutine 泄漏
- 使用 `go tool pprof http://localhost:6060/debug/pprof/goroutine`

---

## 2. Channel 泄漏

### 问题描述

Channel 泄漏通常发生在：
- 发送者向无缓冲 channel 发送数据，但没有接收者
- 接收者等待接收数据，但没有发送者
- Channel 未正确关闭，导致使用 `range` 的 goroutine 永远阻塞

### 错误示例

#### 示例 1：发送者阻塞

```go
func senderChannelLeak() {
    ch := make(chan int) // 无缓冲 channel
    
    go func() {
        // 发送数据，但没有接收者
        ch <- 42
        fmt.Println("This will never be printed")
    }()
    
    // 主 goroutine 退出，发送者 goroutine 永远阻塞
}
```

**问题**：无缓冲 channel 需要发送者和接收者同时准备好。如果没有接收者，发送者会永远阻塞。

#### 示例 2：忘记关闭 Channel

```go
func unclosedChannelLeak() {
    ch := make(chan int)
    
    go func() {
        for i := 0; i < 10; i++ {
            ch <- i
        }
        // 忘记关闭 channel
    }()
    
    go func() {
        // 使用 range 等待 channel 关闭
        for val := range ch {
            fmt.Println("Received:", val)
        }
        // 由于 channel 未关闭，这个 goroutine 永远阻塞
    }()
}
```

**问题**：使用 `range` 遍历 channel 时，会一直等待直到 channel 关闭。如果 channel 未关闭，goroutine 会永远阻塞。

### 正确做法

#### 使用 done channel 或 context

```go
func correctChannelExample() {
    ch := make(chan int)
    done := make(chan bool)
    
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
    
    go func() {
        for val := range ch {
            fmt.Println("Received:", val)
        }
    }()
    
    time.Sleep(1 * time.Second)
    close(done)
}
```

#### 使用超时机制

```go
func correctChannelWithTimeout() {
    ch := make(chan int, 10) // 使用缓冲 channel
    
    go func() {
        for i := 0; i < 100; i++ {
            select {
            case ch <- i:
            case <-time.After(1 * time.Second):
                // 超时处理，避免永久阻塞
                return
            }
        }
        close(ch)
    }()
}
```

### 最佳实践

- 发送者负责关闭 channel（遵循"谁创建谁关闭"的原则）
- 使用缓冲 channel 减少阻塞
- 使用 `select` 配合超时或 done channel 避免永久阻塞
- 明确 channel 的生命周期管理

---

## 3. 定时器泄漏

### 问题描述

`time.Ticker` 和 `time.Timer` 创建后，如果不调用 `Stop()` 方法，会导致：
- 底层的 goroutine 无法退出
- 相关的 channel 无法被 GC 回收
- 资源持续占用

### 错误示例

#### 示例 1：忘记停止 Ticker

```go
func tickerLeakExample() {
    ticker := time.NewTicker(1 * time.Second)
    
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()
    
    // 忘记调用 ticker.Stop()
    // ticker 会一直运行，导致内存泄漏
}
```

**问题**：`time.Ticker` 创建后会启动一个后台 goroutine，如果不调用 `Stop()`，这个 goroutine 永远不会退出。

#### 示例 2：多个 Ticker 泄漏

```go
func multipleTickerLeak() {
    for i := 0; i < 100; i++ {
        ticker := time.NewTicker(1 * time.Second)
        
        go func(id int, t *time.Ticker) {
            for range t.C {
                fmt.Printf("Ticker %d ticked\n", id)
            }
        }(i, ticker)
        
        // 每个 ticker 都没有被停止
    }
}
```

**问题**：创建了多个 ticker 但都没有停止，导致大量 goroutine 泄漏。

### 正确做法

#### 使用 defer 确保停止

```go
func correctTickerExample() {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop() // 确保 ticker 被停止
    
    go func() {
        for t := range ticker.C {
            fmt.Println("Tick at", t)
        }
    }()
    
    time.Sleep(5 * time.Second)
    // defer 会确保 ticker.Stop() 被调用
}
```

#### 使用 context 控制生命周期

```go
func correctTickerWithContext(ctx context.Context) {
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()
    
    go func() {
        for {
            select {
            case t := <-ticker.C:
                fmt.Println("Tick at", t)
            case <-ctx.Done():
                return
            }
        }
    }()
}
```

### 注意事项

- **Ticker vs Timer**：`time.Timer` 只触发一次，`time.Ticker` 会重复触发
- **time.After**：如果只需要延迟执行，使用 `time.After` 更简单，它会自动清理
- **资源管理**：始终使用 `defer ticker.Stop()` 或 `defer timer.Stop()`

---

## 4. 全局变量引用泄漏

### 问题描述

全局变量（包括包级变量）持有对象的引用，导致这些对象无法被 GC 回收。常见场景：
- 全局 map 或 slice 不断增长，从不清理
- 全局 cache 持有大量数据
- 全局 channel 中堆积未消费的数据

### 错误示例

#### 示例 1：全局 Map 不断增长

```go
var globalCache = make(map[string][]byte)

func globalMapLeak() {
    for i := 0; i < 10000; i++ {
        key := fmt.Sprintf("key-%d", i)
        data := make([]byte, 1024*1024) // 1MB
        
        globalCache[key] = data
        // 即使不再使用这些数据，由于全局变量持有引用
        // GC 无法回收这些内存
    }
}
```

**问题**：全局 map 不断添加数据，但从不删除，导致内存持续增长。

#### 示例 2：全局 Slice 持有引用

```go
var globalSlice []*LargeObject

type LargeObject struct {
    Data [1024 * 1024]byte // 1MB
    ID   int
}

func globalSliceLeak() {
    for i := 0; i < 1000; i++ {
        obj := &LargeObject{ID: i}
        globalSlice = append(globalSlice, obj)
        // 即使 obj 不再使用，globalSlice 持有引用
        // 导致所有对象无法被 GC 回收
    }
}
```

### 正确做法

#### 定期清理或限制大小

```go
func correctLimitedCache() {
    const maxSize = 1000
    cache := make(map[string][]byte)
    
    // 限制缓存大小
    if len(cache) >= maxSize {
        // 删除最旧的条目（LRU 策略）
        for k := range cache {
            delete(cache, k)
            break
        }
    }
}
```

#### 使用带过期时间的缓存

```go
type cacheEntry struct {
    data      []byte
    timestamp int64
}

func correctCacheWithTTL() {
    cache := make(map[string]*cacheEntry)
    
    // 定期清理过期条目
    go func() {
        ticker := time.NewTicker(1 * time.Minute)
        defer ticker.Stop()
        
        for range ticker.C {
            now := time.Now().Unix()
            for k, v := range cache {
                if now-v.timestamp > 3600 { // 1小时过期
                    delete(cache, k)
                }
            }
        }
    }()
}
```

#### 使用局部变量

```go
func correctLocalVariable() {
    // 使用局部变量，函数返回后自动释放
    localCache := make(map[string][]byte)
    // ... 使用 localCache
    // 函数返回后，localCache 可以被 GC 回收
}
```

### 最佳实践

- 避免使用全局变量存储大量数据
- 如果必须使用全局缓存，实现 LRU 或 TTL 机制
- 定期清理不再需要的数据
- 使用局部变量替代全局变量

---

## 5. 切片泄漏

### 问题描述

切片泄漏发生在从大切片中切出小切片时。小切片仍然引用大切片的底层数组，导致大切片无法被 GC 回收。

### 错误示例

#### 示例 1：切片引用底层数组

```go
func sliceLeakExample() {
    // 创建一个包含 1000 个元素的大切片
    largeSlice := make([]int, 1000)
    
    // 只使用前 10 个元素
    smallSlice := largeSlice[:10]
    
    // 即使 largeSlice 不再使用，但由于 smallSlice 持有对底层数组的引用
    // 整个 largeSlice 的底层数组（1000 个元素）无法被 GC 回收
    // 只有 10 个元素被使用，但 1000 个元素都占用内存
}
```

**问题**：Go 的切片是引用类型，`smallSlice` 和 `largeSlice` 共享同一个底层数组。只要 `smallSlice` 存在，整个底层数组就无法被 GC 回收。

#### 示例 2：函数返回切片

```go
func getFirstN(n int) []int {
    largeSlice := make([]int, 10000)
    
    // 返回前 n 个元素
    // 但返回的切片仍然引用整个底层数组
    return largeSlice[:n]
}
```

**问题**：返回的切片引用整个 10000 个元素的底层数组，即使只需要 n 个元素。

### 正确做法

#### 使用 copy 创建独立切片

```go
func correctSliceCopy() {
    largeSlice := make([]int, 1000)
    
    // 使用 copy 创建独立的切片
    smallSlice := make([]int, 10)
    copy(smallSlice, largeSlice[:10])
    
    // 现在 smallSlice 是独立的，不引用 largeSlice 的底层数组
    // largeSlice 可以被 GC 回收
}
```

#### 函数返回时使用 copy

```go
func correctGetFirstN(n int) []int {
    largeSlice := make([]int, 10000)
    
    // 创建新切片并复制数据
    result := make([]int, n)
    copy(result, largeSlice[:n])
    return result
    // 返回的切片不引用 largeSlice，largeSlice 可以被 GC 回收
}
```

#### 使用 append 创建新切片

```go
func correctSliceAppend() {
    largeSlice := make([]int, 1000)
    
    // 使用 append 创建新切片
    smallSlice := append([]int(nil), largeSlice[:10]...)
    
    // smallSlice 是新的切片，不引用 largeSlice
}
```

### 检测方法

- 使用 `runtime.ReadMemStats()` 监控内存使用
- 使用 `pprof` 工具分析内存分配
- 注意切片容量（cap）和长度（len）的关系

---

## 6. 接口泄漏

### 问题描述

接口变量持有大对象的引用，即使只需要对象的一小部分数据。这通常发生在：
- 将大结构体赋值给 `interface{}`
- 接口 slice 或 map 存储大对象
- 接口持有整个对象而不是指针

### 错误示例

#### 示例 1：接口持有大结构体

```go
type LargeStruct struct {
    ID   int
    Data [1024 * 1024]byte // 1MB 数据
    Name string
}

func interfaceLeakExample() {
    var obj interface{}
    
    largeObj := LargeStruct{
        ID:   1,
        Name: "test",
    }
    
    // 接口变量持有整个结构体的引用
    obj = largeObj
    
    // 即使只需要 Name 字段，但由于接口持有整个 largeObj
    // 整个 1MB 的数据无法被 GC 回收
}
```

**问题**：将值类型赋值给接口时，会复制整个值。即使只需要部分字段，整个结构体都会占用内存。

#### 示例 2：接口 Slice 持有大对象

```go
func interfaceSliceLeak() {
    var objects []interface{}
    
    for i := 0; i < 100; i++ {
        obj := LargeStruct{
            ID:   i,
            Name: fmt.Sprintf("obj-%d", i),
        }
        objects = append(objects, obj)
        // 每个对象 1MB，100 个对象就是 100MB
    }
}
```

### 正确做法

#### 只存储需要的数据

```go
type SmallStruct struct {
    ID   int
    Name string
}

func correctInterfaceExample() {
    var obj interface{}
    
    largeObj := LargeStruct{
        ID:   1,
        Name: "test",
    }
    
    // 只存储需要的小结构体
    smallObj := SmallStruct{
        ID:   largeObj.ID,
        Name: largeObj.Name,
    }
    obj = smallObj
    
    // 现在接口只持有小结构体，大结构体可以被 GC 回收
}
```

#### 使用指针但及时清理

```go
func correctInterfaceWithPointer() {
    var obj *LargeStruct
    
    largeObj := &LargeStruct{
        ID:   1,
        Name: "test",
    }
    obj = largeObj
    
    // 使用完后显式设置为 nil
    fmt.Println(obj.Name)
    obj = nil
    // 现在 largeObj 可以被 GC 回收
}
```

#### 使用最小数据结构

```go
func correctMinimalInterface() {
    type MinimalData struct {
        ID   int
        Name string
    }
    
    objects := make([]MinimalData, 0, 100)
    
    for i := 0; i < 100; i++ {
        // 只存储必要的数据
        obj := MinimalData{
            ID:   i,
            Name: fmt.Sprintf("obj-%d", i),
        }
        objects = append(objects, obj)
        // 每个对象只有几十字节，而不是 1MB
    }
}
```

### 最佳实践

- 避免将大结构体直接赋值给接口
- 只存储实际需要的数据字段
- 使用指针时，及时清理引用
- 设计最小化的数据结构

---

## 内存泄漏检测工具

### 1. pprof

Go 标准库提供的性能分析工具：

```go
import _ "net/http/pprof"

func main() {
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    // ... 你的代码
}
```

使用方式：
```bash
# 查看 goroutine
go tool pprof http://localhost:6060/debug/pprof/goroutine

# 查看堆内存
go tool pprof http://localhost:6060/debug/pprof/heap

# 查看所有 goroutine 的堆栈
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

### 2. runtime 包

```go
import "runtime"

func monitorMemory() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    
    fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
    fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
    fmt.Printf("Sys = %v KB\n", m.Sys/1024)
    fmt.Printf("NumGC = %v\n", m.NumGC)
    fmt.Printf("NumGoroutine = %v\n", runtime.NumGoroutine())
}
```

### 3. goleak

第三方工具，专门用于检测 goroutine 泄漏：

```bash
go get go.uber.org/goleak
```

```go
import "go.uber.org/goleak"

func TestMyFunction(t *testing.T) {
    defer goleak.VerifyNone(t)
    // 你的测试代码
}
```

---

## 总结

### 常见内存泄漏模式

1. **Goroutine 泄漏**：最常见，goroutine 无法退出
2. **Channel 泄漏**：发送者/接收者阻塞
3. **定时器泄漏**：Ticker/Timer 未停止
4. **全局变量泄漏**：全局变量持有大量引用
5. **切片泄漏**：小切片引用大切片底层数组
6. **接口泄漏**：接口持有大对象完整引用

### 预防措施

1. **使用 defer**：确保资源被正确释放（ticker.Stop(), close(channel)）
2. **使用 context**：统一管理 goroutine 生命周期
3. **限制缓存大小**：实现 LRU 或 TTL 机制
4. **使用 copy**：创建独立切片，避免引用底层数组
5. **最小化数据结构**：只存储需要的数据
6. **定期监控**：使用 pprof 和 runtime 包监控内存和 goroutine

### 代码审查清单

- [ ] 所有启动的 goroutine 都有退出机制
- [ ] 所有 channel 都有明确的关闭逻辑
- [ ] 所有 Ticker/Timer 都调用了 Stop()
- [ ] 全局变量有大小限制或清理机制
- [ ] 返回的切片不引用大数组
- [ ] 接口只存储必要的数据

---

## 参考资源

- [Go Memory Leaks](https://go101.org/article/memory-leaking.html)
- [Effective Go - Memory Management](https://go.dev/doc/effective_go#memory)
- [pprof Documentation](https://pkg.go.dev/net/http/pprof)
- [Go Concurrency Patterns](https://go.dev/blog/pipelines)
