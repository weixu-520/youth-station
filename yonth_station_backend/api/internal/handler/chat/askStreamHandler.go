package chat

import (
	"net/http"
	"strings"

	chatLogic "yonth_station_backend/api/internal/logic/chat"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/pkg/utils"
)

func AskStreamHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// CORS（流式绕过 go-zero，需手动设置）
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		if r.Method == "OPTIONS" { w.WriteHeader(204); return }
		// JWT 校验
		token := strings.TrimPrefix(r.Header.Get("Authorization"), "Bearer ")
		if _, err := utils.ParseToken(token, svcCtx.Config.Auth.AccessSecret); err != nil {
			http.Error(w, "unauthorized", 401)
			return
		}
		l := chatLogic.NewAskStreamLogic(svcCtx)
		l.AskStream(w, r)
	}
}
