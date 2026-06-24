// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1/user/info
package user

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

type GetUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetUserInfoLogic) GetUserInfo() (resp *types.BaseResponse, err error) {
	//从jwt中获取用户id
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	//查缓存
	cacheKey := fmt.Sprintf("user:info:%d", userId)
	cached, err := utils.GetString(l.ctx, l.svcCtx.Redis, cacheKey)
	if err == nil && cached != "" {
		var userInfo types.UserInfoResponse
		if err := json.Unmarshal([]byte(cached), &userInfo); err == nil {
			return &types.BaseResponse{Code: 0, Message: "success", Data: &userInfo}, nil
		}
	}

	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{Code: 404, Message: "用户不存在"}, nil
	}

	// 脱敏处理
	maskedPhone := utils.MaskPhone(user.Phone)
	maskedIdCard := utils.MaskIdCard(user.IdCard)

	userInfo := types.UserInfoResponse{
		UserId:       user.Id,
		UserName:     user.UserName,
		Phone:        maskedPhone,
		IdCard:       maskedIdCard,
		BirthDate:    user.BirthDate,
		Gender:       int32(user.Gender),
		Education:    int32(user.Education),
		School:       user.School,
		GraduateYear: int32(user.GraduateYear),
		HukouCity:    user.HukouCity,
		Status:       int32(user.Status),
		IsAdmin:      user.IsAdmin,
		CreatedAt:    user.CreatedAt.Unix(),
	}

	// 写入缓存，10分钟过期
	data, _ := json.Marshal(userInfo)
	_ = utils.SetEx(l.ctx, l.svcCtx.Redis, cacheKey, string(data), 600)

	return &types.BaseResponse{Code: 0, Message: "success", Data: &userInfo}, nil
}
