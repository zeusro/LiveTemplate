package goroutine

import "testing"

// TestSliceCorrectExamples 验证切片相关的正确用法示例。
func TestSliceCorrectExamples(t *testing.T) {
	correctSliceCopy()

	result := correctGetFirstN(10)
	if len(result) != 10 {
		t.Fatalf("expected length 10, got %d", len(result))
	}

	correctSliceAppend()
	correctSliceTruncate()
}

