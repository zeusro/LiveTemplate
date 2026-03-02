# IncMap —— 支持随时递增的并发安全 Map

## 题目描述

封装一个随时可以递增的 map，支持对同一个 key 的：

- **读取**（Load）
- **写入**（Store）
- **删除**（Delete）
- **递增**（Incr / IncrBy）

所有操作必须是并发安全的。

---

## 实现方案：分片 + per-key Seqlock

### 数据结构

```go
// entry 持有单个 key 的值与顺序锁版本号
type entry struct {
    value   atomic.Value // 始终存储 int64
    version atomic.Int64 // seqlock 版本号：偶数=稳定，奇数=写中
}

type shard struct {
    mu    sync.RWMutex
    items map[string]*entry
}

type IncMap struct {
    shards [32]*shard
}
```

### 三层并发控制

```
┌──────────────────────────────────────────────────┐
│  第1层：分片（sharding）                           │
│  FNV hash(key) % 32 → 定位到对应 shard            │
│  不同 key 大概率落在不同 shard，天然隔离            │
├──────────────────────────────────────────────────┤
│  第2层：shard.mu（RWMutex）                        │
│  仅保护 items map 的结构变更（新增/删除 entry）     │
│  持锁时间极短，只做一次 map 读写                   │
├──────────────────────────────────────────────────┤
│  第3层：entry seqlock（atomic.Value + version）    │
│  保护单个 key 的值读写，完全无锁自旋               │
│  读操作不阻塞写操作，写操作互斥但不阻塞读          │
└──────────────────────────────────────────────────┘
```

---

## 核心：Seqlock（顺序锁）原理

Seqlock 是 Linux 内核中经典的读写并发模式，核心思想：

- **版本号为偶数**：数据处于稳定状态，读写均可进入
- **版本号为奇数**：有写操作正在进行，读者自旋等待

### 写协议（store / incrBy）

```
1. 读取当前版本号 v（必须为偶数，否则自旋）
2. CAS(v → v+1)，将版本号改为奇数，"宣告"写开始
   └── CAS 失败说明有其他写者抢先，重试
3. 执行写操作（此时读者会自旋等待）
4. version.Store(v+2)，恢复为偶数，"释放"写权
```

```go
func (e *entry) store(val int64) {
    for {
        v := e.version.Load()
        if v&1 == 1 { runtime.Gosched(); continue }   // 等待其他写者
        if !e.version.CompareAndSwap(v, v+1) { continue } // 抢写权
        e.value.Store(val)          // 临界区：写数据
        e.version.Store(v + 2)      // 释放写权
        return
    }
}
```

### 读协议（load）

```
1. 读取版本号 v（必须为偶数，否则自旋）
2. 读取 atomic.Value 中的值
3. 再次读取版本号，若与步骤1相同则读取成功
   └── 若版本号变化，说明读期间发生了写，重试
```

```go
func (e *entry) load() int64 {
    for {
        v := e.version.Load()
        if v&1 == 1 { runtime.Gosched(); continue }   // 写者持锁中
        val := e.value.Load().(int64)
        if e.version.Load() == v { return val }        // 读到一致数据
        // 版本变化，重试
    }
}
```

### IncrBy 的原子性保证

`IncrBy` 是"读-改-写"三步操作，必须在持有写权（版本号为奇数）的情况下完成：

```go
func (e *entry) incrBy(delta int64) int64 {
    for {
        v := e.version.Load()
        if v&1 == 1 { runtime.Gosched(); continue }
        if !e.version.CompareAndSwap(v, v+1) { continue } // 独占写权
        old := e.value.Load().(int64)  // 安全读旧值（此时无其他写者）
        newVal := old + delta
        e.value.Store(newVal)
        e.version.Store(v + 2)
        return newVal
    }
}
```

写权（奇数版本号）相当于互斥锁：同一时刻只有一个 goroutine 能通过 CAS 将版本号从偶数改为奇数，从而独占对该 entry 的修改权。

---

## 为什么不用 sync/atomic 直接做

| 方案 | 问题 |
|------|------|
| `atomic.AddInt64` | 需要固定内存地址；map 的 value 不可寻址，无法直接用 |
| `sync.Map` + `*atomic.Int64` | 可行，但每个 key 多一次堆分配，且无法演示 seqlock 设计 |
| 单把 `sync.Mutex` | 所有 key 共用一把锁，高并发下竞争严重 |
| **本方案（seqlock）** | 读操作完全不阻塞，写操作仅锁定单个 entry，粒度最细 |

---

## getOrCreate 的 Double-Check 模式

```go
func (m *IncMap) getOrCreate(key string) *entry {
    s := m.getShard(key)

    s.mu.RLock()                 // 快速路径：读锁查找
    e, ok := s.items[key]
    s.mu.RUnlock()
    if ok { return e }

    s.mu.Lock()                  // 慢路径：写锁创建
    defer s.mu.Unlock()
    if e, ok = s.items[key]; ok { // Double-check：防止并发重复创建
        return e
    }
    e = newEntry(0)
    s.items[key] = e
    return e
}
```

不做 Double-check 的后果：两个 goroutine 同时通过读锁检查"key 不存在"，升级为写锁后各自创建 entry，后创建的会覆盖先创建的，导致已完成的递增值丢失。

---

## 各操作的锁持有情况

| 操作 | shard.mu 持有时长 | entry seqlock |
|------|------------------|---------------|
| `Load` | 极短（map 查找） | 读协议自旋 |
| `Store` | 极短（getOrCreate） | 写协议 CAS |
| `Delete` | 极短（map delete） | 无 |
| `Incr` / `IncrBy` | 极短（getOrCreate） | 写协议 CAS |

**关键优势**：`shard.mu` 从不在值读写期间持有；不同 key（不同 shard）的操作完全并行。

---

## 与单锁版本的对比

| 维度 | 单锁 IncrMap | 分片 + Seqlock IncMap |
|------|--------------|-----------------------|
| 锁粒度 | 整个 map 一把 RWMutex | 32 分片 × 每 key 一个 seqlock |
| 读并发 | 同分片读互不阻塞 | 读不阻塞写（seqlock） |
| 写并发 | 全局写锁，所有操作串行 | 不同 key 完全并行 |
| 内存开销 | 低 | 每 key 多一个 `entry` 对象 |
| 实现复杂度 | 简单 | 较高（seqlock + double-check） |
| 适用场景 | key 少、低并发 | key 多、高并发计数器 |

---

## 常见追问

**Q：seqlock 的读自旋会不会浪费 CPU？**

在极低概率的写冲突场景下（写操作极短，通常只有几十纳秒），自旋几乎感知不到。若写操作耗时较长，应改用 `sync.Mutex`。

**Q：`atomic.Value` 为什么比 `atomic.Int64` 更通用？**

`atomic.Int64` 只能存 `int64`；`atomic.Value` 可存任意相同类型的值（如结构体快照），更适合扩展为泛型 map。本题用 `int64` 是为了演示该模式。

**Q：Delete 后再 Incr 同一个 key 会怎样？**

`Delete` 从 map 中移除 entry，`Incr` 会通过 `getOrCreate` 创建新 entry（从 0 开始），不会 panic。

**Q：如何扩展支持 TTL（过期）？**

在 `entry` 中额外存储 `expireAt int64`（Unix 纳秒），`Load` 时检查是否过期，过期则返回 `(0, false)` 并惰性删除（或后台 GC 扫描）。
