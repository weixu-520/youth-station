// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type PaymentNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(req *types.PaymentNotifyRequest) (resp *types.BaseResponse, err error) {
	// 实际项目中应校验支付平台签名，这里简化处理
	var app model.Application
	if err := l.svcCtx.DB.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		logx.Errorf("查询申请失败: %v", err)
		return &types.BaseResponse{Code: 404, Message: "申请不存在"}, nil
	}

	if app.Status != 1 { // 只有审核通过的申请才能支付
		return &types.BaseResponse{Code: 400, Message: "当前状态不允许支付"}, nil
	}
	if app.DepositStatus == 1 {
		return &types.BaseResponse{Code: 400, Message: "押金已支付"}, nil
	}

	// 更新押金状态和支付记录
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(&app).Updates(map[string]interface{}{
			"deposit_status": 1,
		}).Error; err != nil {
			return err
		}
		payment := &model.Payment{
			ApplicationId: req.ApplicationId,
			TradeNo:       req.TradeNo,
			Amount:        req.Amount,
			PayType:       1, // 支付
			Status:        1, // 成功
			PayTime:       req.PayTime,
		}
		if err := tx.Create(payment).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("支付回调处理失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "处理失败"}, nil
	}

	return &types.BaseResponse{Code: 0, Message: "success"}, nil
}
