package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"call/internal/asr"
	"call/internal/svc"

	"github.com/gorilla/websocket"
	"github.com/zeromicro/go-zero/core/logx"
)

const (
	// 写超时
	writeWait = 10 * time.Second
	// 读超时
	pongWait = 60 * time.Second
	// ping 间隔
	pingPeriod = (pongWait * 9) / 10
	// 最大消息大小
	maxMessageSize = 10 * 1024 * 1024 // 10MB
)

// WebSocketHandler WebSocket 处理器
type WebSocketHandler struct {
	svcCtx *svc.ServiceContext
	upgrader websocket.Upgrader
}

// NewWebSocketHandler 创建 WebSocket 处理器
func NewWebSocketHandler(svcCtx *svc.ServiceContext) *WebSocketHandler {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  svcCtx.Config.WebSocket.ReadBufferSize,
		WriteBufferSize: svcCtx.Config.WebSocket.WriteBufferSize,
		CheckOrigin: func(r *http.Request) bool {
			return svcCtx.Config.WebSocket.CheckOrigin
		},
	}

	return &WebSocketHandler{
		svcCtx:   svcCtx,
		upgrader: upgrader,
	}
}

// HandleWebSocket 处理 WebSocket 连接
func (h *WebSocketHandler) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// 升级连接
	conn, err := h.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.Errorf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	// 从查询参数获取会话信息
	sessionID := r.URL.Query().Get("session_id")
	callID := r.URL.Query().Get("call_id")
	phoneNumber := r.URL.Query().Get("phone_number")

	if sessionID == "" {
		// 如果没有 session_id，创建新会话
		if callID == "" {
			conn.WriteJSON(map[string]interface{}{
				"error": "call_id is required",
			})
			return
		}

		ctx := r.Context()
		session, err := h.svcCtx.SessionMgr.CreateSession(ctx, callID, phoneNumber)
		if err != nil {
			logx.Errorf("Create session failed: %v", err)
			conn.WriteJSON(map[string]interface{}{
				"error": "failed to create session",
			})
			return
		}
		sessionID = session.ID

		// 发送会话开始回调
		h.svcCtx.Callback.OnSessionStart(ctx, sessionID, callID, phoneNumber)
	}

	// 设置连接参数
	conn.SetReadDeadline(time.Now().Add(pongWait))
	conn.SetReadLimit(maxMessageSize)
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	// 创建音频流通道
	audioStream := make(chan []byte, 100)
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	// 启动 ASR 处理
	results, err := h.svcCtx.ASRService.ProcessStream(ctx, sessionID, audioStream)
	if err != nil {
		logx.Errorf("Start ASR stream failed: %v", err)
		conn.WriteJSON(map[string]interface{}{
			"error": "failed to start ASR processing",
		})
		return
	}

	// 启动 ping ticker
	pingTicker := time.NewTicker(pingPeriod)
	defer pingTicker.Stop()

	// 启动写协程
	go h.writePump(conn, results, sessionID)

	// 读循环
	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				logx.Errorf("WebSocket error: %v", err)
			}
			break
		}

		if messageType == websocket.BinaryMessage {
			// 二进制消息（音频数据）
			select {
			case audioStream <- message:
				// 更新会话
				h.svcCtx.SessionMgr.AddAudioChunk(ctx, sessionID)
			case <-ctx.Done():
				return
			default:
				logx.Infof("Audio stream buffer full for session %s", sessionID)
			}
		} else if messageType == websocket.TextMessage {
			// 文本消息（控制消息）
			var msg map[string]interface{}
			if err := json.Unmarshal(message, &msg); err != nil {
				logx.Errorf("Invalid message format: %v", err)
				continue
			}

			action := msg["action"]
			switch action {
			case "close":
				// 关闭会话
				close(audioStream)
				return
			case "ping":
				// 心跳
				conn.WriteJSON(map[string]interface{}{
					"action": "pong",
				})
			}
		}
	}

	// 关闭音频流
	close(audioStream)
	
	// 完成会话
	session, err := h.svcCtx.SessionMgr.GetSession(ctx, sessionID)
	if err == nil {
		h.svcCtx.SessionMgr.CompleteSession(ctx, sessionID, "WebSocket closed")
		h.svcCtx.Callback.OnSessionComplete(ctx, sessionID, session.CallID, session.PhoneNumber, session.Result)
	}
}

// writePump 写协程
func (h *WebSocketHandler) writePump(conn *websocket.Conn, results <-chan *asr.ASRResult, sessionID string) {
	defer conn.Close()

	for {
		select {
		case result, ok := <-results:
			conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 发送识别结果
			if err := conn.WriteJSON(result); err != nil {
				logx.Errorf("Write ASR result failed: %v", err)
				return
			}

			// 更新会话
			if result.IsFinal {
				ctx := context.Background()
				session, err := h.svcCtx.SessionMgr.GetSession(ctx, sessionID)
				if err == nil {
					h.svcCtx.SessionMgr.CompleteSession(ctx, sessionID, result.Text)
					h.svcCtx.Callback.OnSessionComplete(ctx, sessionID, session.CallID, session.PhoneNumber, result.Text)
				}
			}
		}
	}
}
