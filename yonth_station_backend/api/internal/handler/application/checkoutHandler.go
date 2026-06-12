// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yonth_station_backend/api/internal/logic/application"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
)

func CheckoutHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CheckoutRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := application.NewCheckoutLogic(r.Context(), svcCtx)
		resp, err := l.Checkout(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
