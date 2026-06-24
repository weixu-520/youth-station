// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package application

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetMyApplicationsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetMyApplicationsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetMyApplicationsLogic {
	return &GetMyApplicationsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetMyApplicationsLogic) GetMyApplications(req *types.ApplicationListRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	query := l.svcCtx.DB.Model(&model.Application{}).Where("user_id = ?", userId)
	// 约定 -1 表示不筛选，0、1、2... 都是有效筛选值
	if req.Status != nil && *req.Status != -1 {
		query = query.Where("status = ?", *req.Status)
	}
	if req.StationId != 0 {
		query = query.Where("station_id = ?", req.StationId)
	}

	//统计申请总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("统计申请总数失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	//查询当前页的数据
	var apps []model.Application
	offset := (req.Page - 1) * req.PageSize
	//结果按申请时间降序排序
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Order("applied_at DESC").Find(&apps).Error; err != nil {
		logx.Errorf("查询申请列表失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	// 获取驿站名称映射
	stationIds := make([]int64, 0)
	for _, app := range apps {
		stationIds = append(stationIds, app.StationId)
	}

	var stations []model.Station
	stationNameMap := make(map[int64]string)
	if len(stationIds) > 0 {
		l.svcCtx.DB.Where("id IN ?", stationIds).Find(&stations)
		for _, s := range stations {
			stationNameMap[s.Id] = s.StationName
		}
	}
	list := make([]types.ApplicationRecord, 0, len(apps))
	for _, app := range apps {
		status := int32(app.Status)
		list = append(list, types.ApplicationRecord{
			ApplicationId: app.Id,
			UserId:        app.UserId,
			UserName:      "", // 本接口不需要用户名，前端自行查或忽略
			StationId:     app.StationId,
			StationName:   stationNameMap[app.StationId],
			CheckinDate:   app.CheckinDate,
			CheckoutDate:  app.CheckoutDate,
			Status:        &status,
			StatusDesc:    utils.GetStatusDesc(int8(app.Status)),
			VisitPurpose:  int32(app.VisitPurpose),
			RejectReason:  app.RejectReason,
			DepositAmount: app.DepositAmount,
			DepositStatus: int32(app.DepositStatus),
			CheckinAt:     app.CheckinAt,
			CheckoutAt:    app.CheckoutAt,
			AppliedAt:     app.AppliedAt,
			UpdatedAt:     app.UpdatedAt,
		})
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data: &types.ApplicationListResponse{
			PageResponse: types.PageResponse{
				Total:    total,
				Page:     req.Page,
				PageSize: req.PageSize,
				// 不要设置 List 字段
			},
			List: list, // 外层 List 字段
		},
	}, nil
}
