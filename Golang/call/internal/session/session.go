package session

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	SessionPrefix = "asr:session:"
	SessionTTL    = 5 * time.Minute
)

type Session struct {
	ID          string    `json:"id"`
	CallID      string    `json:"call_id"`
	PhoneNumber string    `json:"phone_number"`
	Status      string    `json:"status"` // active, processing, completed, failed
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	AudioChunks int       `json:"audio_chunks"`
	Result      string    `json:"result,omitempty"`
}

type Manager struct {
	redis *redis.Client
}

func NewManager(redisClient *redis.Client) *Manager {
	return &Manager{
		redis: redisClient,
	}
}

// CreateSession 创建新会话
func (m *Manager) CreateSession(ctx context.Context, callID, phoneNumber string) (*Session, error) {
	session := &Session{
		ID:          uuid.New().String(),
		CallID:      callID,
		PhoneNumber: phoneNumber,
		Status:      "active",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
		AudioChunks: 0,
	}

	key := SessionPrefix + session.ID
	data, err := json.Marshal(session)
	if err != nil {
		return nil, fmt.Errorf("marshal session: %w", err)
	}

	err = m.redis.Set(ctx, key, data, SessionTTL).Err()
	if err != nil {
		return nil, fmt.Errorf("save session to redis: %w", err)
	}

	// 同时存储 call_id 到 session_id 的映射
	callKey := SessionPrefix + "call:" + callID
	err = m.redis.Set(ctx, callKey, session.ID, SessionTTL).Err()
	if err != nil {
		logx.Errorf("save call mapping failed: %v", err)
	}

	return session, nil
}

// GetSession 获取会话
func (m *Manager) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	key := SessionPrefix + sessionID
	data, err := m.redis.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found")
		}
		return nil, fmt.Errorf("get session from redis: %w", err)
	}

	var session Session
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, fmt.Errorf("unmarshal session: %w", err)
	}

	return &session, nil
}

// GetSessionByCallID 通过 CallID 获取会话
func (m *Manager) GetSessionByCallID(ctx context.Context, callID string) (*Session, error) {
	callKey := SessionPrefix + "call:" + callID
	sessionID, err := m.redis.Get(ctx, callKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("session not found for call_id: %s", callID)
		}
		return nil, fmt.Errorf("get session id from redis: %w", err)
	}

	return m.GetSession(ctx, sessionID)
}

// UpdateSession 更新会话
func (m *Manager) UpdateSession(ctx context.Context, session *Session) error {
	session.UpdatedAt = time.Now()
	key := SessionPrefix + session.ID
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal session: %w", err)
	}

	err = m.redis.Set(ctx, key, data, SessionTTL).Err()
	if err != nil {
		return fmt.Errorf("update session in redis: %w", err)
	}

	return nil
}

// AddAudioChunk 增加音频块计数
func (m *Manager) AddAudioChunk(ctx context.Context, sessionID string) error {
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.AudioChunks++
	return m.UpdateSession(ctx, session)
}

// CompleteSession 完成会话
func (m *Manager) CompleteSession(ctx context.Context, sessionID, result string) error {
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Status = "completed"
	session.Result = result
	return m.UpdateSession(ctx, session)
}

// FailSession 标记会话失败
func (m *Manager) FailSession(ctx context.Context, sessionID, reason string) error {
	session, err := m.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	session.Status = "failed"
	session.Result = reason
	return m.UpdateSession(ctx, session)
}

// DeleteSession 删除会话
func (m *Manager) DeleteSession(ctx context.Context, sessionID string) error {
	key := SessionPrefix + sessionID
	return m.redis.Del(ctx, key).Err()
}

// ListActiveSessions 列出活跃会话（用于监控）
func (m *Manager) ListActiveSessions(ctx context.Context) ([]*Session, error) {
	pattern := SessionPrefix + "*"
	keys, err := m.redis.Keys(ctx, pattern).Result()
	if err != nil {
		return nil, fmt.Errorf("list session keys: %w", err)
	}

	var sessions []*Session
	for _, key := range keys {
		// 跳过 call: 映射键
		if len(key) > len(SessionPrefix+"call:") && key[:len(SessionPrefix+"call:")] == SessionPrefix+"call:" {
			continue
		}

		data, err := m.redis.Get(ctx, key).Bytes()
		if err != nil {
			continue
		}

		var session Session
		if err := json.Unmarshal(data, &session); err != nil {
			continue
		}

		if session.Status == "active" || session.Status == "processing" {
			sessions = append(sessions, &session)
		}
	}

	return sessions, nil
}
