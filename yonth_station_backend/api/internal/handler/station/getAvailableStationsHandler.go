// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"yonth_station_backend/api/internal/logic/station"
	"yonth_station_backend/api/internal/svc"
)

func GetAvailableStationsHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := station.NewGetAvailableStationsLogic(r.Context(), svcCtx)
		resp, err := l.GetAvailableStations()
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
