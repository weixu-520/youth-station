// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"
	"errors"
	"fmt"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type ApplyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewApplyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ApplyLogic {
	return &ApplyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ApplyLogic) Apply(req *types.ApplicationRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	// 2. 查询驿站信息（检查是否存在及状态）
	var station model.Station
	if err := l.svcCtx.DB.Where("id = ? AND status = 1", req.StationId).First(&station).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 400, Message: "驿站不存在或已关闭"}, nil
		}
		logx.Errorf("查询驿站失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 3. 校验入住日期和退房日期（最多7天）
	checkin, err := time.Parse("2006-01-02", req.CheckinDate)
	if err != nil {
		return &types.BaseResponse{Code: 400, Message: "入住日期格式错误"}, nil
	}
	checkout, err := time.Parse("2006-01-02", req.CheckoutDate)
	if err != nil {
		return &types.BaseResponse{Code: 400, Message: "退房日期格式错误"}, nil
	}
	if checkout.Before(checkin) || checkout.Equal(checkin) {
		return &types.BaseResponse{Code: 400, Message: "退房日期必须晚于入住日期"}, nil
	}
	days := int(checkout.Sub(checkin).Hours() / 24)
	if days > 7 {
		return &types.BaseResponse{Code: 400, Message: "最多可连续入住7天"}, nil
	}

	// 4. 初始化 Redis 配额（key 不存在或值无效时，从 DB 同步）
	quotaKey := fmt.Sprintf("station:quota:%d", req.StationId)
	roomsKey := fmt.Sprintf("station:rooms:%d", req.StationId)
	initQuota := fmt.Sprintf("%d", station.RemainingQuota)
	initRooms := fmt.Sprintf("%d", station.AvailableRooms)
	// 只在 key 不存在时用 SETNX 写入，防止覆盖正常值
	l.svcCtx.Redis.Setnx(quotaKey, initQuota)
	l.svcCtx.Redis.Setnx(roomsKey, initRooms)
	// 兜底：如果 key 存在但值 ≤ 0（被旧数据污染），强制覆盖
	if val, _ := l.svcCtx.Redis.Get(quotaKey); val == "0" || val == "-1" {
		l.svcCtx.Redis.Set(quotaKey, initQuota)
	}
	if val, _ := l.svcCtx.Redis.Get(roomsKey); val == "0" || val == "-1" {
		l.svcCtx.Redis.Set(roomsKey, initRooms)
	}

	// 5. Lua 脚本原子扣减配额和房间数
	luaScript := `
	local quota = redis.call('DECR', KEYS[1])
	local rooms = redis.call('DECR', KEYS[2])
	if quota < 0 then
		redis.call('INCR', KEYS[1])
		redis.call('INCR', KEYS[2])
		return -1
	end
	if rooms < 0 then
		redis.call('INCR', KEYS[1])
		redis.call('INCR', KEYS[2])
		return -1
	end
	return 0
	`
	keys := []string{quotaKey, roomsKey}
	result, err := l.svcCtx.Redis.Eval(luaScript, keys)
	if err != nil {
		logx.Errorf("Redis Lua 脚本执行失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统繁忙"}, nil
	}
	if val, ok := result.(int64); ok && val == -1 {
		return &types.BaseResponse{Code: 400, Message: "配额不足"}, nil
	}

	// 6. 根据来访目的校验必填字段
	if req.VisitPurpose == 1 { // 求职
		if req.InterviewInfo == nil || req.InterviewInfo.Type == 0 {
			return &types.BaseResponse{Code: 400, Message: "求职必须提供面试证明"}, nil
		}
	} else if req.VisitPurpose == 2 { // 创业
		if req.BusinessPlan == "" {
			return &types.BaseResponse{Code: 400, Message: "创业必须提供创业计划"}, nil
		}
	}

	// 8. 检查用户是否已有未完成的申请（同一用户同时只能有一个待审核或已通过未入住的申请）
	var existingApp model.Application
	err = l.svcCtx.DB.Where("user_id = ? AND status IN (0,1,4)", userId).First(&existingApp).Error
	if err == nil {
		return &types.BaseResponse{Code: 400, Message: "您已有未完成的申请，请先处理"}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.Errorf("查询已有申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 9. 创建申请记录
	now := time.Now().Unix()
	app := &model.Application{
		UserId:        userId,
		StationId:     req.StationId,
		CheckinDate:   req.CheckinDate,
		CheckoutDate:  req.CheckoutDate,
		Status:        0, // 待审核
		VisitPurpose:  int8(req.VisitPurpose),
		Remark:        req.Remark,
		AppliedAt:     now,
		UpdatedAt:     now,
		DepositStatus: 0,
	}
	//加入面试证明相关信息
	if req.InterviewInfo != nil {
		app.InterviewProofType = int8(req.InterviewInfo.Type)
		app.InterviewProofContent = req.InterviewInfo.Content
		app.InterviewProofFileUrl = req.InterviewInfo.FileUrl
	}
	//加入创业证明相关信息
	if req.BusinessPlan != "" {
		app.BusinessPlan = req.BusinessPlan
	}
	// 使用事务：创建申请并扣除驿站配额和可预约房间数
	err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(app).Error; err != nil {
			return err
		}
		// 扣减配额和房间数
		if err := tx.Model(&model.Station{}).Where("id = ?", req.StationId).
			Updates(map[string]interface{}{
				"remaining_quota": station.RemainingQuota - 1,
				"available_rooms": station.AvailableRooms - 1,
			}).Error; err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		logx.Errorf("创建申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "提交失败"}, nil
	}

	return &types.BaseResponse{
		Code:    0,
		Message: "申请提交成功，等待审核",
		Data: &types.ApplyResponse{
			ApplicationId: app.Id,
			Status:        0,
			Message:       "申请提交成功，等待审核",
		},
	}, nil

}
