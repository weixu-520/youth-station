// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 获取驿站详情:0.0.0.0:8888/api/v1/station/detail/:stationId
package station

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type GetStationDetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStationDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStationDetailLogic {
	return &GetStationDetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStationDetailLogic) GetStationDetail(req *types.GetStationDetailRequest) (resp *types.BaseResponse, err error) {
	cacheKey := fmt.Sprintf("station:detail:%d", req.StationId)
	// 查缓存
	cached, err := utils.GetString(l.ctx, l.svcCtx.Redis, cacheKey)
	if err == nil && cached != "" {
		var detail types.StationDetailResponse
		if err := json.Unmarshal([]byte(cached), &detail); err == nil {
			return &types.BaseResponse{
				Code:    0,
				Message: "success",
				Data:    &detail,
			}, nil
		}
	}
	// 缓存未命中，查询数据库
	var station model.Station
	if err := l.svcCtx.DB.Where("id = ?", req.StationId).First(&station).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 404, Message: "驿站不存在"}, nil
		}
		logx.Errorf("查询驿站失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	// 组装详情响应
	status := int32(station.Status)
	detail := types.StationDetailResponse{
		StationInfo: types.StationInfo{
			StationId:      station.Id,
			StationName:    station.StationName,
			District:       station.District,
			Address:        station.Address,
			Latitude:       station.Latitude,
			Longitude:      station.Longitude,
			ContactPhone:   station.ContactPhone,
			BusinessHours:  station.BusinessHours,
			TotalRooms:     int32(station.TotalRooms),
			AvailableRooms: int32(station.AvailableRooms),
			Status:         &status,
			Description:    station.Description,
			Amenities:      station.Amenities,
			NearbyMetro:    station.NearbyMetro,
			ImageUrl:       station.ImageUrl,
		},
		AverageRating:  station.AvgRating,
		TotalReviews:   int32(station.TotalReviews),
		WeeklyQuota:    int32(station.WeeklyQuota),
		RemainingQuota: int32(station.RemainingQuota),
	}
	// 写入缓存（过期 10 分钟）
	data, _ := json.Marshal(detail)
	_ = utils.SetEx(l.ctx, l.svcCtx.Redis, cacheKey, string(data), 600)

	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    &detail,
	}, nil
}
