// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

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

type AuditApplicationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAuditApplicationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AuditApplicationLogic {
	return &AuditApplicationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AuditApplicationLogic) AuditApplication(req *types.AuditRequest) (resp *types.BaseResponse, err error) {
	// 1. 检查管理员权限
	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}

	// 2. 获取当前管理员用户信息（用于记录审核人）
	adminId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "无法获取管理员信息"}, nil
	}
	var adminUser model.User
	if err := l.svcCtx.DB.Where("id = ?", adminId).First(&adminUser).Error; err != nil {
		logx.Errorf("查询管理员信息失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	adminName := adminUser.UserName // 使用用户名作为审核人

	// 3. 查询申请记录
	var app model.Application
	if err := l.svcCtx.DB.Where("id = ?", req.ApplicationId).First(&app).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 404, Message: "申请不存在"}, nil
		}
		logx.Errorf("查询申请失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	if app.Status != 0 {
		return &types.BaseResponse{Code: 400, Message: "当前状态无法审核"}, nil
	}

	now := time.Now().Unix()
	//拒绝申请的情况
	if req.Result == 2 {
		if req.RejectReason == "" {
			return &types.BaseResponse{Code: 400, Message: "拒绝原因不能为空"}, nil
		}
		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&app).Updates(map[string]interface{}{
				"status":        2,
				"reject_reason": req.RejectReason,
				"audit_by":      adminName,
				"audit_at":      now,
				"updated_at":    now,
			}).Error; err != nil {
				return err
			}
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
			logx.Errorf("审核拒绝失败: %v", err)
			return &types.BaseResponse{Code: 500, Message: "操作失败"}, nil
		}
		return &types.BaseResponse{Code: 0, Message: "已拒绝申请"}, nil
	}

	// 5. 处理通过
	if req.Result == 1 {
		// 分配房间
		var room model.Room
		if err := l.svcCtx.DB.Where("station_id = ? AND status = 0", app.StationId).First(&room).Error; err != nil {
			logx.Errorf("分配房间失败: %v", err)
			return &types.BaseResponse{Code: 500, Message: "暂无可用房间"}, nil
		}
		// 设置押金金额（示例 200 元）
		depositAmount := int64(20000)

		err = l.svcCtx.DB.Transaction(func(tx *gorm.DB) error {
			if err := tx.Model(&app).Updates(map[string]interface{}{
				"status":         1,
				"room_id":        room.Id,
				"deposit_amount": depositAmount,
				"audit_by":       adminName,
				"audit_at":       now,
				"updated_at":     now,
			}).Error; err != nil {
				return err
			}
			if err := tx.Model(&room).Update("status", 1).Error; err != nil {
				return err
			}
			return nil
		})
		if err != nil {
			logx.Errorf("审核通过失败: %v", err)
			return &types.BaseResponse{Code: 500, Message: "操作失败"}, nil
		}
		return &types.BaseResponse{Code: 0, Message: "已通过申请", Data: map[string]interface{}{"depositAmount": depositAmount}}, nil
	}
	return &types.BaseResponse{Code: 400, Message: "无效的审核结果"}, nil
}
