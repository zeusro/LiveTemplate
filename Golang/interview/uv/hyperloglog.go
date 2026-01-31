// Package uv 提供基于 HyperLogLog 的 UV（独立访客）基数估计。
// 使用极少内存（约 16KB，p=14）估计海量数据中的不重复元素个数，标准误差约 0.81%。
package uv

import (
	"hash"
	"hash/fnv"
	"math"
)

const (
	// DefaultP 默认精度参数，寄存器数量 m = 2^14 = 16384，约 16KB
	DefaultP = 14
	// MaxP 最大精度，寄存器数量 2^20
	MaxP = 20
)

// HyperLogLog 基于 HyperLogLog 算法的基数估计器（Flajolet et al.）
type HyperLogLog struct {
	p   uint8   // 精度，寄存器索引占 p 位
	m   uint32  // 寄存器数量 2^p
	reg []uint8 // 每个寄存器存 ρ(w)，最多 6 位
}

// New 创建一个精度为 p 的 HyperLogLog，p 建议 14（约 16KB），范围 [4, 20]
func New(p uint8) *HyperLogLog {
	if p < 4 {
		p = 4
	}
	if p > MaxP {
		p = MaxP
	}
	m := uint32(1) << p
	return &HyperLogLog{
		p:   p,
		m:   m,
		reg: make([]uint8, m),
	}
}

// NewDefault 使用默认精度 p=14 创建
func NewDefault() *HyperLogLog {
	return New(DefaultP)
}

// Add 将元素加入集合（重复加入不会重复计数）
func (h *HyperLogLog) Add(data []byte) {
	x := h.hash64(data)
	h.addHash(x)
}

// AddString 将字符串加入集合
func (h *HyperLogLog) AddString(s string) {
	h.Add([]byte(s))
}

// AddUint64 将 uint64 加入集合（如用户 ID）；内部会先做 64 位混合再参与估计
func (h *HyperLogLog) AddUint64(v uint64) {
	h.addHash(mix64(v))
}

func (h *HyperLogLog) addHash(x uint64) {
	// 前 p 位作为寄存器下标
	j := uint32(x >> (64 - h.p))
	// 剩余位中从左到右第一个 1 的位置（1-based）
	w := x << h.p
	rho := rho(w, 64-int(h.p))
	if rho > h.reg[j] {
		h.reg[j] = rho
	}
}

// rho 返回 w 中最高位起第一个 1 的位置（1-based），若全 0 返回 maxBits+1
func rho(w uint64, maxBits int) uint8 {
	if w == 0 {
		return uint8(maxBits + 1)
	}
	pos := 0
	for (w & (1 << 63)) == 0 && pos < maxBits {
		w <<= 1
		pos++
	}
	return uint8(pos + 1)
}

// Cardinality 返回当前估计的基数（不重复元素个数）
func (h *HyperLogLog) Cardinality() uint64 {
	sum := 0.0
	zeros := 0
	for i := uint32(0); i < h.m; i++ {
		if h.reg[i] == 0 {
			zeros++
		}
		sum += math.Pow(2, -float64(h.reg[i]))
	}
	alpha := alphaM(h.m)
	est := alpha * float64(h.m*h.m) / sum

	// 小范围修正（线性计数）
	if est <= 2.5*float64(h.m) && zeros > 0 {
		est = float64(h.m) * math.Log(float64(h.m)/float64(zeros))
	}
	// 大范围修正
	if est > (1<<32)/30.0 {
		est = -math.Pow(2, 32) * math.Log(1-est/math.Pow(2, 32))
	}
	return uint64(est + 0.5)
}

// alpha_m ≈ 0.7213 / (1 + 1.079/m)，用于无偏估计
func alphaM(m uint32) float64 {
	switch m {
	case 16:
		return 0.673
	case 32:
		return 0.697
	case 64:
		return 0.709
	default:
		return 0.7213 / (1 + 1.079/float64(m))
	}
}

// Merge 将另一个 HLL 合并到当前（用于多机/多分片合并 UV）
func (h *HyperLogLog) Merge(other *HyperLogLog) {
	if other == nil || h.p != other.p {
		return
	}
	for i := uint32(0); i < h.m; i++ {
		if other.reg[i] > h.reg[i] {
			h.reg[i] = other.reg[i]
		}
	}
}

// Reset 清空所有寄存器
func (h *HyperLogLog) Reset() {
	for i := range h.reg {
		h.reg[i] = 0
	}
}

// hash64 将 data 映射为 64 位，用于 HLL 内部
func (h *HyperLogLog) hash64(data []byte) uint64 {
	var hasher hash.Hash64 = fnv.New64a()
	_, _ = hasher.Write(data)
	return hasher.Sum64()
}

// mix64 对 64 位整数做混合，使高位分布均匀（避免连续 ID 落入同一寄存器）
func mix64(x uint64) uint64 {
	x ^= x >> 33
	x *= 0xff51afd7ed558ccd
	x ^= x >> 33
	x *= 0xc4ceb9fe1a85ec53
	x ^= x >> 33
	return x
}

// Registers 返回寄存器数量（用于文档/测试）
func (h *HyperLogLog) Registers() uint32 { return h.m }

// Precision 返回精度 p
func (h *HyperLogLog) Precision() uint8 { return h.p }
