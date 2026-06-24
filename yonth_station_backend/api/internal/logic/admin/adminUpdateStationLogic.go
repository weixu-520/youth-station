// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"
	"encoding/json"
	"fmt"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminUpdateStationLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminUpdateStationLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminUpdateStationLogic {
	return &AdminUpdateStationLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminUpdateStationLogic) AdminUpdateStation(req *types.StationInfo) (resp *types.BaseResponse, err error) {
	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}
	var station model.Station
	if err := l.svcCtx.DB.Where("id = ?", req.StationId).First(&station).Error; err != nil {
		logx.Errorf("查询驿站失败: %v", err)
		return &types.BaseResponse{Code: 404, Message: "驿站不存在"}, nil
	}
	updates := map[string]interface{}{
		"station_name":    req.StationName,
		"district":        req.District,
		"address":         req.Address,
		"latitude":        req.Latitude,
		"longitude":       req.Longitude,
		"contact_phone":   req.ContactPhone,
		"business_hours":  req.BusinessHours,
		"total_rooms":     req.TotalRooms,
		"available_rooms": req.AvailableRooms,
		"status":          *req.Status,
		"description":     req.Description,
		"amenities":       req.Amenities,
		"nearby_metro":    req.NearbyMetro,
		"image_url":       req.ImageUrl,
	}
	if err := l.svcCtx.DB.Model(&station).Updates(updates).Error; err != nil {
		logx.Errorf("更新驿站失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "更新失败"}, nil
	}
	// 重新加载更新后的驿站数据（确保与数据库一致）
	var updatedStation model.Station
	if err := l.svcCtx.DB.Where("id = ?", req.StationId).First(&updatedStation).Error; err != nil {
		logx.Errorf("重新查询驿站失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "更新成功但缓存更新失败"}, nil
	}

	// 1. 删除旧的详情缓存
	detailCacheKey := fmt.Sprintf("station:detail:%d", req.StationId)
	_ = utils.DeleteCache(l.ctx, l.svcCtx.Redis, detailCacheKey)
	status := int32(updatedStation.Status)

	// 2. 写入新的详情缓存（过期时间 10 分钟）
	detail := types.StationDetailResponse{
		StationInfo: types.StationInfo{
			StationId:      updatedStation.Id,
			StationName:    updatedStation.StationName,
			District:       updatedStation.District,
			Address:        updatedStation.Address,
			Latitude:       updatedStation.Latitude,
			Longitude:      updatedStation.Longitude,
			ContactPhone:   updatedStation.ContactPhone,
			BusinessHours:  updatedStation.BusinessHours,
			TotalRooms:     int32(updatedStation.TotalRooms),
			AvailableRooms: int32(updatedStation.AvailableRooms),
			Status:         &status,
			Description:    updatedStation.Description,
			Amenities:      updatedStation.Amenities,
			NearbyMetro:    updatedStation.NearbyMetro,
			ImageUrl:       updatedStation.ImageUrl,
		},
		AverageRating:  updatedStation.AvgRating,
		TotalReviews:   int32(updatedStation.TotalReviews),
		WeeklyQuota:    int32(updatedStation.WeeklyQuota),
		RemainingQuota: int32(updatedStation.RemainingQuota),
	}
	data, _ := json.Marshal(detail)
	_ = utils.SetEx(l.ctx, l.svcCtx.Redis, detailCacheKey, string(data), 600)
	// 3. 列表缓存不主动删除，依赖其自身的过期时间（如 5 分钟）
	// 因为列表缓存包含多种筛选组合，无法精确清除，让它们自然过期即可。
	return &types.BaseResponse{Code: 0, Message: "更新成功"}, nil
}
