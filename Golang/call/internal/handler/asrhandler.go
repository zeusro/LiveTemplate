package handler

import (
	"context"
	"io"

	"call/internal/session"
	"call/internal/svc"
	"call/pb"

	"github.com/zeromicro/go-zero/core/logx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ASRHandler gRPC 处理器
type ASRHandler struct {
	pb.UnimplementedASRServiceServer
	svcCtx *svc.ServiceContext
}

// NewASRHandler 创建 gRPC 处理器
func NewASRHandler(svcCtx *svc.ServiceContext) *ASRHandler {
	return &ASRHandler{
		svcCtx: svcCtx,
	}
}

// ProcessAudioStream 处理音频流
func (h *ASRHandler) ProcessAudioStream(stream pb.ASRService_ProcessAudioStreamServer) error {
	var sessionID string
	audioChan := make(chan []byte, 100)

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			// 流结束
			close(audioChan)
			break
		}
		if err != nil {
			logx.Errorf("Receive audio chunk failed: %v", err)
			return status.Errorf(codes.Internal, "receive failed: %v", err)
		}

		if sessionID == "" {
			sessionID = chunk.SessionId
			if sessionID == "" {
				return status.Error(codes.InvalidArgument, "session_id is required")
			}

			// 启动 ASR 处理
			results, err := h.svcCtx.ASRService.ProcessStream(stream.Context(), sessionID, audioChan)
			if err != nil {
				return status.Errorf(codes.Internal, "start ASR processing failed: %v", err)
			}

			// 启动结果发送协程
			go func() {
				for result := range results {
					pbResult := &pb.ASRResult{
						SessionId:  result.SessionID,
						Text:       result.Text,
						Confidence: result.Confidence,
						Timestamp:  result.Timestamp.Unix(),
						IsFinal:    result.IsFinal,
					}
					if err := stream.Send(pbResult); err != nil {
						logx.Errorf("Send ASR result failed: %v", err)
						return
					}

					// 更新会话
					if result.IsFinal {
						session, err := h.svcCtx.SessionMgr.GetSession(context.Background(), sessionID)
						if err == nil {
							h.svcCtx.SessionMgr.CompleteSession(context.Background(), sessionID, result.Text)
							h.svcCtx.Callback.OnSessionComplete(context.Background(), sessionID, session.CallID, session.PhoneNumber, result.Text)
						}
					}
				}
			}()
		}

		// 发送音频数据到处理通道
		select {
		case audioChan <- chunk.AudioData:
			h.svcCtx.SessionMgr.AddAudioChunk(stream.Context(), sessionID)
		case <-stream.Context().Done():
			return stream.Context().Err()
		default:
			logx.Infof("Audio stream buffer full for session %s", sessionID)
		}
	}

	return nil
}

// ProcessAudio 处理单个音频块
func (h *ASRHandler) ProcessAudio(ctx context.Context, req *pb.AudioChunk) (*pb.ASRResult, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	result, err := h.svcCtx.ASRService.ProcessAudio(ctx, req.SessionId, req.AudioData)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "process audio failed: %v", err)
	}

	// 更新会话
	h.svcCtx.SessionMgr.AddAudioChunk(ctx, req.SessionId)

	return &pb.ASRResult{
		SessionId:  result.SessionID,
		Text:       result.Text,
		Confidence: result.Confidence,
		Timestamp:  result.Timestamp.Unix(),
		IsFinal:    result.IsFinal,
	}, nil
}

// CreateSession 创建会话
func (h *ASRHandler) CreateSession(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	if req.CallId == "" {
		return nil, status.Error(codes.InvalidArgument, "call_id is required")
	}

	session, err := h.svcCtx.SessionMgr.CreateSession(ctx, req.CallId, req.PhoneNumber)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "create session failed: %v", err)
	}

	// 发送会话开始回调
	h.svcCtx.Callback.OnSessionStart(ctx, session.ID, req.CallId, req.PhoneNumber)

	return &pb.SessionResponse{
		SessionId: session.ID,
		Success:   true,
		Message:   "Session created successfully",
	}, nil
}

// CloseSession 关闭会话
func (h *ASRHandler) CloseSession(ctx context.Context, req *pb.SessionRequest) (*pb.SessionResponse, error) {
	if req.SessionId == "" {
		return nil, status.Error(codes.InvalidArgument, "session_id is required")
	}

	err := h.svcCtx.ASRService.CloseSession(ctx, req.SessionId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "close session failed: %v", err)
	}

	session, err := h.svcCtx.SessionMgr.GetSession(ctx, req.SessionId)
	if err == nil {
		h.svcCtx.Callback.OnSessionComplete(ctx, req.SessionId, session.CallID, session.PhoneNumber, session.Result)
	}

	return &pb.SessionResponse{
		SessionId: req.SessionId,
		Success:   true,
		Message:   "Session closed successfully",
	}, nil
}

// GetSessionStatus 获取会话状态
func (h *ASRHandler) GetSessionStatus(ctx context.Context, req *pb.SessionRequest) (*pb.SessionStatus, error) {
	var session *session.Session
	var err error

	if req.SessionId != "" {
		session, err = h.svcCtx.SessionMgr.GetSession(ctx, req.SessionId)
	} else if req.CallId != "" {
		session, err = h.svcCtx.SessionMgr.GetSessionByCallID(ctx, req.CallId)
	} else {
		return nil, status.Error(codes.InvalidArgument, "session_id or call_id is required")
	}

	if err != nil {
		return nil, status.Errorf(codes.NotFound, "session not found: %v", err)
	}

	return &pb.SessionStatus{
		SessionId:   session.ID,
		CallId:      session.CallID,
		PhoneNumber: session.PhoneNumber,
		Status:      session.Status,
		AudioChunks: int32(session.AudioChunks),
		Result:      session.Result,
		CreatedAt:   session.CreatedAt.Unix(),
		UpdatedAt:   session.UpdatedAt.Unix(),
	}, nil
}
