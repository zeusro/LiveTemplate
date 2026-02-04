package goroutine

import "testing"

// TestInterfaceCorrectExamples 验证接口相关的正确用法示例可以正常运行。
func TestInterfaceCorrectExamples(t *testing.T) {
	correctInterfaceExample()
	correctInterfaceWithPointer()
	correctInterfaceCleanup()
	correctMinimalInterface()
}

