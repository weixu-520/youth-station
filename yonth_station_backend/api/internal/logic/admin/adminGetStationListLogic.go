// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package admin

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type AdminGetStationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetStationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetStationListLogic {
	return &AdminGetStationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetStationListLogic) AdminGetStationList(req *types.StationListRequest) (resp *types.BaseResponse, err error) {

	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}

	query := l.svcCtx.DB.Model(&model.Station{})
	if req.District != "" {
		query = query.Where("district = ?", req.District)
	}
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("统计驿站总数失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	offset := (req.Page - 1) * req.PageSize
	var stations []model.Station
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&stations).Error; err != nil {
		logx.Errorf("查询驿站列表失败: %v", err)
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
		Data: &types.StationListResponse{
			PageResponse: types.PageResponse{
				Total:    total,
				Page:     req.Page,
				PageSize: req.PageSize,
			},
			List: list,
		},
	}, nil
}
