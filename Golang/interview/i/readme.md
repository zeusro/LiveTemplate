
---

## 分片Map算法 (ShardedMap)

本文档介绍 `Golang/interview/i/map.go` 中实现的分片Map算法，这是一种高性能的并发安全Map实现。

### 算法概述

分片Map（ShardedMap）是一种通过将数据分散到多个分片（shard）中来提高并发性能的数据结构。相比使用单一全局锁的传统并发Map，分片Map通过细粒度锁机制显著降低了锁竞争，提升了多线程环境下的性能。

### 核心思想

1. **分片策略**：将整个Map分成多个独立的分片（默认32个）
2. **哈希分布**：使用哈希函数将key均匀分布到不同的分片
3. **细粒度锁**：每个分片拥有独立的读写锁，减少锁竞争
4. **并行操作**：支持并行遍历和统计操作

### 数据结构

```go
// 分片结构：每个分片包含一个map和一把读写锁
type shard struct {
    items map[string]interface{}
    mu    sync.RWMutex
}

// 分片Map：包含多个分片
type ShardedMap struct {
    shards []*shard  // 默认32个分片
}
```

### 关键算法

#### 1. 分片选择算法

使用FNV-1a哈希算法计算key的分片索引：

```go
func (sm *ShardedMap) getShard(key string) *shard {
    h := fnv.New32a()
    h.Write([]byte(key))
    hash := h.Sum32()
    return sm.shards[hash % shardCount]  // 取模运算确定分片
}
```

**算法特点：**
- 使用FNV-1a哈希算法，速度快且分布均匀
- 分片数量建议为2的幂（如32），便于位运算优化
- 相同key总是映射到同一个分片，保证一致性

#### 2. 存储操作 (Store)

```go
func (sm *ShardedMap) Store(key string, value interface{}) {
    shard := sm.getShard(key)  // 1. 定位到对应分片
    shard.mu.Lock()            // 2. 获取写锁（仅锁定该分片）
    defer shard.mu.Unlock()
    shard.items[key] = value   // 3. 写入数据
}
```

**性能优势：**
- 只锁定一个分片，其他分片可继续并发操作
- 不同分片的操作完全并行，无锁竞争

#### 3. 读取操作 (Load)

```go
func (sm *ShardedMap) Load(key string) (interface{}, bool) {
    shard := sm.getShard(key)
    shard.mu.RLock()           // 使用读锁，支持并发读
    defer shard.mu.RUnlock()
    val, ok := shard.items[key]
    return val, ok
}
```

**性能优势：**
- 使用读锁，同一分片的多个读操作可以并发执行
- 不同分片的读操作完全并行

#### 4. 并行遍历算法 (Range)

```go
func (sm *ShardedMap) Range(f func(key string, value interface{}) bool) {
    var wg sync.WaitGroup
    wg.Add(shardCount)
    
    // 为每个分片启动一个goroutine
    for i := 0; i < shardCount; i++ {
        go func(idx int) {
            defer wg.Done()
            shard := sm.shards[idx]
            shard.mu.RLock()
            defer shard.mu.RUnlock()
            
            // 遍历当前分片
            for k, v := range shard.items {
                if !f(k, v) {
                    return  // 提前退出
                }
            }
        }(i)
    }
    
    wg.Wait()  // 等待所有goroutine完成
}
```

**算法特点：**
- 每个分片由独立的goroutine处理，实现真正的并行遍历
- 使用WaitGroup同步，确保所有分片遍历完成
- 支持提前退出（通过回调函数返回false）

#### 5. 并行统计算法 (Size)

```go
func (sm *ShardedMap) Size() int64 {
    var total int64
    var wg sync.WaitGroup
    wg.Add(shardCount)
    
    for i := 0; i < shardCount; i++ {
        go func(idx int) {
            defer wg.Done()
            shard := sm.shards[idx]
            shard.mu.RLock()
            defer shard.mu.RUnlock()
            total += int64(len(shard.items))  // 累加各分片大小
        }(i)
    }
    wg.Wait()
    return total
}
```

**注意：** 当前实现中`total`的累加不是原子操作，在高并发场景下可能存在数据竞争。实际应用中应使用`atomic.AddInt64`或互斥锁保护。

### 性能分析

#### 优势

1. **降低锁竞争**：32个分片意味着理论上可以将锁竞争降低到原来的1/32
2. **提高并发度**：不同分片的操作完全并行，互不干扰
3. **读操作优化**：使用读写锁，读操作可以并发执行
4. **并行遍历**：Range操作并行处理所有分片，大幅提升遍历速度

#### 适用场景

- **高并发读写**：大量并发的Store/Load操作
- **大规模数据**：数据量较大，需要高效遍历
- **读多写少**：充分利用读写锁的优势

#### 局限性

1. **内存开销**：需要维护多个分片和锁，内存占用略高于单Map
2. **跨分片操作**：无法保证跨分片操作的原子性
3. **Size统计**：当前实现存在并发安全问题，需要改进

### 算法复杂度

- **时间复杂度：**
  - Store/Load/Delete: O(1) 平均情况（哈希计算 + 单次Map操作）
  - Range: O(n) 其中n为总元素数，但通过并行处理，实际耗时约为 O(n/32)
  - Size: O(shardCount) 并行统计，实际耗时约为单次锁获取时间

- **空间复杂度：**
  - O(n + shardCount) 其中n为元素数量，shardCount为分片数

### 优化建议

1. **分片数量调优**：根据实际并发量和CPU核心数调整分片数量
2. **哈希算法选择**：FNV-1a适合大多数场景，也可考虑更快的哈希算法
3. **Size方法改进**：使用原子操作或定期维护计数器
4. **初始容量**：根据预期数据量设置合理的初始容量，减少扩容开销

### 代码示例

```go
// 创建分片Map
m := NewShardedMap()

// 并发写入
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(id int) {
        defer wg.Done()
        for j := 0; j < 500; j++ {
            key := fmt.Sprintf("user-%d-%d", id, j)
            m.Store(key, rand.Intn(10000))
        }
    }(i)
}
wg.Wait()

// 并行遍历
m.Range(func(key string, value interface{}) bool {
    // 处理每个键值对
    return true  // 返回false可提前退出
})
```

### 总结

分片Map算法通过将数据分散到多个独立的分片，并使用细粒度锁机制，有效降低了并发环境下的锁竞争，显著提升了高并发场景下的性能。这是一种经典的"分而治之"思想在并发数据结构中的应用。
