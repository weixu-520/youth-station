// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"
	"errors"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CheckinLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckinLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckinLogic {
	return &CheckinLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckinLogic) Checkin(req *types.CheckinRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	// 管理员权限校验
	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}

	var app model.Application
	if err := l.svcCtx.DB.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 404, Message: "申请不存在"}, nil
		}
		logx.Errorf("查询申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	if app.Status != 1 {
		return &types.BaseResponse{Code: 400, Message: "申请尚未通过，无法入住"}, nil
	}
	if app.DepositStatus != 1 {
		return &types.BaseResponse{Code: 400, Message: "押金未缴纳，请先支付押金"}, nil
	}
	if app.CheckinAt != 0 {
		return &types.BaseResponse{Code: 400, Message: "已办理入住，请勿重复操作"}, nil
	}

	now := time.Now().Unix()
	if err := l.svcCtx.DB.Model(&app).Updates(map[string]interface{}{
		"status":     4,
		"checkin_at": now,
		"updated_at": now,
	}).Error; err != nil {
		logx.Errorf("办理入住失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "操作失败"}, nil
	}

	return &types.BaseResponse{Code: 0, Message: "入住成功"}, nil
}
