package chat

import (
	"encoding/json"
	"net/http"

	chatLogic "yonth_station_backend/api/internal/logic/chat"
	"yonth_station_backend/api/internal/svc"
)

func GetUsersHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := chatLogic.NewGetUsersLogic(r.Context(), svcCtx)
		resp, _ := l.GetUsers()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
