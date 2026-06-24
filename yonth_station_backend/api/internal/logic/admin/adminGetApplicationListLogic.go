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

type AdminGetApplicationListLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAdminGetApplicationListLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AdminGetApplicationListLogic {
	return &AdminGetApplicationListLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AdminGetApplicationListLogic) AdminGetApplicationList(req *types.ApplicationListRequest) (resp *types.BaseResponse, err error) {
	//验证管理员权限
	isAdmin, err := utils.GetIsAdmin(l.ctx)
	if err != nil || !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限，仅管理员可操作"}, nil
	}
	//正在营业且存在的驿站,
	query := l.svcCtx.DB.Model(&model.Application{})
	if req.Status != nil {
		query = query.Where("status = ?", req.Status)
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

	//查看当前页数的申请
	offset := (req.Page - 1) * req.PageSize //偏移量
	var apps []model.Application
	if err := query.Offset(int(offset)).Limit(int(req.PageSize)).Order("applied_at DESC").Find(&apps).Error; err != nil {
		logx.Errorf("查询申请列表失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	stationIds := make([]int64, 0)
	userIds := make([]int64, 0)
	for _, app := range apps {
		stationIds = append(stationIds, app.StationId)
		userIds = append(userIds, app.UserId)
	}
	//将驿站id与驿站名进行映射
	stationMap := make(map[int64]string)
	if len(stationIds) > 0 {
		var stations []model.Station
		l.svcCtx.DB.Where("id IN ?", stationIds).Find(&stations)
		for _, s := range stations {
			stationMap[s.Id] = s.StationName
		}
	}
	//将用户id与用户名进行映射
	userMap := make(map[int64]string)
	if len(userIds) > 0 {
		var users []model.User
		l.svcCtx.DB.Select("id, user_name").Where("id IN ?", userIds).Find(&users)
		for _, u := range users {
			userMap[u.Id] = u.UserName
		}
	}
	//整理申请记录
	list := make([]types.ApplicationRecord, 0, len(apps))
	for _, app := range apps {
		status := int32(app.Status)
		list = append(list, types.ApplicationRecord{
			ApplicationId: app.Id,
			UserId:        app.UserId,
			UserName:      userMap[app.UserId],
			StationId:     app.StationId,
			StationName:   stationMap[app.StationId],
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
			},
			List: list,
		},
	}, nil

}
