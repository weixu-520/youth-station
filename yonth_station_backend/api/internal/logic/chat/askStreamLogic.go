package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AskStreamLogic struct {
	logx.Logger
	svcCtx *svc.ServiceContext
}

func NewAskStreamLogic(svcCtx *svc.ServiceContext) *AskStreamLogic {
	return &AskStreamLogic{svcCtx: svcCtx}
}

func (l *AskStreamLogic) AskStream(w http.ResponseWriter, r *http.Request) {
	var req types.ChatRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", 400)
		return
	}
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", 500)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	sessionId := req.SessionId
	if sessionId == "" {
		sessionId = generateSessionId()
	}
	fmt.Fprintf(w, "data: {\"type\":\"session\",\"sessionId\":\"%s\"}\n\n", sessionId)
	flusher.Flush()

	// 取最近 10 条对话历史
	var historyMsgs []model.ChatHistory
	l.svcCtx.DB.Where("session_id = ?", sessionId).
		Order("created_at DESC").Limit(10).Find(&historyMsgs)
	history := ""
	for i := len(historyMsgs) - 1; i >= 0; i-- {
		history += fmt.Sprintf("用户：%s\n客服：%s\n", historyMsgs[i].Question, historyMsgs[i].Answer)
	}

	// 调用 RAG 流式查询（带历史）
	stream, err := l.svcCtx.RAGWorkflow.QueryStreamWithHistory(r.Context(), req.Question, history)
	if err != nil {
		fmt.Fprintf(w, "data: {\"type\":\"error\",\"msg\":\"%s\"}\n\n", err.Error())
		flusher.Flush()
		return
	}
	defer stream.Close()

	var fullAnswer string
	for {
		msg, err := stream.Recv()
		if err != nil {
			break
		}
		if msg.Content != "" {
			fullAnswer += msg.Content
			data, _ := json.Marshal(map[string]string{"type": "text", "content": msg.Content})
			fmt.Fprintf(w, "data: %s\n\n", string(data))
			flusher.Flush()
		}
	}
	fmt.Fprintf(w, "data: {\"type\":\"done\"}\n\n")
	flusher.Flush()

	go saveHistory(l.svcCtx, r.Context(), sessionId, req.Question, fullAnswer)
}

func saveHistory(svcCtx *svc.ServiceContext, ctx context.Context, sessionId, question, answer string) {
	userId, _ := utils.GetUserId(ctx)
	record := &model.ChatHistory{
		SessionId: sessionId, UserId: userId,
		Question: question, Answer: answer, CreatedAt: time.Now(),
	}
	_ = svcCtx.DB.Create(record).Error
}
