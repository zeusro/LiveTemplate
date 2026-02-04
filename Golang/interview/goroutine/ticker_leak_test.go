package goroutine

import "testing"

// TestRunTickerLeakExamples 仅验证示例入口函数可以被调用。
// 注意：该示例内部包含 time.Sleep，运行时间略长。
func TestRunTickerLeakExamples(t *testing.T) {
	RunTickerLeakExamples()
}

