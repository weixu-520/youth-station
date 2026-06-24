// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 获取驿站列表
// 访问路径：0.0.0.0:8888/api/v1/station/list
package station

import (
	"context"
	"encoding/json"
	"fmt"
	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/pkg/utils"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetStationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStationListLogic {
	return &GetStationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetStationListLogic) GetStationList(req *types.StationListRequest) (resp *types.BaseResponse, err error) {
	// 构建缓存 key
	cacheKey := fmt.Sprintf("station:list:%d:%d:%s:%d", req.Page, req.PageSize, req.District, req.Status)
	// 尝试从缓存读取
	cached, err := utils.GetString(l.ctx, l.svcCtx.Redis, cacheKey)
	if err == nil && cached != "" {
		var result types.StationListResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &types.BaseResponse{
				Code:    0,
				Message: "success",
				Data:    &result,
			}, nil
		}
	}

	// 缓存未命中，查询数据库
	query := l.svcCtx.DB.Model(&model.Station{})
	//判断所填区域是否为空
	if req.District != "" {
		query = query.Where("district = ?", req.District)
	}
	//确定驿站未关闭，处于运营状态
	if req.Status != 0 {
		query = query.Where("status = ?", req.Status)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("统计驿站总数失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 分页查询
	var stations []model.Station
	offset := (req.Page - 1) * req.PageSize
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Find(&stations).Error; err != nil {
		logx.Errorf("查询驿站列表失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 转换为响应格式
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

	result := &types.StationListResponse{
		PageResponse: types.PageResponse{
			Total:    total,
			Page:     req.Page,
			PageSize: req.PageSize,
		},
		List: list,
	}
	// 写入缓存（过期时间 5 分钟）
	data, _ := json.Marshal(result)
	_ = utils.SetEx(l.ctx, l.svcCtx.Redis, cacheKey, string(data), 60)

	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data:    result,
	}, nil

}
