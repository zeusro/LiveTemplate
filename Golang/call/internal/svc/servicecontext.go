package svc

import (
	"time"

	"call/internal/asr"
	"call/internal/callback"
	"call/internal/config"
	"call/internal/session"

	"github.com/go-redis/redis/v8"
)

type ServiceContext struct {
	Config      config.Config
	Redis       *redis.Client
	SessionMgr  *session.Manager
	ASRService  asr.Service
	Callback    *callback.Client
}

func NewServiceContext(c config.Config) *ServiceContext {
	// 初始化 Redis Client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Host,
		Password: c.Redis.Pass,
		DB:       0, // 默认使用 DB 0
	})

	// 初始化会话管理器
	sessionMgr := session.NewManager(redisClient)

	// 初始化 ASR 服务
	asrService := asr.NewMockASRService(
		c.ASR.MaxConcurrentSessions,
		parseDuration(c.ASR.SessionTimeout),
		c.RateLimit.QPS,
	)

	// 初始化回调客户端
	callbackClient := callback.NewClient(
		c.Callback.Endpoints,
		parseDuration(c.Callback.Timeout),
		c.Callback.RetryTimes,
		parseDuration(c.Callback.RetryInterval),
	)

	return &ServiceContext{
		Config:     c,
		Redis:      redisClient,
		SessionMgr: sessionMgr,
		ASRService: asrService,
		Callback:   callbackClient,
	}
}

func parseDuration(s string) time.Duration {
	d, err := time.ParseDuration(s)
	if err != nil {
		// 默认值
		return 5 * time.Minute
	}
	return d
}
