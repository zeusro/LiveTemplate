package incrmap

import (
	"hash/fnv"
	"runtime"
	"sync"
	"sync/atomic"
)

const shardCount = 32

// entry 持有单个 key 的值与顺序锁版本号。
//
// version 用作 seqlock：
//   - 偶数：稳定状态，可安全读写
//   - 奇数：写操作进行中，读者需自旋等待
//
// 这样 shard.mu 只需要保护 items map 的结构变更（新增/删除 key），
// 而对已存在 key 的读写完全通过 entry 上的原子操作完成，无需持有 shard 锁。
type entry struct {
	value   atomic.Value // 始终存储 int64，由 seqlock 保护一致性
	version atomic.Int64 // seqlock 版本号
}

func newEntry(v int64) *entry {
	e := &entry{}
	// atomic.Value 要求首次 Store 前不能 Load，这里提前初始化。
	e.value.Store(v)
	return e
}

// load 使用 seqlock 读协议：读前后各读一次版本号，一致则返回。
func (e *entry) load() int64 {
	for {
		v := e.version.Load()
		if v&1 == 1 {
			// 有写操作正在进行，让出 CPU 后重试
			runtime.Gosched()
			continue
		}
		val := e.value.Load().(int64)
		// 若版本号未变，说明读期间无写操作，读到的值是一致的
		if e.version.Load() == v {
			return val
		}
	}
}

// store 使用 seqlock 写协议：CAS 将版本号设为奇数以独占写权，写完后恢复为偶数。
func (e *entry) store(val int64) {
	for {
		v := e.version.Load()
		if v&1 == 1 {
			runtime.Gosched()
			continue
		}
		if !e.version.CompareAndSwap(v, v+1) {
			continue
		}
		// 持有写权（version 为奇数），其他写者无法进入
		e.value.Store(val)
		e.version.Store(v + 2) // 释放写权，版本号推进
		return
	}
}

// incrBy 原子地将值增加 delta，返回新值。
// 同样使用 seqlock 独占写，再读出旧值后加 delta 写回。
func (e *entry) incrBy(delta int64) int64 {
	for {
		v := e.version.Load()
		if v&1 == 1 {
			runtime.Gosched()
			continue
		}
		if !e.version.CompareAndSwap(v, v+1) {
			continue
		}
		old := e.value.Load().(int64)
		newVal := old + delta
		e.value.Store(newVal)
		e.version.Store(v + 2)
		return newVal
	}
}

// -----------------------------------------------------------------------

type shard struct {
	mu    sync.RWMutex
	items map[string]*entry
}

// IncMap 是使用分片 + per-key seqlock 实现的并发安全递增 map。
//
// 与单锁 IncrMap 的关键区别：
//   - shard.mu 仅保护 items map 的结构（新增/删除 key），粒度更粗但持锁时间极短
//   - 对同一 key 的读/写/递增通过 entry 的 seqlock 完成，彼此互不阻塞
//   - 不同 shard 的操作完全并行，无锁竞争
type IncMap struct {
	shards [shardCount]*shard
}

func NewIncMap() *IncMap {
	m := &IncMap{}
	for i := range m.shards {
		m.shards[i] = &shard{items: make(map[string]*entry)}
	}
	return m
}

func (m *IncMap) getShard(key string) *shard {
	h := fnv.New32a()
	h.Write([]byte(key))
	return m.shards[h.Sum32()%shardCount]
}

// getOrCreate 先以读锁快速查找，找不到时升级为写锁创建。
// Double-check 防止在读锁释放到写锁获取之间的并发创建。
func (m *IncMap) getOrCreate(key string) *entry {
	s := m.getShard(key)

	s.mu.RLock()
	e, ok := s.items[key]
	s.mu.RUnlock()
	if ok {
		return e
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	if e, ok = s.items[key]; ok {
		return e
	}
	e = newEntry(0)
	s.items[key] = e
	return e
}

// Load 读取 key 的值。若 key 不存在返回 (0, false)。
// 仅在查找 entry 时持有读锁，值读取本身无锁（seqlock）。
func (m *IncMap) Load(key string) (int64, bool) {
	s := m.getShard(key)
	s.mu.RLock()
	e, ok := s.items[key]
	s.mu.RUnlock()
	if !ok {
		return 0, false
	}
	return e.load(), true
}

// Store 设置 key 的值，若 key 不存在则创建。
func (m *IncMap) Store(key string, value int64) {
	m.getOrCreate(key).store(value)
}

// Delete 删除 key。
func (m *IncMap) Delete(key string) {
	s := m.getShard(key)
	s.mu.Lock()
	delete(s.items, key)
	s.mu.Unlock()
}

// Incr 将 key 的值加 1，返回新值。key 不存在时从 0 开始。
func (m *IncMap) Incr(key string) int64 {
	return m.IncrBy(key, 1)
}

// IncrBy 将 key 的值加 delta，返回新值。
func (m *IncMap) IncrBy(key string, delta int64) int64 {
	return m.getOrCreate(key).incrBy(delta)
}

// Len 返回所有分片中 key 的总数。
func (m *IncMap) Len() int {
	total := 0
	for i := range m.shards {
		s := m.shards[i]
		s.mu.RLock()
		total += len(s.items)
		s.mu.RUnlock()
	}
	return total
}

// Ensure IncMap satisfies the same interface as the original IncrMap.
var _ interface {
	Load(string) (int64, bool)
	Store(string, int64)
	Delete(string)
	Incr(string) int64
	IncrBy(string, int64) int64
	Len() int
} = (*IncMap)(nil)

// keep atomic import used (atomic.Value / atomic.Int64 are used inside entry)
var _ = atomic.Value{}
