package callback

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/zeromicro/go-zero/core/logx"
)

// CallbackData 回调数据
type CallbackData struct {
	SessionID  string                 `json:"session_id"`
	CallID     string                 `json:"call_id"`
	PhoneNumber string                `json:"phone_number"`
	Status     string                 `json:"status"`
	Result     string                 `json:"result,omitempty"`
	Timestamp  time.Time              `json:"timestamp"`
	Extra      map[string]interface{} `json:"extra,omitempty"`
}

// Client 回调客户端
type Client struct {
	endpoints     []string
	timeout       time.Duration
	retryTimes    int
	retryInterval time.Duration
	httpClient    *http.Client
}

// NewClient 创建回调客户端
func NewClient(endpoints []string, timeout time.Duration, retryTimes int, retryInterval time.Duration) *Client {
	return &Client{
		endpoints:     endpoints,
		timeout:       timeout,
		retryTimes:    retryTimes,
		retryInterval: retryInterval,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// SendCallback 发送回调
func (c *Client) SendCallback(ctx context.Context, data *CallbackData) error {
	if len(c.endpoints) == 0 {
		logx.Infof("No callback endpoints configured, skipping callback for session %s", data.SessionID)
		return nil
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal callback data: %w", err)
	}

	var lastErr error
	for _, endpoint := range c.endpoints {
		for i := 0; i < c.retryTimes; i++ {
			err := c.sendRequest(ctx, endpoint, jsonData)
			if err == nil {
				logx.Infof("Callback sent successfully to %s for session %s", endpoint, data.SessionID)
				return nil
			}
			lastErr = err
			if i < c.retryTimes-1 {
				time.Sleep(c.retryInterval)
			}
		}
	}

	return fmt.Errorf("failed to send callback after %d retries: %w", c.retryTimes, lastErr)
}

// SendCallbackAsync 异步发送回调
func (c *Client) SendCallbackAsync(ctx context.Context, data *CallbackData) {
	go func() {
		if err := c.SendCallback(ctx, data); err != nil {
			logx.Errorf("Async callback failed for session %s: %v", data.SessionID, err)
		}
	}()
}

func (c *Client) sendRequest(ctx context.Context, endpoint string, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint, bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("callback returned status %d: %s", resp.StatusCode, string(body))
	}

	return nil
}

// OnSessionStart 会话开始回调
func (c *Client) OnSessionStart(ctx context.Context, sessionID, callID, phoneNumber string) {
	data := &CallbackData{
		SessionID:  sessionID,
		CallID:     callID,
		PhoneNumber: phoneNumber,
		Status:     "started",
		Timestamp:  time.Now(),
	}
	c.SendCallbackAsync(ctx, data)
}

// OnSessionUpdate 会话更新回调
func (c *Client) OnSessionUpdate(ctx context.Context, sessionID, callID, phoneNumber, result string) {
	data := &CallbackData{
		SessionID:  sessionID,
		CallID:     callID,
		PhoneNumber: phoneNumber,
		Status:     "processing",
		Result:     result,
		Timestamp:  time.Now(),
	}
	c.SendCallbackAsync(ctx, data)
}

// OnSessionComplete 会话完成回调
func (c *Client) OnSessionComplete(ctx context.Context, sessionID, callID, phoneNumber, result string) {
	data := &CallbackData{
		SessionID:  sessionID,
		CallID:     callID,
		PhoneNumber: phoneNumber,
		Status:     "completed",
		Result:     result,
		Timestamp:  time.Now(),
	}
	c.SendCallbackAsync(ctx, data)
}

// OnSessionFailed 会话失败回调
func (c *Client) OnSessionFailed(ctx context.Context, sessionID, callID, phoneNumber, reason string) {
	data := &CallbackData{
		SessionID:  sessionID,
		CallID:     callID,
		PhoneNumber: phoneNumber,
		Status:     "failed",
		Result:     reason,
		Timestamp:  time.Now(),
	}
	c.SendCallbackAsync(ctx, data)
}
