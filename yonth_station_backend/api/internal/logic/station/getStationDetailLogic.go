// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"context"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStationDetailLogic {
	return &GetStationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStationDetailLogic) GetStationDetail() (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
