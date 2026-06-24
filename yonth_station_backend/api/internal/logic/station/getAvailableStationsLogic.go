// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 获取可预约驿站列表:0.0.0.0:8888/api/v1/station/available
package station

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetAvailableStationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetAvailableStationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetAvailableStationsLogic {
	return &GetAvailableStationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetAvailableStationsLogic) GetAvailableStations() (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	var stations []model.Station
	// 查询条件：运营中(status=1) 且 剩余配额>0 且 可预约房间数>0
	if err := l.svcCtx.DB.Where("status = 1 AND remaining_quota > 0 AND available_rooms > 0").
		Find(&stations).Error; err != nil {
		logx.Errorf("查询可预约驿站失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	list := make([]types.StationInfo, 0, len(stations))
	for _, s := range stations {
		status := int32(s.Status)
		list = append(list, types.StationInfo{
			StationId:      s.Id,
			StationName:    s.StationName,
			District:       s.District,
			Address:        s.Address,
			Latitude:       s.Latitude,
			Longitude:      s.Longitude,
			ContactPhone:   s.ContactPhone,
			BusinessHours:  s.BusinessHours,
			TotalRooms:     int32(s.TotalRooms),
			AvailableRooms: int32(s.AvailableRooms),
			Status:         &status,
			Description:    s.Description,
			Amenities:      s.Amenities,
			NearbyMetro:    s.NearbyMetro,
			ImageUrl:       s.ImageUrl,
		})
	}

	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    list,
	}, nil
}
