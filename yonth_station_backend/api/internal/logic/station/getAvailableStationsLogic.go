// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"context"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableStationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableStationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableStationsLogic {
	return &GetAvailableStationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableStationsLogic) GetAvailableStations() (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
