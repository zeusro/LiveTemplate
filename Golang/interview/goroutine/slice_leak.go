package main

import (
	"fmt"
)

// 错误示例：切片泄漏
// 问题：大切片持有小切片的引用，导致小切片无法被 GC 回收

// 错误示例：从大切片中切出小切片，但保留对大切片的引用
func sliceLeakExample() {
	// 创建一个包含 1000 个元素的大切片
	largeSlice := make([]int, 1000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 只使用前 10 个元素
	smallSlice := largeSlice[:10]
	
	// 即使 largeSlice 不再使用，但由于 smallSlice 持有对底层数组的引用
	// 整个 largeSlice 的底层数组（1000 个元素）无法被 GC 回收
	// 只有 10 个元素被使用，但 1000 个元素都占用内存
	
	fmt.Println(smallSlice)
}

// 错误示例：函数返回切片时泄漏
func getFirstN(n int) []int {
	largeSlice := make([]int, 10000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 返回前 n 个元素
	// 但返回的切片仍然引用整个底层数组
	return largeSlice[:n]
}

// 错误示例：多个小切片引用大切片
func multipleSliceLeak() {
	largeSlice := make([]byte, 1024*1024*100) // 100MB
	
	// 创建多个小切片
	slices := make([][]byte, 100)
	for i := 0; i < 100; i++ {
		start := i * 1024 * 1024
		end := start + 1024
		slices[i] = largeSlice[start:end] // 每个切片只有 1KB
	}
	
	// 即使只需要 100KB 的数据（100 * 1KB）
	// 但由于所有切片都引用 largeSlice 的底层数组
	// 整个 100MB 的数组无法被 GC 回收
}

// 正确示例：使用 copy 创建独立的切片
func correctSliceCopy() {
	largeSlice := make([]int, 1000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 使用 copy 创建独立的切片
	smallSlice := make([]int, 10)
	copy(smallSlice, largeSlice[:10])
	
	// 现在 smallSlice 是独立的，不引用 largeSlice 的底层数组
	// largeSlice 可以被 GC 回收
	fmt.Println(smallSlice)
}

// 正确示例：函数返回时使用 copy
func correctGetFirstN(n int) []int {
	largeSlice := make([]int, 10000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 创建新切片并复制数据
	result := make([]int, n)
	copy(result, largeSlice[:n])
	return result
	// 返回的切片不引用 largeSlice，largeSlice 可以被 GC 回收
}

// 正确示例：使用 append 创建新切片
func correctSliceAppend() {
	largeSlice := make([]int, 1000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 使用 append 创建新切片
	smallSlice := append([]int(nil), largeSlice[:10]...)
	
	// smallSlice 是新的切片，不引用 largeSlice
	fmt.Println(smallSlice)
}

// 正确示例：显式截断大切片
func correctSliceTruncate() {
	largeSlice := make([]int, 1000)
	for i := range largeSlice {
		largeSlice[i] = i
	}
	
	// 只保留需要的部分
	largeSlice = largeSlice[:10]
	// 现在 largeSlice 只引用前 10 个元素
	// 后面的元素可以被 GC 回收（如果底层数组被重新分配）
	
	// 更好的方式是创建新切片
	largeSlice = append([]int(nil), largeSlice[:10]...)
}

func main() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// sliceLeakExample()
	// result := getFirstN(10)
	// fmt.Println(result)
	// multipleSliceLeak()
	
	// 正确示例
	correctSliceCopy()
	correctGetFirstN(10)
	correctSliceAppend()
	correctSliceTruncate()
}
