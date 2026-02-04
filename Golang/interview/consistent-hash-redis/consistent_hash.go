// Package chredis 基于一致性哈希将请求路由到多个 Redis 集群，实现请求分离与平滑扩缩容。
package chredis

import (
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

// Hash 计算 key 的哈希值，默认使用 CRC32。
type Hash func(data []byte) uint32

// Ring 一致性哈希环。
type Ring struct {
	hash     Hash
	replicas int               // 每个物理节点的虚拟节点数
	keys     []uint32          // 有序的环上位置
	ring     map[uint32]string // 位置 -> 节点名
	nodes    map[string]struct{}
	mu       sync.RWMutex
}

// New 创建一致性哈希环。replicas 为每个物理节点对应的虚拟节点数量，越大分布越均匀。
func New(replicas int, fn Hash) *Ring {
	if fn == nil {
		fn = crc32.ChecksumIEEE
	}
	return &Ring{
		hash:     fn,
		replicas: replicas,
		ring:     make(map[uint32]string),
		nodes:    make(map[string]struct{}),
	}
}

// Add 向环上添加节点（可多次调用以动态扩缩容）。
func (r *Ring) Add(nodes ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, node := range nodes {
		if _, ok := r.nodes[node]; ok {
			continue
		}
		r.nodes[node] = struct{}{}
		for i := 0; i < r.replicas; i++ {
			h := r.hash([]byte(node + strconv.Itoa(i)))
			r.keys = append(r.keys, h)
			r.ring[h] = node
		}
	}
	sort.Slice(r.keys, func(i, j int) bool { return r.keys[i] < r.keys[j] })
}

// Remove 从环上移除节点。
func (r *Ring) Remove(nodes ...string) {
	r.mu.Lock()
	defer r.mu.Unlock()
	for _, node := range nodes {
		if _, ok := r.nodes[node]; !ok {
			continue
		}
		delete(r.nodes, node)
		for i := 0; i < r.replicas; i++ {
			h := r.hash([]byte(node + strconv.Itoa(i)))
			delete(r.ring, h)
		}
	}
	r.rebuildKeys()
}

func (r *Ring) rebuildKeys() {
	r.keys = r.keys[:0]
	for k := range r.ring {
		r.keys = append(r.keys, k)
	}
	sort.Slice(r.keys, func(i, j int) bool { return r.keys[i] < r.keys[j] })
}

// Get 根据 key 返回应路由到的节点名。
func (r *Ring) Get(key string) string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	if len(r.keys) == 0 {
		return ""
	}
	h := r.hash([]byte(key))
	idx := sort.Search(len(r.keys), func(i int) bool { return r.keys[i] >= h })
	if idx == len(r.keys) {
		idx = 0
	}
	return r.ring[r.keys[idx]]
}

// Nodes 返回当前环上所有物理节点名。
func (r *Ring) Nodes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	out := make([]string, 0, len(r.nodes))
	for n := range r.nodes {
		out = append(out, n)
	}
	sort.Strings(out)
	return out
}
