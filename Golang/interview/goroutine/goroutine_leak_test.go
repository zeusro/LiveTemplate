package goroutine

import "testing"

// TestRunGoroutineLeakExamples 验证 goroutine 示例入口函数可以被调用。
// 注意：内部包含 time.Sleep，运行时间略长。
func TestRunGoroutineLeakExamples(t *testing.T) {
	RunGoroutineLeakExamples()
}

