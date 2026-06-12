// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yonth_station_backend/api/internal/logic/application"
	"yonth_station_backend/api/internal/svc"
)

func GetApplicationDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := application.NewGetApplicationDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetApplicationDetail()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
