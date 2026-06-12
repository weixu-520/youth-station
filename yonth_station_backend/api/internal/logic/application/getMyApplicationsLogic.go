// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyApplicationsLogic {
	return &GetMyApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyApplicationsLogic) GetMyApplications(req *types.ApplicationListRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
