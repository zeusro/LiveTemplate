# ASR Service

An Automatic Speech Recognition (ASR) service implemented based on the go-zero framework, supporting both WebSocket and gRPC protocols, with complete session management, flow control, timeout handling, and callback functionality.

## Features

### Core Features
- ✅ **ASR Service Wrapper**: Supports both WebSocket and gRPC protocols
- ✅ **Session Management**: One independent session per call, stored in Redis
- ✅ **Flow Control**: QPS limits and concurrent session limits
- ✅ **Timeout Handling**: Automatic session timeout cleanup
- ✅ **Error Handling**: Comprehensive error handling and logging
- ✅ **Callback Functionality**: Supports callbacks to call centers or business systems

### Technology Stack
- **Framework**: go-zero (microservices framework)
- **Protocols**: WebSocket / gRPC
- **Storage**: Redis (session and state management)
- **Concurrency**: High-concurrency processing with goroutines
- **Logging**: go-zero logging system

## Project Structure

```
call/
├── etc/
│   └── asr.yaml              # Configuration file
├── internal/
│   ├── asr/                  # ASR core service
│   │   └── service.go        # ASR service implementation
│   ├── callback/             # Callback handling
│   │   └── client.go         # Callback client
│   ├── config/               # Configuration
│   │   └── config.go         # Configuration structure
│   ├── handler/              # Handlers
│   │   ├── asrhandler.go     # gRPC handler
│   │   └── websockethandler.go # WebSocket handler
│   ├── session/              # Session management
│   │   └── session.go        # Session manager
│   └── svc/                  # Service context
│       └── servicecontext.go # Service context
├── asr.proto                 # gRPC protocol definition
├── asrserver.go              # HTTP/WebSocket server entry
├── asrrpcserver.go           # gRPC server entry
├── go.mod                    # Go module definition
├── Makefile                  # Build script
└── readme.md                 # This document
```

## Quick Start

### 1. Requirements

- Go 1.23+
- Redis 6.0+
- protoc (for generating gRPC code)

### 2. Install Dependencies

```bash
go mod download
```

### 3. Install protoc and Plugins

Install protoc compiler:

**macOS:**
```bash
brew install protobuf
```

**Linux:**
```bash
apt-get install protobuf-compiler  # Debian/Ubuntu
yum install protobuf-compiler      # CentOS/RHEL
```

Install Go plugins:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

### 4. Generate gRPC Code

**Important: You must generate the proto files before compiling the project!**

```bash
make proto
```

Or execute manually:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    asr.proto
```

This will generate a `pb/` directory in the current directory, containing the generated Go code.

### 5. Configure Redis

Ensure Redis service is running with default configuration:
- Host: 127.0.0.1:6379
- DB: 0

### 6. Configuration File

Edit the `etc/asr.yaml` configuration file:

```yaml
Name: asr-service
Host: 0.0.0.0
Port: 8080
Mode: dev

Redis:
  Host: 127.0.0.1:6379
  Pass: ""
  DB: 0

# ... other configurations
```

### 7. Build and Run

#### Run HTTP/WebSocket Server

```bash
make run-http
# or
go run asrserver.go -f etc/asr.yaml
```

The server will start at `http://0.0.0.0:8080`

#### Run gRPC Server

```bash
make run-rpc
# or
go run asrrpcserver.go -f etc/asr.yaml
```

The gRPC server will start at `0.0.0.0:9090`

## API Documentation

### WebSocket API

#### Connect WebSocket

```
ws://localhost:8080/ws?call_id=<call_id>&phone_number=<phone_number>
```

**Query Parameters:**
- `call_id` (required): Call ID
- `phone_number` (optional): Phone number
- `session_id` (optional): If provided, will use existing session

#### Message Format

**Send Audio Data:**
- Type: Binary Message
- Content: Audio byte stream (PCM, WAV, etc.)

**Send Control Message:**
- Type: Text Message (JSON)
- Format:
  ```json
  {
    "action": "close"  // or "ping"
  }
  ```

**Receive Recognition Results:**
- Type: Text Message (JSON)
- Format:
  ```json
  {
    "session_id": "xxx",
    "text": "Recognition result",
    "confidence": 0.95,
    "timestamp": "2026-01-23T10:00:00Z",
    "is_final": false
  }
  ```

### gRPC API

#### Service Definition

See the `asr.proto` file

#### Main Methods

1. **CreateSession** - Create a new session
2. **ProcessAudioStream** - Process audio stream (bidirectional stream)
3. **ProcessAudio** - Process a single audio chunk
4. **CloseSession** - Close a session
5. **GetSessionStatus** - Get session status

### HTTP API

#### Health Check

```bash
GET /health
```

#### Query Session

```bash
GET /session/:sessionId
```

## Session Management

### Session Lifecycle

1. **Create**: Automatically created via `CreateSession` or when WebSocket connects
2. **Active**: Receiving and processing audio data
3. **Completed**: Recognition completed or actively closed
4. **Failed**: Marked as failed when processing exceptions occur

### Session Storage

- Storage Location: Redis
- Key Format: `asr:session:{session_id}`
- TTL: 5 minutes (configurable)
- Data Structure: JSON

### Session Status

- `active`: Active state, receiving audio
- `processing`: Processing
- `completed`: Completed
- `failed`: Failed

## Flow Control

### Configuration

```yaml
RateLimit:
  QPS: 100        # Requests per second
  Burst: 200      # Burst requests

ASR:
  MaxConcurrentSessions: 1000  # Maximum concurrent sessions
```

### Flow Control Strategy

- **QPS Limit**: Use go-zero's rate limiter to control requests per second
- **Concurrency Limit**: Limit the number of sessions processed simultaneously
- **Buffer Limit**: Audio stream buffer size limit

## Timeout Handling

### Timeout Configuration

```yaml
ASR:
  SessionTimeout: 300s  # Session timeout duration
```

### Timeout Mechanism

- **Session Timeout**: If a session doesn't receive new audio data within the specified time, it will automatically timeout
- **Connection Timeout**: WebSocket connection timeout (pongWait)
- **Processing Timeout**: ASR processing timeout

## Callback Functionality

### Callback Configuration

```yaml
Callback:
  Timeout: 5s
  RetryTimes: 3
  RetryInterval: 1s
  Endpoints:
    - http://localhost:8081/callback
```

### Callback Events

1. **Session Start**: `OnSessionStart`
2. **Session Update**: `OnSessionUpdate` (real-time recognition results)
3. **Session Complete**: `OnSessionComplete`
4. **Session Failed**: `OnSessionFailed`

### Callback Data Format

```json
{
  "session_id": "xxx",
  "call_id": "xxx",
  "phone_number": "xxx",
  "status": "started|processing|completed|failed",
  "result": "Recognition result",
  "timestamp": "2026-01-23T10:00:00Z",
  "extra": {}
}
```

## Integrating Real ASR Engine

The current implementation uses a mock ASR service. To integrate a real ASR engine (such as Baidu, Alibaba Cloud, iFlytek, etc.), you need to:

1. Implement the `RealASREngine` interface
2. Replace the `mockASRProcess` method in `MockASRService`
3. Configure the corresponding API Key and Secret

Example:

```go
// Integrate Baidu ASR
func (s *MockASRService) processWithBaiduASR(audioData []byte) (*ASRResult, error) {
    // Call Baidu ASR API
    // ...
}
```

## Monitoring and Logging

### Log Configuration

```yaml
Log:
  ServiceName: asr-service
  Mode: file
  Path: logs
  Level: info
  Compress: true
  KeepDays: 7
```

### Monitoring Metrics

- Active session count
- Sessions in processing
- QPS
- Error rate
- Average processing time

## Performance Optimization

### High Concurrency Processing

- Use goroutine pools to handle concurrent requests
- Use channels to buffer audio streams
- Redis connection pool optimization

### Resource Management

- Release session resources promptly
- Control concurrent session count
- Set buffer size appropriately

## Deployment

### Docker Deployment

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

### Kubernetes Deployment

Refer to `deployment.yaml` and `service.yaml` configuration files

## Troubleshooting

### Common Issues

1. **Redis Connection Failed**
   - Check if Redis service is running
   - Check Redis address in configuration file

2. **Session Creation Failed**
   - Check if maximum concurrent sessions limit is reached
   - Check Redis connection status

3. **Empty Recognition Results**
   - Check audio data format
   - Check ASR service configuration

## Development Guide

### Adding New Features

1. Add processing logic in the corresponding handler
2. Update proto file (if needed)
3. Add unit tests
4. Update documentation

### Code Standards

- Follow Go code standards
- Use go-zero best practices
- Add necessary comments and documentation

## License

MIT License

## Contributing

Issues and Pull Requests are welcome!
