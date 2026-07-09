package chat

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AskLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAskLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AskLogic {
	return &AskLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AskLogic) Ask(req *types.ChatRequest) (resp *types.BaseResponse, err error) {
	// 1. 获取或生成会话 ID
	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = generateSessionId()
	}

	// 2. 执行 RAG 查询
	answer, err := l.svcCtx.RAGWorkflow.Query(l.ctx, req.Question)
	if err != nil {
		logx.Errorf("RAG query failed: %v", err)
		return &types.BaseResponse{Code: 500, Message: fmt.Sprintf("智能客服暂时不可用：%v", err)}, nil
	}

	// 3. 异步保存对话记录
	go l.saveHistory(sessionId, req.Question, answer)

	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data: &types.ChatResponse{
			Answer:    answer,
			SessionId: sessionId,
		},
	}, nil
}

func (l *AskLogic) saveHistory(sessionId, question, answer string) {
	record := &model.ChatHistory{
		SessionId: sessionId,
		UserId:    getUserIdFromCtx(l.ctx), // 需实现获取当前用户 ID
		Question:  question,
		Answer:    answer,
		CreatedAt: time.Now(),
	}
	_ = l.svcCtx.DB.Create(record).Error
}

func generateSessionId() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
func getUserIdFromCtx(ctx context.Context) int64 {
	userId, _ := utils.GetUserId(ctx)
	return userId
}
