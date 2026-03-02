# TimestampTreeMap —— 带时间戳的并发安全有序 Map

## 题目描述

设计一种有序 Map，在支持 `Load / Store / Delete / Incr / IncrBy` 的同时：

- **键有序**：任意时刻可以查询最小键、最大键、区间扫描、前驱/后继
- **时间戳**：每个 key 记录**创建时间**（不可变）与**最后修改时间**（每次写更新）
- **时间过滤**：支持按修改时间筛选 key（`UpdatedAfter`）
- **并发安全**：所有操作在 `go test -race` 下无 data race

---

## 与 IncMap 的核心区别

| 维度 | IncMap（分片 Hash） | TimestampTreeMap（LLRB 树） |
|------|--------------------|-----------------------------|
| 底层结构 | 32 个 `map[string]*entry` 分片 | 单棵左倾红黑树（LLRB） |
| key 顺序 | 无序 | 升序（字典序） |
| 时间戳 | 无 | `CreatedAt`（不可变）+ `UpdatedAt`（每写更新） |
| 有序查询 | 不支持 | Min / Max / Floor / Ceiling / Range |
| 时间过滤 | 不支持 | `UpdatedAfter(since)` |
| 单次操作复杂度 | O(1) 均摊 | O(log n) |
| 并发写粒度 | 32 分片 × per-key seqlock | 单 RWMutex + per-entry seqlock |

---

## 数据结构

### tsEntry —— 带时间戳的 entry

```go
type tsEntry struct {
    version   atomic.Int64 // seqlock: 偶数=稳定，奇数=写中
    value     atomic.Int64
    createdAt int64        // unix nano，创建时写入，此后不变
    updatedAt atomic.Int64 // unix nano，每次 Store/IncrBy 更新
}
```

`createdAt` 是普通 `int64`（非 atomic），因为它在 `newTsEntry` 中写入后通过
树的写锁发布给其他 goroutine——写锁释放 happens-before 任何后续读锁获取，保证可见性。

`updatedAt` 必须是 atomic，因为它在 seqlock 临界区内被写、在 seqlock 读协议中被读。

### llrbNode —— 左倾红黑树节点

```go
type llrbNode struct {
    key   string
    entry *tsEntry
    color bool      // llrbRed / llrbBlack
    left  *llrbNode
    right *llrbNode
    sz    int       // 子树节点数（用于 Len）
}
```

### TimestampTreeMap

```go
type TimestampTreeMap struct {
    mu   sync.RWMutex // 保护树结构（插入/删除/旋转）
    root *llrbNode
}
```

---

## 并发模型：双层保护

```
┌────────────────────────────────────────────────────────────────┐
│  第 1 层：sync.RWMutex（树结构锁）                              │
│  写锁：仅用于插入新节点、删除节点、LLRB 旋转/重新着色           │
│  读锁：查找 entry 指针、遍历树（Range/Keys/UpdatedAfter）       │
├────────────────────────────────────────────────────────────────┤
│  第 2 层：per-entry seqlock（值 + updatedAt 一致性）            │
│  值更新（Store/IncrBy 到已存在的 key）只需读锁 + seqlock       │
│  不同 key 的值更新可完全并发，互不阻塞                         │
└────────────────────────────────────────────────────────────────┘
```

### 快慢路径（Store / IncrBy 为例）

```
1. RLock → 树内查找 entry → RUnlock
   ├── 找到（快路径）：直接走 seqlock 写协议，不阻塞其他读/写
   └── 未找到（慢路径）：
       a. WLock
       b. Double-check（防止两个 goroutine 同时进入慢路径）
       c. LLRB 插入 + 旋转 + 重新着色
       d. WUnlock
```

### 读操作（Load / Range / UpdatedAfter）

- 仅持读锁遍历树，不阻塞其他读者
- 找到 entry 后用 seqlock 读协议保证 value+updatedAt 一致
- `Range` / `UpdatedAfter` 在整个遍历期间持读锁，保证一致性快照

---

## 左倾红黑树（LLRB）

LLRB 是 Sedgewick（2008）简化版红黑树，满足以下不变量：

1. 根节点为黑色
2. 不存在连续两条红链接
3. 完美黑链接平衡（所有路径黑色节点数相同）
4. **红链接只能左倾**（左子节点为红 → 对应 2-3 树的 3-节点）

相比标准红黑树，LLRB 将旋转情况从 6 种降至 3 种：

| 情形 | 操作 |
|------|------|
| 右子为红、左子非红 | `rotateLeft` |
| 左子为红、左孙为红 | `rotateRight` |
| 左子为红、右子为红 | `flipColors`（分裂 4-节点） |

每次插入/删除最后都调用 `balance(h)` 统一处理上述三种情形。

### 插入（Put）

```
llrbPut(h, key, val):
    if h == nil → 返回新红节点
    cmp = compare(key, h.key)
    if cmp < 0 → 递归向左
    if cmp > 0 → 递归向右
    else       → 调用 updateFn 更新已有 entry（不创建新节点）
    return balance(h)
```

### 删除（Delete）

删除时从根向下保证"当前节点或其兄弟节点是红色"（即处于 3-节点/4-节点），
这样叶子节点才可以安全移除而不破坏黑平衡。

- 目标在左子树：`moveRedLeft` 后递归向左
- 目标在右子树或当前节点：`moveRedRight` 后处理，或用右子树最小节点替换后删除最小

---

## 时间戳语义

```go
type EntryInfo struct {
    Key       string
    Value     int64
    CreatedAt time.Time // 第一次 Store/IncrBy 时设置，之后永不变更
    UpdatedAt time.Time // 每次 Store/IncrBy 都刷新
}
```

```
初次插入：CreatedAt = UpdatedAt = now
再次 Store：CreatedAt 不变，UpdatedAt = now
IncrBy：    CreatedAt 不变，UpdatedAt = now
Delete 后重新 Store：视为全新 key，CreatedAt = now
```

### UpdatedAfter 使用场景

- **增量同步**：记录上次同步时间 `checkpoint`，调用 `UpdatedAfter(checkpoint, ...)` 只处理新变更
- **缓存失效**：找出一段时间内被修改的 key，批量推送更新
- **审计日志**：遍历指定时间段内的写操作

---

## API 对照

| 方法 | 说明 |
|------|------|
| `Load(key)` | 读取值，不存在返回 `(0, false)` |
| `LoadInfo(key)` | 读取值 + 两个时间戳 |
| `Store(key, val)` | 写入值，CreatedAt 首次设置后不变 |
| `Delete(key)` | 删除 key，不存在时 no-op |
| `Incr(key)` | 值 +1，返回新值 |
| `IncrBy(key, delta)` | 值 +delta，返回新值 |
| `Len()` | key 总数 |
| `Min()` | 字典序最小 key 的 EntryInfo |
| `Max()` | 字典序最大 key 的 EntryInfo |
| `Floor(target)` | 最大的 key ≤ target |
| `Ceiling(target)` | 最小的 key ≥ target |
| `Range(lo, hi, fn)` | 区间 [lo, hi] 升序遍历，fn 返回 false 提前终止 |
| `Keys()` | 返回全部 key 的升序切片 |
| `UpdatedAfter(since, fn)` | 遍历 UpdatedAt > since 的 entry |

---

## 常见追问

**Q：为什么不像 IncMap 那样分片？**

分片对哈希 Map 效果好，但树只有一个根，分片后无法跨片做全局有序遍历（Min/Max/Range
需要合并多棵树，复杂度上升）。本题的核心价值在于有序性，单树 + 读锁并发模型是合理权衡。

**Q：Range 持读锁期间会阻塞写操作吗？**

`Range` 持读锁，新插入/删除（需写锁）会等待；但对**已有 key 的值更新**（Store/IncrBy 快路径）
只需读锁 + seqlock，可与 Range 并发。

**Q：createdAt 为什么不用 atomic？**

`newTsEntry` 在持写锁期间调用，写锁的释放 happens-before 后续任何读锁获取，
因此 `createdAt` 对所有后续读者可见，无需额外的原子操作。

**Q：Delete 后重新 Incr，CreatedAt 是新的还是旧的？**

新的。Delete 移除整个节点，Incr 会触发 `llrbPut` 创建新节点，调用 `newTsEntry`，
`CreatedAt` 重置为当前时间。

**Q：如何扩展 TTL？**

在 `tsEntry` 中加 `expireAt atomic.Int64`（unix nano，0=永不过期）；
`tsLoad` 返回前检查 `expireAt`，若已过期返回 `(0, false)` 并在写锁下惰性删除节点。
