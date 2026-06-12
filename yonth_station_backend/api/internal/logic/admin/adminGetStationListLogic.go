// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetStationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetStationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetStationListLogic {
	return &AdminGetStationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetStationListLogic) AdminGetStationList(req *types.StationListRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
