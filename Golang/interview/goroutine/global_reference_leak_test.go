package goroutine

import "testing"

// TestGlobalReferenceCorrectExamples 验证全局引用相关的正确用法示例。
func TestGlobalReferenceCorrectExamples(t *testing.T) {
	correctLocalVariable()
	correctGlobalCacheWithCleanup()
	correctLimitedCache()
	correctExplicitCleanup()
}

