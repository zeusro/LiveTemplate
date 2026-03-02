package incrmap

import (
	"sync"
	"testing"
)

func TestIncMap_BasicOperations(t *testing.T) {
	m := NewIncMap()

	// Store + Load
	m.Store("a", 10)
	v, ok := m.Load("a")
	if !ok || v != 10 {
		t.Fatalf("expected 10, got %d (ok=%v)", v, ok)
	}

	// Incr on existing key
	newVal := m.Incr("a")
	if newVal != 11 {
		t.Fatalf("expected 11 after Incr, got %d", newVal)
	}

	// Incr on non-existing key starts from 0
	newVal = m.Incr("new_key")
	if newVal != 1 {
		t.Fatalf("expected 1 for new key, got %d", newVal)
	}

	// IncrBy with negative delta
	newVal = m.IncrBy("a", -5)
	if newVal != 6 {
		t.Fatalf("expected 6 after IncrBy(-5), got %d", newVal)
	}

	// Delete
	m.Delete("a")
	_, ok = m.Load("a")
	if ok {
		t.Fatal("expected key 'a' to be deleted")
	}

	// Len
	if m.Len() != 1 {
		t.Fatalf("expected len=1, got %d", m.Len())
	}
}

func TestIncMap_LoadMissing(t *testing.T) {
	m := NewIncMap()
	_, ok := m.Load("ghost")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestIncMap_DeleteNonExisting(t *testing.T) {
	m := NewIncMap()
	m.Delete("not_there") // 不应 panic
}

// TestIncMap_ConcurrentIncr 验证高并发下 Incr 的原子性：
// N 个 goroutine 各执行 M 次 Incr，最终值应等于 N*M。
func TestIncMap_ConcurrentIncr(t *testing.T) {
	const (
		goroutines = 100
		incrTimes  = 1000
	)

	m := NewIncMap()
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < incrTimes; j++ {
				m.Incr("counter")
			}
		}()
	}

	wg.Wait()

	want := int64(goroutines * incrTimes)
	got, _ := m.Load("counter")
	if got != want {
		t.Fatalf("concurrent Incr: expected %d, got %d", want, got)
	}
}

// TestIncMap_ConcurrentReadWrite 同时进行读写删，不应发生 data race。
// 配合 go test -race 运行。
func TestIncMap_ConcurrentReadWrite(t *testing.T) {
	m := NewIncMap()
	m.Store("x", 0)

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(3)
		go func() {
			defer wg.Done()
			m.Incr("x")
		}()
		go func() {
			defer wg.Done()
			m.Load("x")
		}()
		go func() {
			defer wg.Done()
			m.Store("x", 0)
		}()
	}
	wg.Wait()
}

// TestIncMap_MultiKeyParallel 多个不同 key 并发操作，验证分片隔离。
func TestIncMap_MultiKeyParallel(t *testing.T) {
	const keys = 20
	const incrTimes = 500

	m := NewIncMap()
	var wg sync.WaitGroup
	wg.Add(keys)

	for i := 0; i < keys; i++ {
		go func(idx int) {
			defer wg.Done()
			key := string(rune('A' + idx))
			for j := 0; j < incrTimes; j++ {
				m.Incr(key)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < keys; i++ {
		key := string(rune('A' + i))
		got, ok := m.Load(key)
		if !ok || got != incrTimes {
			t.Errorf("key %s: expected %d, got %d (ok=%v)", key, incrTimes, got, ok)
		}
	}
}
