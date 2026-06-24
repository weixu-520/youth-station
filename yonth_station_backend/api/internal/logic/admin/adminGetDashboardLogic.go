// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"encoding/json"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetDashboardLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetDashboardLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetDashboardLogic {
	return &AdminGetDashboardLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetDashboardLogic) AdminGetDashboard() (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}
	//查缓存
	cacheKey := "dashboard:stats"
	cached, err := utils.GetString(l.ctx, l.svcCtx.Redis, cacheKey)
	if err == nil && cached != "" {
		var data map[string]interface{}
		if err := json.Unmarshal([]byte(cached), &data); err == nil {
			return &types.BaseResponse{Code: 0, Message: "success", Data: data}, nil
		}
	}

	// 计算统计数据（与之前相同）
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()).Unix()

	var totalApps, todayApps, pendingApps, checkedInApps int64
	var totalStations, activeStations int64

	l.svcCtx.DB.Model(&model.Application{}).Count(&totalApps)
	l.svcCtx.DB.Model(&model.Application{}).Where("applied_at >= ?", todayStart).Count(&todayApps)
	l.svcCtx.DB.Model(&model.Application{}).Where("status = 0").Count(&pendingApps)
	l.svcCtx.DB.Model(&model.Application{}).Where("status = 4").Count(&checkedInApps)

	l.svcCtx.DB.Model(&model.Station{}).Count(&totalStations)
	l.svcCtx.DB.Model(&model.Station{}).Where("status = 1").Count(&activeStations)

	data := map[string]interface{}{
		"totalApplications":     totalApps,
		"todayApplications":     todayApps,
		"pendingApplications":   pendingApps,
		"checkedInApplications": checkedInApps,
		"totalStations":         totalStations,
		"activeStations":        activeStations,
	}

	// 写入缓存，5分钟过期
	jsonData, _ := json.Marshal(data)
	_ = utils.SetEx(l.ctx, l.svcCtx.Redis, cacheKey, string(jsonData), 300)

	return &types.BaseResponse{Code: 0, Message: "success", Data: data}, nil
}
