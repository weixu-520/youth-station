// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"
	"errors"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type CancelApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelApplicationLogic {
	return &CancelApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelApplicationLogic) CancelApplication(req *types.CancelApplicationRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	var app model.Application
	if err := l.svcCtx.DB.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 404, Message: "申请不存在"}, nil
		}
		logx.Errorf("查询申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	// 只能取消自己的申请，且状态必须为待审核(0)
	if app.UserId != userId {
		return &types.BaseResponse{Code: 403, Message: "无权操作"}, nil
	}
	if app.Status != 0 {
		return &types.BaseResponse{Code: 400, Message: "当前状态无法取消"}, nil
	}

	// 更新状态为已取消(3)，并恢复驿站配额和房间数
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&app).Update("status", 3).Error; err != nil {
			return err
		}
		// 恢复驿站剩余配额和可预约房间数
		if err := tx.Model(&model.Station{}).Where("id = ?", app.StationId).
			Updates(map[string]interface{}{
				"remaining_quota": gorm.Expr("remaining_quota + 1"),
				"available_rooms": gorm.Expr("available_rooms + 1"),
			}).Error; err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		logx.Errorf("取消申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "取消失败"}, nil
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "已取消申请",
		Data:    nil,
	}, nil
}
