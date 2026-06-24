// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"net/http"

	"yonth_station_backend/api/internal/logic/station"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func GetStationDetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetStationDetailRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := station.NewGetStationDetailLogic(r.Context(), svcCtx)
		resp, err := l.GetStationDetail(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
