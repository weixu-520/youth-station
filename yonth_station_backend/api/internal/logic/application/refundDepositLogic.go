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

type RefundDepositLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundDepositLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundDepositLogic {
	return &RefundDepositLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundDepositLogic) RefundDeposit(req *types.RefundDepositRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
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

	if app.Status != 5 {
		return &types.BaseResponse{Code: 400, Message: "未退房，无法退还押金"}, nil
	}
	if app.DepositStatus != 1 {
		return &types.BaseResponse{Code: 400, Message: "押金未缴纳或已退还"}, nil
	}

	// 生成退款流水号（简单示例）
	refundNo := "R" + time.Now().Format("20060102150405")
	now := time.Now().Unix()

	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&app).Update("deposit_status", 2).Error; err != nil {
			return err
		}
		payment := &model.Payment{
			ApplicationId: req.ApplicationId,
			TradeNo:       refundNo,
			Amount:        app.DepositAmount,
			PayType:       2, // 退款
			Status:        4, // 已退款
			RefundNo:      refundNo,
			RefundTime:    now,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("退还押金失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "操作失败"}, nil
	}

	return &types.BaseResponse{Code: 0, Message: "押金已退还"}, nil
}
