package main

import (
	"fmt"
	"sync"
)

// 错误示例：全局变量持有引用导致内存泄漏
// 问题：全局变量持有大对象的引用，导致这些对象无法被 GC 回收

var globalCache = make(map[string][]byte)
var globalCacheMu sync.Mutex

// 错误示例：全局 map 不断增长，从不清理
func globalMapLeak() {
	for i := 0; i < 10000; i++ {
		key := fmt.Sprintf("key-%d", i)
		// 创建 1MB 的数据
		data := make([]byte, 1024*1024)
		
		globalCacheMu.Lock()
		globalCache[key] = data
		globalCacheMu.Unlock()
		
		// 即使不再使用这些数据，由于全局变量持有引用
		// GC 无法回收这些内存
	}
}

// 错误示例：全局 slice 持有引用
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

// 错误示例：全局 channel 持有引用
var globalChan = make(chan *LargeObject, 1000)

func globalChannelLeak() {
	for i := 0; i < 1000; i++ {
		obj := &LargeObject{ID: i}
		globalChan <- obj
		// 如果 channel 中的数据不被消费
		// 这些对象会一直占用内存
	}
}

// 正确示例：使用局部变量，函数返回后自动释放
func correctLocalVariable() {
	localCache := make(map[string][]byte)
	
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		data := make([]byte, 1024*1024)
		localCache[key] = data
	}
	
	// 函数返回后，localCache 可以被 GC 回收
}

// 正确示例：定期清理全局缓存
func correctGlobalCacheWithCleanup() {
	// 使用带过期时间的缓存
	type cacheEntry struct {
		data      []byte
		timestamp int64
	}
	
	cache := make(map[string]*cacheEntry)
	
	// 定期清理过期条目
	go func() {
		for {
			// 清理逻辑
			// 删除过期的条目
		}
	}()
}

// 正确示例：使用弱引用或限制大小
func correctLimitedCache() {
	const maxSize = 1000
	cache := make(map[string][]byte)
	
	// 限制缓存大小
	if len(cache) >= maxSize {
		// 删除最旧的条目
		for k := range cache {
			delete(cache, k)
			break
		}
	}
}

// 正确示例：显式清理
func correctExplicitCleanup() {
	globalCacheMu.Lock()
	defer globalCacheMu.Unlock()
	
	// 清理不再需要的条目
	for key := range globalCache {
		// 根据业务逻辑决定是否删除
		delete(globalCache, key)
	}
}

func main() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// globalMapLeak()
	// globalSliceLeak()
	// globalChannelLeak()
	
	// 正确示例
	correctLocalVariable()
	correctGlobalCacheWithCleanup()
	correctLimitedCache()
	correctExplicitCleanup()
}
