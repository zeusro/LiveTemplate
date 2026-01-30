# Go 全栈开发工程师线上测评 - 参考答案

> 题目来源：[Go全栈开发工程师（Go + Vue）线上测评](https://woshihaiwaichat.com/recruitment-evaluation.html)

---

## Part 1: 思维与架构设计

### 题目 1：受限环境下的高并发处理

**场景**：Go 服务通过 HTTP 接收海量请求（Payload 为 URL），抓取 URL 内容并分析入库。上游峰值 10k QPS，不允许引入 Redis/Kafka/RabbitMQ，单机进程内存消化流量，允许极端情况下丢弃任务，但不能 OOM 或 Crash。

#### 1. 优雅停机 (Graceful Shutdown)，减少排队任务丢失

- **监听退出信号**：在 main 中监听 `SIGTERM`/`SIGINT`，收到后不再接受新请求。
- **停止接收新任务**：关闭 HTTP Server 的 `Listen`（调用 `Server.Shutdown(ctx)`），新请求直接返回 503 或拒绝。
- **等待存量任务**：使用带超时的 `context`（如 30s）调用 `Server.Shutdown(ctx)`，等待当前正在处理的请求完成。
- **队列排空策略**：
  - 用 **done channel** 或 **context** 通知 worker 池“不再接受新任务”。
  - Worker 只消费已有队列，不退出直到队列为空或超时。
  - 若 Shutdown 超时时间到，仍有未处理任务，可记录到本地文件（如 JSONL）做“冷恢复”，下次启动时优先重放；若不允许持久化，则丢弃并打日志。
- **关键点**：先停入口（Shutdown），再等 worker 排空队列，最后退出进程；超时时间根据业务可接受延迟设定。

#### 2. 进程内“缓冲队列”的数据结构选择

- **推荐：有界、固定大小的带缓冲 Channel（或基于 RingBuffer 的队列）**
  - **原因**：
    - 有界可防止无限堆积导致 OOM：队列满时直接丢弃新任务（或返回 503），内存上界 = 队列容量 × 单任务元数据大小。
    - Channel 语义简单，多个 worker 可直接 `for range ch` 消费，无需手写锁。
  - **不选无缓冲 Channel**：无缓冲时生产者必须等消费者就绪，高 QPS 下大量 goroutine 阻塞，易导致 goroutine 暴涨、内存与调度压力大。
  - **不选“无限 Slice 追加”**：无界队列在高 QPS 下会无限增长，容易 OOM。
  - **RingBuffer 备选**：若需要“固定内存、覆盖最旧任务”的语义，可用 RingBuffer（固定大小数组 + 头尾指针），满时覆盖最旧任务并计数丢弃，实现有界、O(1) 入队出队。
  - **建议**：队列容量根据单机内存和单任务大小计算（如 10k 任务 × 1KB ≈ 10MB），并配合 worker 数量（如 CPU 核数 × 2）做消费能力与延迟的平衡。

---

### 题目 2：Vue 大数据量实时渲染

**场景**：后台监控面板，WebSocket 每秒约 500 条日志，前端保留最近 10,000 条供回滚查看。

#### 1. 数据层面优化（减少响应式开销）

- **不深度监听大数组**：  
  - 用 `Object.freeze()` 冻结单条日志对象，或使用 `shallowRef` / `markRaw` 存列表，让 Vue 不对数组元素做深度响应式。
  - 若用 `ref`，存“普通对象数组”且不需要对每条日志字段做响应式时，可在接入数据时 `Object.freeze(log)` 或 `markRaw(log)`，避免每条日志都变成响应式。
- **用 `shallowRef` 存列表**：  
  - `const logs = shallowRef<Log[]>([])`，只对 `logs` 本身做响应式，替换整个数组时触发更新，而不对数组内部 10,000 个元素做依赖收集，大幅降低开销。
- **增量更新**：  
  - 不每次 `logs.value = [...newList]` 全量替换，而是只 `push` 新条目并必要时 `splice` 掉超出 10,000 的旧条目，配合下面虚拟列表只渲染可见区域，减少无效的依赖追踪。
- **必要时脱离响应式**：  
  - 若某块逻辑只做“读一遍列表再计算”，可用 `toRaw(logs.value)` 或单独维护一份非响应式副本，减少响应式系统参与。

#### 2. 渲染技术方案（避免 10,000 条 DOM 卡死）

- **虚拟滚动 (Virtual Scroll)**  
  - 思路：只渲染当前视口内可见的若干条（例如几十条），列表容器有固定高度并设 `overflow: auto`，根据滚动位置计算“可见区间”的起止索引，只对这段区间内的数据渲染 DOM，其余用占位 div 撑开高度。
  - 实现要点：
    - 单条高度固定（或可测）：用 `itemHeight * length` 得到总高度；若不定高，可用“预估高度 + 测量”或按区间缓存高度。
    - 滚动时根据 `scrollTop` 计算：`startIndex = Math.floor(scrollTop / itemHeight)`，`endIndex = startIndex + visibleCount`，`visibleCount = Math.ceil(containerHeight / itemHeight) + buffer`（可加少量 buffer 减少白屏）。
    - 可见区数据：`list.slice(startIndex, endIndex)`；列表内容区用 `transform: translateY(startIndex * itemHeight)` 或顶部 padding 保持滚动位置与条数一致。
  - 效果：DOM 数量恒定在几十个，而不是 10,000 个，滚动和更新都流畅。
- **可选库**：`vue-virtual-scroller`、`@vueuse/core` 的 `useVirtualList` 等，本质都是上述“按可见区间切片渲染”的实现。

---

## Part 2: 代码审计与找茬

题目中的 Go 代码存在以下**至少 3 处**严重问题：

### 问题 1：并发读写全局 map（Concurrent Map Access）

- **代码**：`stats[deviceID]++` 与 `stats[deviceID]` 的读发生在多个 goroutine 中，且 `stats` 为全局 `map[string]int`，无任何同步。
- **后果**：Go 的 map 非并发安全，并发写或“写与读”并发会直接 **panic: concurrent map writes** 或 **concurrent map read and map write**，导致进程崩溃。
- **修复**：使用 `sync.Mutex` 保护整个 map，或使用 `sync.Map`；读增删都在锁内或通过 `sync.Map` 的 API 进行。

### 问题 2：Goroutine 内使用 `time.After` 导致资源泄漏

- **代码**：在 `select` 中使用 `case <-time.After(2 * time.Second)`。当 `ctx.Done()` 先触发（如客户端提前断开）时，`time.After` 创建的定时器不会被回收，要等到 2 秒后才会被 GC 回收。
- **后果**：高并发下大量请求断开时，会堆积大量未到期的 timer，导致内存和调度开销增加，属于典型的“timer 泄漏”。
- **修复**：使用 `time.NewTimer(2*time.Second)`，在 `select` 外 `defer timer.Stop()`，或在 `select` 中收到 `ctx.Done()` 时调用 `timer.Stop()`，确保 timer 被及时释放。

### 问题 3：Request 的 Context 与 Goroutine 生命周期

- **代码**：在 `go func()` 中使用 `ctx := r.Context()`。Handler 返回后，`r` 可能被 Server 回收或复用，而 goroutine 仍在运行。
- **后果**：虽然 `r.Context()` 在 handler 内已经拿到，在 goroutine 里只用这个 context 做取消是安全的（客户端断开时会取消），但若 goroutine 内再访问 `r` 的其他字段（如 `r.URL`、`r.Body`），就会在 handler 返回后访问已无效的 Request，导致数据竞争或未定义行为。当前代码只用了 `deviceID`（在启动 goroutine 前已取）和 `r.Context()`，所以“只取 context”这部分是可接受的，但最佳实践是**在启动 goroutine 前**把需要的参数（如 `deviceID`）和 `r.Context()` 都拷贝到局部变量，避免在 goroutine 里再引用 `r`。
- **修复**：在 `go func()` 外先执行 `ctx := r.Context(); id := deviceID`，goroutine 内只使用 `ctx` 和 `id`，不再使用 `r`。

### 问题 4：缺少优雅停机与生产级 HTTP 配置

- **代码**：直接 `http.ListenAndServe(":8080", nil)`，无超时、无优雅退出。
- **后果**：进程被 kill 时连接被强制断开，无法排空正在处理的请求；且未设置 `ReadHeaderTimeout`、`ReadTimeout`、`WriteTimeout` 等，易受慢请求或恶意连接影响。
- **修复**：使用 `http.Server` 结构体，配置 `ReadHeaderTimeout`、`ReadTimeout`、`WriteTimeout`、`IdleTimeout`，并在 main 中监听 `SIGTERM`/`SIGINT`，调用 `Server.Shutdown(context.WithTimeout(...))` 实现优雅停机。

---

以上 Part 1、Part 2 的参考答案可与仓库内以下代码对照使用：

- Part 2 错误示例与修复实现：见 `part2-audit/bad/main.go`（题目原始代码）、`part2-audit/fixed/main.go`（修复版）。
- Part 3 简易服务器状态监控看板：见 `monitor-dashboard/` 目录（Go SSE 后端 + Vue 3 + TS 前端）。

---

## Part 3: AI 协作实战（Mini Project）实现说明

- **后端**（`monitor-dashboard/backend/main.go`）：  
  - `/api/metrics` 返回单次随机 CPU/内存；  
  - `/api/stream` 为 SSE，由全局 ticker 按 `intervalSec` 向所有连接广播；  
  - `/api/interval` 支持 query `seconds` 或 JSON body `{"seconds": n}` 修改间隔（1–60 秒）。  
- **前端**（`monitor-dashboard/frontend/`）：  
  - Vue 3 + TypeScript，使用 ECharts 绘制 CPU/内存双曲线；  
  - 使用 `EventSource('/api/stream')` 接收 SSE；  
  - 下拉框修改间隔后调用 `PUT /api/interval?seconds=x`，后端立即生效，下一轮 tick 即按新间隔推送。  
- 运行方式与 AI 使用报告模板见 `monitor-dashboard/README.md`。
