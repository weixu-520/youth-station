package chat

import (
	"encoding/json"
	"net/http"

	"yonth_station_backend/api/internal/logic/chat"
	"yonth_station_backend/api/internal/svc"
)

func GetHistoryHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chat.NewGetHistoryLogic(r.Context(), svcCtx)
		resp, _ := l.GetHistory()
		// 简化的 JSON 响应
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		data, _ := json.Marshal(resp)
		w.Write(data)
	}
}
