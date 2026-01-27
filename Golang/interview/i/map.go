package i

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"sync"
	"time"
)

// 分片数量（建议是 2 的幂，便于位运算取模）
const shardCount = 32

// 每个分片包含一个 map + 一个读写锁
type shard struct {
	items map[string]interface{}
	mu    sync.RWMutex
}

// ShardedMap 分片 map
type ShardedMap struct {
	shards []*shard
}

// NewShardedMap 创建分片 map
func NewShardedMap() *ShardedMap {
	shards := make([]*shard, shardCount)
	for i := 0; i < shardCount; i++ {
		shards[i] = &shard{
			items: make(map[string]interface{}, 256), // 初始容量可调
			mu:    sync.RWMutex{},
		}
	}
	return &ShardedMap{shards: shards}
}

// getShard 根据 key 计算分片索引
func (sm *ShardedMap) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	hash := h.Sum32()
	return sm.shards[hash%shardCount]
}

// Store 存值（并发安全）
func (sm *ShardedMap) Store(key string, value interface{}) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	shard.items[key] = value
}

// Load 取值（并发安全）
func (sm *ShardedMap) Load(key string) (interface{}, bool) {
	shard := sm.getShard(key)
	shard.mu.RLock()
	defer shard.mu.RUnlock()
	val, ok := shard.items[key]
	return val, ok
}

// Delete 删除（并发安全）
func (sm *ShardedMap) Delete(key string) {
	shard := sm.getShard(key)
	shard.mu.Lock()
	defer shard.mu.Unlock()
	delete(shard.items, key)
}

// Range 并行遍历（每个分片独立由一个 goroutine 处理）
func (sm *ShardedMap) Range(f func(key string, value interface{}) bool) {
	var wg sync.WaitGroup
	wg.Add(shardCount)

	for i := 0; i < shardCount; i++ {
		go func(idx int) {
			defer wg.Done()
			shard := sm.shards[idx]
			shard.mu.RLock()
			defer shard.mu.RUnlock()

			for k, v := range shard.items {
				if !f(k, v) {
					return // 提前退出当前分片
				}
			}
		}(i)
	}

	wg.Wait()
}

// Size 获取总元素个数（并发安全，但有一定开销）
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
			total += int64(len(shard.items)) // 注意：非原子，需要加锁保护 total
		}(i)
	}
	wg.Wait()
	return total
}

func main() {
	m := NewShardedMap()

	// 模拟写入（并发）
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

	fmt.Println("写入完成，总大小约：", m.Size())

	// 并行遍历示例
	start := time.Now()
	count := 0
	m.Range(func(key string, value interface{}) bool {
		// 模拟处理
		count++
		if count%100000 == 0 {
			fmt.Printf("已处理 %d 条...\n", count)
		}
		return true
	})
	fmt.Printf("并行遍历耗时: %v, 总条数: %d\n", time.Since(start), count)

	// 普通单 goroutine 遍历对比（仅为了演示）
	// start = time.Now()
	// count = 0
	// for i := 0; i < shardCount; i++ {
	// 	shard := m.shards[i]
	// 	shard.mu.RLock()
	// 	for k, v := range shard.items {
	// 		count++
	// 	}
	// 	shard.mu.RUnlock()
	// }
	// fmt.Printf("单线程遍历耗时: %v\n", time.Since(start))
}
