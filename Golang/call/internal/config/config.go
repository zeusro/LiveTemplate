package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/zrpc"
)

type Config struct {
	rest.RestConf
	zrpc.RpcServerConf
	Redis     redis.RedisConf
	WebSocket WebSocketConfig
	ASR       ASRConfig
	RateLimit RateLimitConfig
	Callback  CallbackConfig
	Log       LogConfig
}

type WebSocketConfig struct {
	ReadBufferSize   int    `json:",default=4096"`
	WriteBufferSize  int    `json:",default=4096"`
	CheckOrigin      bool   `json:",default=true"`
	HandshakeTimeout string `json:",default=10s"`
}

type ASRConfig struct {
	MaxConcurrentSessions int    `json:",default=1000"`
	SessionTimeout        string `json:",default=300s"`
	AudioBufferSize       int    `json:",default=8192"`
	MaxMessageSize        int64  `json:",default=10485760"`
}

type RateLimitConfig struct {
	QPS   int `json:",default=100"`
	Burst int `json:",default=200"`
}

type CallbackConfig struct {
	Timeout      string   `json:",default=5s"`
	RetryTimes   int      `json:",default=3"`
	RetryInterval string   `json:",default=1s"`
	Endpoints    []string `json:",optional"`
}

type LogConfig struct {
	ServiceName         string `json:",default=asr-service"`
	Mode               string `json:",default=file"`
	Path               string `json:",default=logs"`
	Level              string `json:",default=info"`
	Compress           bool   `json:",default=true"`
	KeepDays           int    `json:",default=7"`
	StackCooldownMillis int    `json:",default=100"`
}
