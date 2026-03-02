package incrmap

import (
	"fmt"
	"sync"
	"testing"
	"time"
)

// ─── Basic CRUD ───────────────────────────────────────────────────────────────

func TestTTM_StoreLoad(t *testing.T) {
	m := NewTimestampTreeMap()

	m.Store("b", 20)
	m.Store("a", 10)
	m.Store("c", 30)

	cases := []struct {
		key  string
		want int64
	}{{"a", 10}, {"b", 20}, {"c", 30}}
	for _, tc := range cases {
		got, ok := m.Load(tc.key)
		if !ok || got != tc.want {
			t.Errorf("Load(%q) = (%d, %v), want (%d, true)", tc.key, got, ok, tc.want)
		}
	}
}

func TestTTM_LoadMissing(t *testing.T) {
	m := NewTimestampTreeMap()
	_, ok := m.Load("ghost")
	if ok {
		t.Fatal("expected ok=false for missing key")
	}
}

func TestTTM_Delete(t *testing.T) {
	m := NewTimestampTreeMap()
	m.Store("x", 1)
	m.Delete("x")
	_, ok := m.Load("x")
	if ok {
		t.Fatal("key should be gone after Delete")
	}
}

func TestTTM_DeleteNonExisting(t *testing.T) {
	m := NewTimestampTreeMap()
	m.Delete("no_such_key") // must not panic
}

func TestTTM_Incr(t *testing.T) {
	m := NewTimestampTreeMap()

	// New key starts at delta=1.
	if v := m.Incr("k"); v != 1 {
		t.Fatalf("first Incr want 1, got %d", v)
	}
	if v := m.Incr("k"); v != 2 {
		t.Fatalf("second Incr want 2, got %d", v)
	}
	if v := m.IncrBy("k", -3); v != -1 {
		t.Fatalf("IncrBy(-3) want -1, got %d", v)
	}
}

func TestTTM_Len(t *testing.T) {
	m := NewTimestampTreeMap()
	for i := 0; i < 5; i++ {
		m.Store(fmt.Sprintf("k%d", i), int64(i))
	}
	if n := m.Len(); n != 5 {
		t.Fatalf("Len want 5, got %d", n)
	}
	m.Delete("k2")
	if n := m.Len(); n != 4 {
		t.Fatalf("Len after Delete want 4, got %d", n)
	}
}

// ─── Ordering ─────────────────────────────────────────────────────────────────

func TestTTM_Keys_Sorted(t *testing.T) {
	m := NewTimestampTreeMap()
	for _, k := range []string{"dog", "cat", "ant", "bee", "elk"} {
		m.Store(k, 1)
	}
	keys := m.Keys()
	want := []string{"ant", "bee", "cat", "dog", "elk"}
	if len(keys) != len(want) {
		t.Fatalf("Keys len %d, want %d", len(keys), len(want))
	}
	for i, k := range keys {
		if k != want[i] {
			t.Errorf("Keys[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestTTM_MinMax(t *testing.T) {
	m := NewTimestampTreeMap()

	_, ok := m.Min()
	if ok {
		t.Fatal("Min on empty map should return false")
	}

	for _, k := range []string{"mango", "apple", "zucchini", "banana"} {
		m.Store(k, 1)
	}

	min, ok := m.Min()
	if !ok || min.Key != "apple" {
		t.Errorf("Min want apple, got %q (ok=%v)", min.Key, ok)
	}
	max, ok := m.Max()
	if !ok || max.Key != "zucchini" {
		t.Errorf("Max want zucchini, got %q (ok=%v)", max.Key, ok)
	}
}

func TestTTM_FloorCeiling(t *testing.T) {
	m := NewTimestampTreeMap()
	for _, k := range []string{"b", "d", "f", "h"} {
		m.Store(k, 1)
	}

	// Floor: largest key ≤ target
	cases := []struct{ target, want string }{
		{"c", "b"},
		{"d", "d"},
		{"e", "d"},
		{"z", "h"},
	}
	for _, tc := range cases {
		info, ok := m.Floor(tc.target)
		if !ok || info.Key != tc.want {
			t.Errorf("Floor(%q) = %q (ok=%v), want %q", tc.target, info.Key, ok, tc.want)
		}
	}
	_, ok := m.Floor("a") // smaller than all keys
	if ok {
		t.Error("Floor of smallest possible should return false")
	}

	// Ceiling: smallest key ≥ target
	ceilCases := []struct{ target, want string }{
		{"a", "b"},
		{"b", "b"},
		{"c", "d"},
		{"g", "h"},
	}
	for _, tc := range ceilCases {
		info, ok := m.Ceiling(tc.target)
		if !ok || info.Key != tc.want {
			t.Errorf("Ceiling(%q) = %q (ok=%v), want %q", tc.target, info.Key, ok, tc.want)
		}
	}
	_, ok = m.Ceiling("z") // larger than all keys
	if ok {
		t.Error("Ceiling of largest possible should return false")
	}
}

func TestTTM_Range(t *testing.T) {
	m := NewTimestampTreeMap()
	for i, k := range []string{"a", "b", "c", "d", "e"} {
		m.Store(k, int64(i+1))
	}

	var got []string
	m.Range("b", "d", func(info EntryInfo) bool {
		got = append(got, info.Key)
		return true
	})
	want := []string{"b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("Range got %v, want %v", got, want)
	}
	for i, k := range got {
		if k != want[i] {
			t.Errorf("Range[%d] = %q, want %q", i, k, want[i])
		}
	}
}

func TestTTM_Range_EarlyStop(t *testing.T) {
	m := NewTimestampTreeMap()
	for _, k := range []string{"a", "b", "c", "d", "e"} {
		m.Store(k, 1)
	}
	var count int
	m.Range("a", "e", func(_ EntryInfo) bool {
		count++
		return count < 3 // stop after 3rd
	})
	if count != 3 {
		t.Fatalf("early stop: expected 3 calls, got %d", count)
	}
}

// ─── Timestamps ───────────────────────────────────────────────────────────────

func TestTTM_Timestamps_CreatedAtPreserved(t *testing.T) {
	m := NewTimestampTreeMap()

	before := time.Now()
	m.Store("key", 1)
	after := time.Now()

	info, ok := m.LoadInfo("key")
	if !ok {
		t.Fatal("LoadInfo returned false")
	}
	if info.CreatedAt.Before(before) || info.CreatedAt.After(after) {
		t.Errorf("CreatedAt %v not in [%v, %v]", info.CreatedAt, before, after)
	}

	// Store again: CreatedAt must not change.
	created := info.CreatedAt
	time.Sleep(2 * time.Millisecond)
	m.Store("key", 99)

	info2, _ := m.LoadInfo("key")
	if !info2.CreatedAt.Equal(created) {
		t.Errorf("CreatedAt changed after second Store: was %v, now %v", created, info2.CreatedAt)
	}
	if !info2.UpdatedAt.After(created) {
		t.Errorf("UpdatedAt %v should be after CreatedAt %v", info2.UpdatedAt, created)
	}
}

func TestTTM_UpdatedAfter(t *testing.T) {
	m := NewTimestampTreeMap()
	m.Store("old", 1)
	m.Store("also_old", 2)

	checkpoint := time.Now()
	time.Sleep(2 * time.Millisecond)

	m.Store("new1", 10)
	m.Store("new2", 20)
	m.IncrBy("old", 5) // refreshes UpdatedAt

	var updated []string
	m.UpdatedAfter(checkpoint, func(info EntryInfo) bool {
		updated = append(updated, info.Key)
		return true
	})

	// Expect "also_old" to NOT appear (not touched after checkpoint).
	// "old", "new1", "new2" should appear.
	want := map[string]bool{"old": true, "new1": true, "new2": true}
	for _, k := range updated {
		if !want[k] {
			t.Errorf("unexpected key in UpdatedAfter result: %q", k)
		}
		delete(want, k)
	}
	if len(want) > 0 {
		t.Errorf("missing keys in UpdatedAfter result: %v", want)
	}
}

// ─── Concurrency ──────────────────────────────────────────────────────────────

// TestTTM_ConcurrentIncr verifies that N*M increments on the same key
// yield exactly N*M. Run with: go test -race ./...
func TestTTM_ConcurrentIncr(t *testing.T) {
	const (
		goroutines = 100
		times      = 1000
	)
	m := NewTimestampTreeMap()
	var wg sync.WaitGroup
	wg.Add(goroutines)
	for i := 0; i < goroutines; i++ {
		go func() {
			defer wg.Done()
			for j := 0; j < times; j++ {
				m.Incr("counter")
			}
		}()
	}
	wg.Wait()

	got, _ := m.Load("counter")
	if want := int64(goroutines * times); got != want {
		t.Fatalf("ConcurrentIncr: got %d, want %d", got, want)
	}
}

// TestTTM_ConcurrentMultiKey concurrently increments distinct keys;
// no key should interfere with another.
func TestTTM_ConcurrentMultiKey(t *testing.T) {
	const (
		keyCount = 20
		times    = 500
	)
	m := NewTimestampTreeMap()
	var wg sync.WaitGroup
	wg.Add(keyCount)
	for i := 0; i < keyCount; i++ {
		go func(idx int) {
			defer wg.Done()
			key := fmt.Sprintf("key_%02d", idx)
			for j := 0; j < times; j++ {
				m.Incr(key)
			}
		}(i)
	}
	wg.Wait()

	for i := 0; i < keyCount; i++ {
		key := fmt.Sprintf("key_%02d", i)
		got, ok := m.Load(key)
		if !ok || got != times {
			t.Errorf("%s: got %d (ok=%v), want %d", key, got, ok, times)
		}
	}
}

// TestTTM_ConcurrentReadWriteDelete mixes Store, Load, Delete, Incr on the
// same key to exercise all concurrency paths. Run with -race.
func TestTTM_ConcurrentReadWriteDelete(t *testing.T) {
	m := NewTimestampTreeMap()
	m.Store("shared", 0)

	var wg sync.WaitGroup
	for i := 0; i < 60; i++ {
		wg.Add(4)
		go func() { defer wg.Done(); m.Incr("shared") }()
		go func() { defer wg.Done(); m.Load("shared") }()
		go func() { defer wg.Done(); m.Store("shared", 0) }()
		go func() { defer wg.Done(); m.Delete("shared") }()
	}
	wg.Wait()
}

// TestTTM_ConcurrentRange verifies that concurrent Range + Incr does not
// trigger the race detector.
func TestTTM_ConcurrentRange(t *testing.T) {
	m := NewTimestampTreeMap()
	for i := 0; i < 10; i++ {
		m.Store(fmt.Sprintf("%02d", i), int64(i))
	}

	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			m.Range("00", "09", func(_ EntryInfo) bool { return true })
		}()
		go func() {
			defer wg.Done()
			m.Incr("05")
		}()
	}
	wg.Wait()
}
