package incrmap

import (
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// ─── per-key entry: value + seqlock + immutable createdAt + mutable updatedAt ─

type tsEntry struct {
	version   atomic.Int64 // seqlock: even=stable, odd=writing
	value     atomic.Int64
	createdAt int64        // unix nano; written once before publication, never changed
	updatedAt atomic.Int64 // unix nano; updated on every write via seqlock
}

func newTsEntry(val int64) *tsEntry {
	now := time.Now().UnixNano()
	e := &tsEntry{createdAt: now}
	e.value.Store(val)
	e.updatedAt.Store(now)
	return e
}

// tsLoad reads value, createdAt, updatedAt atomically (seqlock read protocol).
func (e *tsEntry) tsLoad() (val, createdAt, updatedAt int64) {
	for {
		v := e.version.Load()
		if v&1 == 1 {
			runtime.Gosched()
			continue
		}
		val = e.value.Load()
		ua := e.updatedAt.Load()
		if e.version.Load() == v {
			return val, e.createdAt, ua
		}
	}
}

// tsStore writes val and refreshes updatedAt under seqlock.
func (e *tsEntry) tsStore(val int64) {
	now := time.Now().UnixNano()
	for {
		v := e.version.Load()
		if v&1 == 1 {
			runtime.Gosched()
			continue
		}
		if !e.version.CompareAndSwap(v, v+1) {
			continue
		}
		e.value.Store(val)
		e.updatedAt.Store(now)
		e.version.Store(v + 2)
		return
	}
}

// tsIncrBy atomically adds delta, refreshes updatedAt, returns new value.
func (e *tsEntry) tsIncrBy(delta int64) int64 {
	now := time.Now().UnixNano()
	for {
		v := e.version.Load()
		if v&1 == 1 {
			runtime.Gosched()
			continue
		}
		if !e.version.CompareAndSwap(v, v+1) {
			continue
		}
		newVal := e.value.Load() + delta
		e.value.Store(newVal)
		e.updatedAt.Store(now)
		e.version.Store(v + 2)
		return newVal
	}
}

// keep atomic import used
var _ = atomic.Value{}

// ─── Left-Leaning Red-Black tree (LLRB) ──────────────────────────────────────
//
// Reference: Sedgewick, "Left-leaning Red-Black Trees" (2008).
// Invariants:
//   - No two consecutive red links on any path.
//   - Perfect black-link balance.
//   - Red links always lean left.

const (
	llrbRed   = true
	llrbBlack = false
)

type llrbNode struct {
	key   string
	entry *tsEntry
	color bool
	left  *llrbNode
	right *llrbNode
	sz    int // subtree node count
}

func newLLRBNode(key string, val int64) *llrbNode {
	return &llrbNode{key: key, entry: newTsEntry(val), color: llrbRed, sz: 1}
}

func isRed(n *llrbNode) bool        { return n != nil && n.color == llrbRed }
func nodeSize(n *llrbNode) int      { if n == nil { return 0 }; return n.sz }
func updateSz(n *llrbNode)          { n.sz = nodeSize(n.left) + nodeSize(n.right) + 1 }

func rotateLeft(h *llrbNode) *llrbNode {
	x := h.right
	h.right = x.left
	x.left = h
	x.color = h.color
	h.color = llrbRed
	updateSz(h)
	updateSz(x)
	return x
}

func rotateRight(h *llrbNode) *llrbNode {
	x := h.left
	h.left = x.right
	x.right = h
	x.color = h.color
	h.color = llrbRed
	updateSz(h)
	updateSz(x)
	return x
}

func flipColors(h *llrbNode) {
	h.color = !h.color
	h.left.color = !h.left.color
	h.right.color = !h.right.color
}

func llrbBalance(h *llrbNode) *llrbNode {
	if isRed(h.right) && !isRed(h.left) {
		h = rotateLeft(h)
	}
	if isRed(h.left) && isRed(h.left.left) {
		h = rotateRight(h)
	}
	if isRed(h.left) && isRed(h.right) {
		flipColors(h)
	}
	updateSz(h)
	return h
}

// llrbPut inserts key=val or calls updateFn on the existing entry.
// Returns (new root, wasCreated).
func llrbPut(h *llrbNode, key string, val int64, updateFn func(*tsEntry)) (*llrbNode, bool) {
	if h == nil {
		return newLLRBNode(key, val), true
	}
	var created bool
	cmp := strings.Compare(key, h.key)
	switch {
	case cmp < 0:
		h.left, created = llrbPut(h.left, key, val, updateFn)
	case cmp > 0:
		h.right, created = llrbPut(h.right, key, val, updateFn)
	default:
		updateFn(h.entry)
	}
	return llrbBalance(h), created
}

func llrbGet(h *llrbNode, key string) *tsEntry {
	for h != nil {
		cmp := strings.Compare(key, h.key)
		switch {
		case cmp < 0:
			h = h.left
		case cmp > 0:
			h = h.right
		default:
			return h.entry
		}
	}
	return nil
}

func llrbMin(h *llrbNode) *llrbNode {
	for h.left != nil {
		h = h.left
	}
	return h
}

func llrbMax(h *llrbNode) *llrbNode {
	for h.right != nil {
		h = h.right
	}
	return h
}

func moveRedLeft(h *llrbNode) *llrbNode {
	flipColors(h)
	if isRed(h.right.left) {
		h.right = rotateRight(h.right)
		h = rotateLeft(h)
		flipColors(h)
	}
	return h
}

func moveRedRight(h *llrbNode) *llrbNode {
	flipColors(h)
	if isRed(h.left.left) {
		h = rotateRight(h)
		flipColors(h)
	}
	return h
}

func deleteMin(h *llrbNode) *llrbNode {
	if h.left == nil {
		return nil
	}
	if !isRed(h.left) && !isRed(h.left.left) {
		h = moveRedLeft(h)
	}
	h.left = deleteMin(h.left)
	return llrbBalance(h)
}

// llrbDelete removes key. Returns (new root, wasDeleted).
func llrbDelete(h *llrbNode, key string) (*llrbNode, bool) {
	if h == nil {
		return nil, false
	}
	var deleted bool
	if strings.Compare(key, h.key) < 0 {
		if h.left == nil {
			return h, false
		}
		if !isRed(h.left) && !isRed(h.left.left) {
			h = moveRedLeft(h)
		}
		h.left, deleted = llrbDelete(h.left, key)
	} else {
		if isRed(h.left) {
			h = rotateRight(h)
		}
		if strings.Compare(key, h.key) == 0 && h.right == nil {
			return nil, true
		}
		if h.right != nil && !isRed(h.right) && !isRed(h.right.left) {
			h = moveRedRight(h)
		}
		if strings.Compare(key, h.key) == 0 {
			m := llrbMin(h.right)
			h.key = m.key
			h.entry = m.entry
			h.right = deleteMin(h.right)
			deleted = true
		} else {
			h.right, deleted = llrbDelete(h.right, key)
		}
	}
	return llrbBalance(h), deleted
}

func llrbFloor(h *llrbNode, key string) *llrbNode {
	if h == nil {
		return nil
	}
	cmp := strings.Compare(key, h.key)
	if cmp == 0 {
		return h
	}
	if cmp < 0 {
		return llrbFloor(h.left, key)
	}
	t := llrbFloor(h.right, key)
	if t != nil {
		return t
	}
	return h
}

func llrbCeiling(h *llrbNode, key string) *llrbNode {
	if h == nil {
		return nil
	}
	cmp := strings.Compare(key, h.key)
	if cmp == 0 {
		return h
	}
	if cmp > 0 {
		return llrbCeiling(h.right, key)
	}
	t := llrbCeiling(h.left, key)
	if t != nil {
		return t
	}
	return h
}

// llrbRange calls fn in-order for keys in [lo, hi]. Returns false if fn stopped early.
func llrbRange(h *llrbNode, lo, hi string, fn func(EntryInfo) bool) bool {
	if h == nil {
		return true
	}
	cmpLo := strings.Compare(lo, h.key)
	cmpHi := strings.Compare(hi, h.key)
	if cmpLo < 0 {
		if !llrbRange(h.left, lo, hi, fn) {
			return false
		}
	}
	if cmpLo <= 0 && cmpHi >= 0 {
		if !fn(makeEntryInfo(h.key, h.entry)) {
			return false
		}
	}
	if cmpHi > 0 {
		if !llrbRange(h.right, lo, hi, fn) {
			return false
		}
	}
	return true
}

// llrbWalk calls fn for every node in ascending key order.
func llrbWalk(h *llrbNode, fn func(EntryInfo) bool) bool {
	if h == nil {
		return true
	}
	if !llrbWalk(h.left, fn) {
		return false
	}
	if !fn(makeEntryInfo(h.key, h.entry)) {
		return false
	}
	return llrbWalk(h.right, fn)
}

func llrbCollectKeys(h *llrbNode, out *[]string) {
	if h == nil {
		return
	}
	llrbCollectKeys(h.left, out)
	*out = append(*out, h.key)
	llrbCollectKeys(h.right, out)
}

// ─── public types ─────────────────────────────────────────────────────────────

// EntryInfo is a value-copy snapshot of a single entry returned by query methods.
type EntryInfo struct {
	Key       string
	Value     int64
	CreatedAt time.Time // immutable; set when the key was first inserted
	UpdatedAt time.Time // refreshed on every Store / IncrBy
}

func makeEntryInfo(key string, e *tsEntry) EntryInfo {
	val, ca, ua := e.tsLoad()
	return EntryInfo{
		Key:       key,
		Value:     val,
		CreatedAt: time.Unix(0, ca),
		UpdatedAt: time.Unix(0, ua),
	}
}

// ─── TimestampTreeMap ─────────────────────────────────────────────────────────

// TimestampTreeMap is a concurrency-safe ordered map backed by a Left-Leaning
// Red-Black tree (LLRB). It differs from IncMap in three key ways:
//
//  1. Keys are stored in sorted order, enabling O(log n) ordered operations:
//     Min, Max, Floor, Ceiling, and Range scans.
//
//  2. Every entry carries two timestamps:
//     • CreatedAt – set once when the key is first inserted, never changes.
//     • UpdatedAt – refreshed on every Store / IncrBy, enabling time-filtered
//     iteration via UpdatedAfter.
//
//  3. Concurrency model differs from sharded IncMap:
//     • A single sync.RWMutex guards tree structural changes (insert / delete /
//     LLRB rebalancing). Only insertions and deletions hold the write lock.
//     • Value-only updates (Store / IncrBy to existing keys) find the entry
//     under the read lock, then mutate it via the per-entry seqlock – so
//     concurrent value updates to distinct keys never block each other.
//     • Read operations (Load, Min, Max, …) hold only the read lock + seqlock.
type TimestampTreeMap struct {
	mu   sync.RWMutex
	root *llrbNode
}

func NewTimestampTreeMap() *TimestampTreeMap {
	return &TimestampTreeMap{}
}

// Load returns the current value for key and whether the key exists.
func (m *TimestampTreeMap) Load(key string) (int64, bool) {
	m.mu.RLock()
	e := llrbGet(m.root, key)
	m.mu.RUnlock()
	if e == nil {
		return 0, false
	}
	val, _, _ := e.tsLoad()
	return val, true
}

// LoadInfo returns the full EntryInfo (value + both timestamps) for key.
func (m *TimestampTreeMap) LoadInfo(key string) (EntryInfo, bool) {
	m.mu.RLock()
	e := llrbGet(m.root, key)
	m.mu.RUnlock()
	if e == nil {
		return EntryInfo{}, false
	}
	return makeEntryInfo(key, e), true
}

// Store sets key to value. If key already exists, CreatedAt is preserved and
// UpdatedAt is refreshed. If the key is new, both timestamps are set to now.
func (m *TimestampTreeMap) Store(key string, value int64) {
	// Fast path: key already in tree – only a read lock is needed.
	m.mu.RLock()
	e := llrbGet(m.root, key)
	m.mu.RUnlock()
	if e != nil {
		e.tsStore(value)
		return
	}

	// Slow path: insert a new node (tree structural change requires write lock).
	m.mu.Lock()
	// Double-check: another goroutine may have inserted between RUnlock and Lock.
	e = llrbGet(m.root, key)
	if e != nil {
		m.mu.Unlock()
		e.tsStore(value)
		return
	}
	m.root, _ = llrbPut(m.root, key, value, func(ex *tsEntry) { ex.tsStore(value) })
	m.root.color = llrbBlack
	m.mu.Unlock()
}

// Delete removes key from the map. No-op if the key does not exist.
func (m *TimestampTreeMap) Delete(key string) {
	m.mu.Lock()
	var deleted bool
	m.root, deleted = llrbDelete(m.root, key)
	if deleted && m.root != nil {
		m.root.color = llrbBlack
	}
	m.mu.Unlock()
}

// Incr increments key by 1 and returns the new value. Creates key=1 if absent.
func (m *TimestampTreeMap) Incr(key string) int64 {
	return m.IncrBy(key, 1)
}

// IncrBy increments key by delta and returns the new value. Creates key=delta if absent.
func (m *TimestampTreeMap) IncrBy(key string, delta int64) int64 {
	m.mu.RLock()
	e := llrbGet(m.root, key)
	m.mu.RUnlock()
	if e != nil {
		return e.tsIncrBy(delta)
	}

	m.mu.Lock()
	e = llrbGet(m.root, key)
	if e != nil {
		m.mu.Unlock()
		return e.tsIncrBy(delta)
	}
	m.root, _ = llrbPut(m.root, key, delta, func(ex *tsEntry) {})
	m.root.color = llrbBlack
	newEntry := llrbGet(m.root, key)
	m.mu.Unlock()
	// The entry was just created with value=delta; return it via seqlock read
	// to ensure the caller observes the published value.
	val, _, _ := newEntry.tsLoad()
	return val
}

// Len returns the total number of keys in the map.
func (m *TimestampTreeMap) Len() int {
	m.mu.RLock()
	n := nodeSize(m.root)
	m.mu.RUnlock()
	return n
}

// Min returns the entry with the lexicographically smallest key.
// Returns (EntryInfo{}, false) if the map is empty.
func (m *TimestampTreeMap) Min() (EntryInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.root == nil {
		return EntryInfo{}, false
	}
	n := llrbMin(m.root)
	return makeEntryInfo(n.key, n.entry), true
}

// Max returns the entry with the lexicographically largest key.
func (m *TimestampTreeMap) Max() (EntryInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	if m.root == nil {
		return EntryInfo{}, false
	}
	n := llrbMax(m.root)
	return makeEntryInfo(n.key, n.entry), true
}

// Floor returns the entry whose key is the largest key ≤ target.
// Returns (EntryInfo{}, false) if no such key exists.
func (m *TimestampTreeMap) Floor(target string) (EntryInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := llrbFloor(m.root, target)
	if n == nil {
		return EntryInfo{}, false
	}
	return makeEntryInfo(n.key, n.entry), true
}

// Ceiling returns the entry whose key is the smallest key ≥ target.
func (m *TimestampTreeMap) Ceiling(target string) (EntryInfo, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	n := llrbCeiling(m.root, target)
	if n == nil {
		return EntryInfo{}, false
	}
	return makeEntryInfo(n.key, n.entry), true
}

// Range calls fn for every entry whose key is in [lo, hi] (inclusive),
// visiting in ascending key order. Returning false from fn stops iteration.
// The read lock is held for the duration of the traversal to guarantee a
// consistent snapshot.
func (m *TimestampTreeMap) Range(lo, hi string, fn func(EntryInfo) bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()
	llrbRange(m.root, lo, hi, fn)
}

// Keys returns all keys in ascending lexicographic order.
func (m *TimestampTreeMap) Keys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	out := make([]string, 0, nodeSize(m.root))
	llrbCollectKeys(m.root, &out)
	return out
}

// UpdatedAfter calls fn for every entry whose UpdatedAt is strictly after
// since, in ascending key order. Returning false from fn stops iteration.
func (m *TimestampTreeMap) UpdatedAfter(since time.Time, fn func(EntryInfo) bool) {
	sinceNano := since.UnixNano()
	m.mu.RLock()
	defer m.mu.RUnlock()
	llrbWalk(m.root, func(info EntryInfo) bool {
		if info.UpdatedAt.UnixNano() > sinceNano {
			return fn(info)
		}
		return true
	})
}
