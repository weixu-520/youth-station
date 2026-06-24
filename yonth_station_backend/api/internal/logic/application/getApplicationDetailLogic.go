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

type GetApplicationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetApplicationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetApplicationDetailLogic {
	return &GetApplicationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetApplicationDetailLogic) GetApplicationDetail(req *types.GetApplicationDetailRequest) (resp *types.BaseResponse, err error) {
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

	// 权限校验：只能查看自己的申请（管理员可看所有，暂不实现）
	if app.UserId != userId {
		return &types.BaseResponse{Code: 403, Message: "无权查看"}, nil
	}

	// 查询驿站名称
	var station model.Station
	stationName := ""
	if err := l.svcCtx.DB.Where("id = ?", app.StationId).First(&station).Error; err == nil {
		stationName = station.StationName
	}

	// 查询用户信息
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", app.UserId).First(&user).Error; err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 脱敏
	maskedIdCard := utils.MaskIdCard(user.IdCard)
	maskedPhone := utils.MaskPhone(user.Phone)

	status := int32(app.Status)
	record := types.ApplicationRecord{
		ApplicationId: app.Id,
		UserId:        app.UserId,
		UserName:      user.UserName,
		StationId:     app.StationId,
		StationName:   stationName,
		CheckinDate:   app.CheckinDate,
		CheckoutDate:  app.CheckoutDate,
		Status:        &status,
		StatusDesc:    utils.GetStatusDesc(app.Status),
		VisitPurpose:  int32(app.VisitPurpose),
		RejectReason:  app.RejectReason,
		DepositAmount: app.DepositAmount,
		DepositStatus: int32(app.DepositStatus),
		CheckinAt:     app.CheckinAt,
		CheckoutAt:    app.CheckoutAt,
		AppliedAt:     app.AppliedAt,
		UpdatedAt:     app.UpdatedAt,
	}

	detail := types.ApplicationDetailResponse{
		ApplicationRecord: record,
		UserIdCard:        maskedIdCard,
		UserPhone:         maskedPhone,
		UserSchool:        user.School,
		UserHukou:         user.HukouCity,
		AuditBy:           app.AuditBy,
		AuditAt:           app.AuditAt,
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    &detail,
	}, nil
}
