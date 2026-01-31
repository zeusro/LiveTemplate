// Package uv 提供 UV 统计示例：使用 HyperLogLog 估计独立访客数。
package uv

import (
	"fmt"
	"math/rand"
)

// Demo 演示：模拟大量访问（含重复），用 HLL 估计 UV 并与真实集合对比
func Demo() {
	h := NewDefault()
	// 模拟 10 万次“访问”，独立访客约 1 万（ID 从 0 到 approxUnique*2-1 随机重复）
	unique := make(map[uint64]struct{})
	const totalVisits = 100_000
	const approxUnique = 5_000
	rand.Seed(42)
	for i := 0; i < totalVisits; i++ {
		visitorID := uint64(rand.Intn(approxUnique * 2))
		unique[visitorID] = struct{}{}
		h.AddUint64(visitorID)
	}
	realUV := len(unique)
	estUV := h.Cardinality()
	fmt.Printf("真实 UV（去重后）: %d\n", realUV)
	fmt.Printf("HLL 估计 UV:       %d\n", estUV)
	fmt.Printf("误差:              %.2f%%\n", 100*float64(int64(estUV)-int64(realUV))/float64(realUV))
	fmt.Printf("寄存器数量:        %d (约 %.1f KB)\n", h.Registers(), float64(h.Registers())/1024)
}

// ExampleUV 可在测试或 main 中调用的示例
func ExampleUV() {
	h := New(14)
	// 模拟按 IP 或 UserID 统计
	h.AddString("192.168.1.1")
	h.AddString("192.168.1.2")
	h.AddString("192.168.1.1") // 重复
	h.AddString("10.0.0.1")
	fmt.Println(h.Cardinality()) // 约 3
}
