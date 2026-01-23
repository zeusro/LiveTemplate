# ASR 服务

基于 go-zero 框架实现的自动语音识别（ASR）服务，支持 WebSocket 和 gRPC 两种协议，提供完整的会话管理、流控、超时处理和回调功能。

## 功能特性

### 核心功能
- ✅ **ASR 服务封装**：支持 WebSocket 和 gRPC 两种协议
- ✅ **会话管理**：每通电话一个独立的 session，使用 Redis 存储
- ✅ **流控机制**：QPS 限制和并发会话数限制
- ✅ **超时处理**：会话超时自动清理
- ✅ **异常处理**：完善的错误处理和日志记录
- ✅ **回调功能**：支持回调呼叫中心或业务系统

### 技术栈
- **框架**：go-zero（微服务框架）
- **协议**：WebSocket / gRPC
- **存储**：Redis（会话和状态管理）
- **并发**：goroutine 高并发处理
- **日志**：go-zero 日志系统

## 项目结构

```
call/
├── etc/
│   └── asr.yaml              # 配置文件
├── internal/
│   ├── asr/                  # ASR 核心服务
│   │   └── service.go        # ASR 服务实现
│   ├── callback/             # 回调处理
│   │   └── client.go         # 回调客户端
│   ├── config/               # 配置
│   │   └── config.go         # 配置结构
│   ├── handler/              # 处理器
│   │   ├── asrhandler.go     # gRPC 处理器
│   │   └── websockethandler.go # WebSocket 处理器
│   ├── session/              # 会话管理
│   │   └── session.go        # 会话管理器
│   └── svc/                  # 服务上下文
│       └── servicecontext.go # 服务上下文
├── asr.proto                 # gRPC 协议定义
├── asrserver.go              # HTTP/WebSocket 服务器入口
├── asrrpcserver.go           # gRPC 服务器入口
├── go.mod                    # Go 模块定义
├── Makefile                  # 构建脚本
└── readme.md                 # 本文档
```

## 快速开始

### 1. 环境要求

- Go 1.23+
- Redis 6.0+
- protoc（用于生成 gRPC 代码）

### 2. 安装依赖

```bash
go mod download
```

### 3. 安装 protoc 和插件

安装 protoc 编译器：

**macOS:**
```bash
brew install protobuf
```

**Linux:**
```bash
apt-get install protobuf-compiler  # Debian/Ubuntu
yum install protobuf-compiler      # CentOS/RHEL
```

安装 Go 插件：
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 4. 生成 gRPC 代码

**重要：必须先生成 proto 文件才能编译项目！**

```bash
make proto
```

或者手动执行：

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    asr.proto
```

这将在当前目录生成 `pb/` 目录，包含生成的 Go 代码。

### 5. 配置 Redis

确保 Redis 服务正在运行，默认配置：
- Host: 127.0.0.1:6379
- DB: 0

### 6. 配置文件

编辑 `etc/asr.yaml` 配置文件：

```yaml
Name: asr-service
Host: 0.0.0.0
Port: 8080
Mode: dev

Redis:
  Host: 127.0.0.1:6379
  Pass: ""
  DB: 0

# ... 其他配置
```

### 7. 构建和运行

#### 运行 HTTP/WebSocket 服务器

```bash
make run-http
# 或
go run asrserver.go -f etc/asr.yaml
```

服务器将在 `http://0.0.0.0:8080` 启动

#### 运行 gRPC 服务器

```bash
make run-rpc
# 或
go run asrrpcserver.go -f etc/asr.yaml
```

gRPC 服务器将在 `0.0.0.0:9090` 启动

## API 文档

### WebSocket API

#### 连接 WebSocket

```
ws://localhost:8080/ws?call_id=<call_id>&phone_number=<phone_number>
```

**查询参数：**
- `call_id`（必需）：通话 ID
- `phone_number`（可选）：电话号码
- `session_id`（可选）：如果提供，将使用现有会话

#### 消息格式

**发送音频数据：**
- 类型：Binary Message
- 内容：音频字节流（PCM、WAV 等格式）

**发送控制消息：**
- 类型：Text Message（JSON）
- 格式：
  ```json
  {
    "action": "close"  // 或 "ping"
  }
  ```

**接收识别结果：**
- 类型：Text Message（JSON）
- 格式：
  ```json
  {
    "session_id": "xxx",
    "text": "识别结果",
    "confidence": 0.95,
    "timestamp": "2026-01-23T10:00:00Z",
    "is_final": false
  }
  ```

### gRPC API

#### 服务定义

参见 `asr.proto` 文件

#### 主要方法

1. **CreateSession** - 创建新会话
2. **ProcessAudioStream** - 处理音频流（双向流）
3. **ProcessAudio** - 处理单个音频块
4. **CloseSession** - 关闭会话
5. **GetSessionStatus** - 获取会话状态

### HTTP API

#### 健康检查

```bash
GET /health
```

#### 查询会话

```bash
GET /session/:sessionId
```

## 会话管理

### 会话生命周期

1. **创建**：通过 `CreateSession` 或 WebSocket 连接时自动创建
2. **活跃**：接收音频数据并处理
3. **完成**：识别完成或主动关闭
4. **失败**：处理异常时标记为失败

### 会话存储

- 存储位置：Redis
- Key 格式：`asr:session:{session_id}`
- TTL：5 分钟（可配置）
- 数据结构：JSON

### 会话状态

- `active`：活跃状态，正在接收音频
- `processing`：处理中
- `completed`：已完成
- `failed`：失败

## 流控机制

### 配置项

```yaml
RateLimit:
  QPS: 100        # 每秒请求数
  Burst: 200      # 突发请求数

ASR:
  MaxConcurrentSessions: 1000  # 最大并发会话数
```

### 流控策略

- **QPS 限制**：使用 go-zero 的限流器控制每秒请求数
- **并发限制**：限制同时处理的会话数量
- **缓冲区限制**：音频流缓冲区大小限制

## 超时处理

### 超时配置

```yaml
ASR:
  SessionTimeout: 300s  # 会话超时时间
```

### 超时机制

- **会话超时**：如果会话在指定时间内没有收到新的音频数据，将自动超时
- **连接超时**：WebSocket 连接超时（pongWait）
- **处理超时**：ASR 处理超时

## 回调功能

### 回调配置

```yaml
Callback:
  Timeout: 5s
  RetryTimes: 3
  RetryInterval: 1s
  Endpoints:
    - http://localhost:8081/callback
```

### 回调事件

1. **会话开始**：`OnSessionStart`
2. **会话更新**：`OnSessionUpdate`（实时识别结果）
3. **会话完成**：`OnSessionComplete`
4. **会话失败**：`OnSessionFailed`

### 回调数据格式

```json
{
  "session_id": "xxx",
  "call_id": "xxx",
  "phone_number": "xxx",
  "status": "started|processing|completed|failed",
  "result": "识别结果",
  "timestamp": "2026-01-23T10:00:00Z",
  "extra": {}
}
```

## 集成真实 ASR 引擎

当前实现使用的是模拟 ASR 服务。要集成真实的 ASR 引擎（如百度、阿里云、讯飞等），需要：

1. 实现 `RealASREngine` 接口
2. 在 `MockASRService` 中替换 `mockASRProcess` 方法
3. 配置相应的 API Key 和 Secret

示例：

```go
// 集成百度 ASR
func (s *MockASRService) processWithBaiduASR(audioData []byte) (*ASRResult, error) {
    // 调用百度 ASR API
    // ...
}
```

## 监控和日志

### 日志配置

```yaml
Log:
  ServiceName: asr-service
  Mode: file
  Path: logs
  Level: info
  Compress: true
  KeepDays: 7
```

### 监控指标

- 活跃会话数
- 处理中的会话数
- QPS
- 错误率
- 平均处理时间

## 性能优化

### 高并发处理

- 使用 goroutine 池处理并发请求
- 使用 channel 缓冲音频流
- Redis 连接池优化

### 资源管理

- 及时释放会话资源
- 控制并发会话数
- 合理设置缓冲区大小

## 部署

### Docker 部署

```dockerfile
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN make build

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/asr-server .
COPY --from=builder /app/etc/asr.yaml ./etc/
CMD ["./asr-server", "-f", "etc/asr.yaml"]
```

### Kubernetes 部署

参考 `deployment.yaml` 和 `service.yaml` 配置文件

## 故障排查

### 常见问题

1. **Redis 连接失败**
   - 检查 Redis 服务是否运行
   - 检查配置文件中的 Redis 地址

2. **会话创建失败**
   - 检查是否达到最大并发会话数
   - 检查 Redis 连接状态

3. **识别结果为空**
   - 检查音频数据格式
   - 检查 ASR 服务配置

## 开发指南

### 添加新功能

1. 在相应的 handler 中添加处理逻辑
2. 更新 proto 文件（如需要）
3. 添加单元测试
4. 更新文档

### 代码规范

- 遵循 Go 代码规范
- 使用 go-zero 的最佳实践
- 添加必要的注释和文档

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！
