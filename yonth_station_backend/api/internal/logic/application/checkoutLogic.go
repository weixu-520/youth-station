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

type CheckoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCheckoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckoutLogic {
	return &CheckoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CheckoutLogic) Checkout(req *types.CheckoutRequest) (resp *types.BaseResponse, err error) {
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

	if app.Status != 4 {
		return &types.BaseResponse{Code: 400, Message: "未入住或已退房"}, nil
	}

	now := time.Now().Unix()
	//办理退房事务
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		// 更新申请状态为已退房
		if err := tx.Model(&app).Updates(map[string]interface{}{
			"status":      5,
			"checkout_at": now,
			"updated_at":  now,
		}).Error; err != nil {
			return err
		}
		// 释放房间（将房间状态置为空闲）
		if app.RoomId != 0 {
			if err := tx.Model(&model.Room{}).Where("id = ?", app.RoomId).Update("status", 0).Error; err != nil {
				return err
			}
		}
		// 恢复驿站的可用房间数（注意：不恢复每周配额 remaining_quota）
		if err := tx.Model(&model.Station{}).Where("id = ?", app.StationId).
			Update("available_rooms", gorm.Expr("available_rooms + 1")).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("办理退房失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "操作失败"}, nil
	}

	return &types.BaseResponse{Code: 0, Message: "退房成功"}, nil
}
