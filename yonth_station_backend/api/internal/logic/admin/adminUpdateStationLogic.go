// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateStationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateStationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateStationLogic {
	return &AdminUpdateStationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateStationLogic) AdminUpdateStation(req *types.StationInfo) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
