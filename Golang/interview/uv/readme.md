# 使用 HyperLogLog 统计 UV

本目录提供基于 **HyperLogLog（HLL）** 算法的 **UV（独立访客）** 基数估计实现，适用于海量访问场景下用极少内存估算不重复用户数。

## 什么是 HyperLogLog？

HyperLogLog 是一种**概率型数据结构**，用于估计一个集合的**基数**（不重复元素个数），特点：

- **内存极小**：标准配置约 **16KB**（精度 p=14，2^14 个寄存器，每寄存器 1 字节），即可估计数十亿级 UV
- **误差可控**：标准误差约 **0.81%**（与精度参数有关）
- **可合并**：多个 HLL 可合并（如多机、多分片）得到全局 UV 估计
- **仅支持“加元素”与“估基数”**：不支持删除单条、不支持列举元素

典型用途：页面 UV、独立 IP 数、去重计数等，在 Redis、大数据系统中广泛使用。

## 算法简述

1. 对每个元素做 **64 位哈希**。
2. 取哈希的**前 p 位**作为寄存器下标（共 2^p 个寄存器）。
3. 在**剩余位**中计算“从左到右第一个 1 的位置”ρ，将该寄存器更新为**历史最大值**。
4. 基数估计公式：  
   \( E = \alpha_m \cdot m^2 / \sum_j 2^{-M[j]} \)  
   其中 \(m = 2^p\)，\(M[j]\) 为第 j 个寄存器的值，\(\alpha_m\) 为无偏常数。
5. 对小基数、大基数有**修正公式**以降低偏差。

## 目录结构

```
uv/
├── hyperloglog.go   # HLL 核心实现（Add / Cardinality / Merge）
├── demo.go          # 演示：模拟访问流，对比真实 UV 与估计值
├── cmd/
│   └── main.go      # 可执行入口，运行 Demo
└── readme.md        # 本文档
```

## API 使用

### 创建估计器

```go
import "interview/uv"

// 默认精度 p=14，约 16KB，标准误差 ~0.81%
h := uv.NewDefault()

// 或指定精度 p（4–20），p 越大内存与精度越高
h := uv.New(16)  // 2^16 个寄存器
```

### 添加访客

```go
// 按用户 ID（uint64）
h.AddUint64(userID)

// 按字符串（如 IP、DeviceID）
h.AddString("192.168.1.1")

// 按任意字节
h.Add([]byte(visitorID))
```

重复添加同一元素不会重复计数，符合 UV 语义。

### 读取估计值

```go
estimatedUV := h.Cardinality()
```

### 合并多分片（可选）

多机或多时间段各自维护一个 HLL，合并得到全局 UV：

```go
h1 := uv.NewDefault()
h2 := uv.NewDefault()
// ... 分别 Add ...
h1.Merge(h2)
globalUV := h1.Cardinality()
```

### 清空

```go
h.Reset()
```

## 运行示例

在 `Golang/interview` 所在模块根目录执行：

```bash
go run ./uv/cmd
```

将输出类似：

```
真实 UV（去重后）: xxxx
HLL 估计 UV:       xxxx
误差:              x.xx%
寄存器数量:        16384 (约 16.0 KB)
```

用于直观感受在约 10 万次访问、约 5000 独立访客下的估计效果。

## 精度与内存

| 精度 p | 寄存器数 m | 内存约 | 标准误差约 |
|--------|-------------|--------|------------|
| 12     | 4096        | ~3 KB  | ~1.6%      |
| 14     | 16384       | ~16 KB | ~0.81%     |
| 16     | 65536       | ~48 KB | ~0.4%      |

本实现中每个寄存器 1 字节（可存 0–63），故内存约 \(m\) 字节；p=14 时约 16 KB。

## 参考

- Flajolet, P., et al. “HyperLogLog: the analysis of a near-optimal cardinality estimation algorithm.” *AofA* (2007).
- [Redis PFCOUNT / PFADD](https://redis.io/commands/pfadd/) 使用 HLL 做基数统计。
