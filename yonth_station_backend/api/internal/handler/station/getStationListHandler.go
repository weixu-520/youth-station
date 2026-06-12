// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yonth_station_backend/api/internal/logic/station"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
)

func GetStationListHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.StationListRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := station.NewGetStationListLogic(r.Context(), svcCtx)
		resp, err := l.GetStationList(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
