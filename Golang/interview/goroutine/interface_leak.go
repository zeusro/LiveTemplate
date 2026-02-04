package goroutine

import (
	"fmt"
	"runtime"
)

// 错误示例：接口泄漏
// 问题：接口变量持有大对象的引用，即使只需要小部分数据

type LargeStruct struct {
	ID   int
	Data [1024 * 1024]byte // 1MB 数据
	Name string
}

// 错误示例：接口持有大结构体的引用
func interfaceLeakExample() {
	var obj interface{}
	
	// 创建大结构体
	largeObj := LargeStruct{
		ID:   1,
		Name: "test",
	}
	
	// 接口变量持有整个结构体的引用
	obj = largeObj
	
	// 即使只需要 Name 字段，但由于接口持有整个 largeObj
	// 整个 1MB 的数据无法被 GC 回收
	fmt.Println(obj.(LargeStruct).Name)
}

// 错误示例：接口 slice 持有大对象
func interfaceSliceLeak() {
	var objects []interface{}
	
	for i := 0; i < 100; i++ {
		obj := LargeStruct{
			ID:   i,
			Name: fmt.Sprintf("obj-%d", i),
		}
		objects = append(objects, obj)
		// 每个对象 1MB，100 个对象就是 100MB
		// 即使只需要 Name 字段，但所有数据都占用内存
	}
}

// 错误示例：接口 map 持有大对象
func interfaceMapLeak() {
	objects := make(map[string]interface{})
	
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		obj := LargeStruct{
			ID:   i,
			Name: fmt.Sprintf("obj-%d", i),
		}
		objects[key] = obj
		// 同样的问题：接口持有整个结构体
	}
}

// 正确示例：只存储需要的数据
type SmallStruct struct {
	ID   int
	Name string
}

func correctInterfaceExample() {
	var obj interface{}
	
	largeObj := LargeStruct{
		ID:   1,
		Name: "test",
	}
	
	// 只存储需要的小结构体
	smallObj := SmallStruct{
		ID:   largeObj.ID,
		Name: largeObj.Name,
	}
	obj = smallObj
	
	// 现在接口只持有小结构体，大结构体可以被 GC 回收
	fmt.Println(obj.(SmallStruct).Name)
}

// 正确示例：使用指针但及时清理
func correctInterfaceWithPointer() {
	var obj *LargeStruct
	
	largeObj := &LargeStruct{
		ID:   1,
		Name: "test",
	}
	obj = largeObj
	
	// 使用完后显式设置为 nil
	fmt.Println(obj.Name)
	obj = nil
	// 现在 largeObj 可以被 GC 回收
}

// 正确示例：使用接口但存储指针，并在不需要时清理
func correctInterfaceCleanup() {
	objects := make(map[string]interface{})
	
	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("key-%d", i)
		obj := &LargeStruct{
			ID:   i,
			Name: fmt.Sprintf("obj-%d", i),
		}
		objects[key] = obj
	}
	
	// 使用完后清理
	for key := range objects {
		delete(objects, key)
	}
	
	// 强制 GC（仅用于演示，生产环境不推荐频繁调用）
	runtime.GC()
}

// 正确示例：使用值类型但只包含必要字段
func correctMinimalInterface() {
	type MinimalData struct {
		ID   int
		Name string
	}
	
	objects := make([]MinimalData, 0, 100)
	
	for i := 0; i < 100; i++ {
		// 只存储必要的数据
		obj := MinimalData{
			ID:   i,
			Name: fmt.Sprintf("obj-%d", i),
		}
		objects = append(objects, obj)
		// 每个对象只有几十字节，而不是 1MB
	}
}

// RunInterfaceLeakExamples 演示本文件中的正确用法示例。
func RunInterfaceLeakExamples() {
	// 演示错误示例（注释掉以避免实际泄漏）
	// interfaceLeakExample()
	// interfaceSliceLeak()
	// interfaceMapLeak()
	
	// 正确示例
	correctInterfaceExample()
	correctInterfaceWithPointer()
	correctInterfaceCleanup()
	correctMinimalInterface()
}
