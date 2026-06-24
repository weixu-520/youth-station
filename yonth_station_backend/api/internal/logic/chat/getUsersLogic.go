package chat

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUsersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUsersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUsersLogic {
	return &GetUsersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUsersLogic) GetUsers() (resp *types.BaseResponse, err error) {
	isAdmin, _ := utils.GetIsAdmin(l.ctx)
	if !isAdmin {
		return &types.BaseResponse{Code: 403, Message: "无权限"}, nil
	}

	// 查询所有发过消息的用户的 ID（去重）
	var userIds []int64
	l.svcCtx.DB.Model(&model.ChatMessage{}).
		Select("DISTINCT from_user_id").
		Where("target_type = 1"). // 用户→管理员
		Pluck("from_user_id", &userIds)

	if len(userIds) == 0 {
		return &types.BaseResponse{Code: 0, Data: []interface{}{}}, nil
	}

	// 批量查用户名称
	var users []model.User
	l.svcCtx.DB.Select("id, user_name").Where("id IN ?", userIds).Find(&users)

	type chatUser struct {
		UserId   int64  `json:"userId"`
		UserName string `json:"userName"`
	}
	list := make([]chatUser, 0, len(users))
	for _, u := range users {
		list = append(list, chatUser{UserId: u.Id, UserName: u.UserName})
	}

	return &types.BaseResponse{Code: 0, Data: list}, nil
}
