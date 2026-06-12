// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1/user/info
package user

import (
	"context"
	"encoding/json"
	"errors"
	"yonth_station_backend/api/gorm/model"

	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type UpdateUserInfoLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateUserInfoLogic) UpdateUserInfo(req *types.UpdateUserInfoRequest) (resp *types.BaseResponse, err error) {
	// 1. 从 context 中获取 userId（go-zero JWT 中间件注入，数字为 json.Number 类型）
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	userIdNum, ok := userIdVal.(json.Number)
	if !ok {
		return &types.BaseResponse{Code: 401, Message: "认证信息异常"}, nil
	}
	userId, err := userIdNum.Int64()
	if err != nil || userId <= 0 {
		return &types.BaseResponse{Code: 401, Message: "认证信息无效"}, nil
	}

	// 2. 查询用户是否存在
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 404, Message: "用户不存在"}, nil
		}
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 3. 构建更新字段（只更新传入的非零值，且与旧值不同的字段）
	updates := make(map[string]interface{})

	if req.UserName != "" && req.UserName != user.UserName {
		// 检查用户名是否已被其他用户使用
		var existUser model.User
		result := l.svcCtx.DB.Where("user_name = ? AND id != ?", req.UserName, userId).First(&existUser)
		if result.Error == nil {
			return &types.BaseResponse{Code: 400, Message: "用户名已被占用"}, nil
		} else if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
			logx.Errorf("查询用户名失败: %v", result.Error)
			return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
		}
		updates["user_name"] = req.UserName
	}
	if req.BirthDate != "" && req.BirthDate != user.BirthDate {
		updates["birth_date"] = req.BirthDate
	}
	if req.Gender != 0 && req.Gender != int32(user.Gender) {
		updates["gender"] = req.Gender
	}
	if req.Education != 0 && req.Education != int32(user.Education) {
		updates["education"] = req.Education
	}
	if req.School != "" && req.School != user.School {
		updates["school"] = req.School
	}
	if req.GraduateYear != 0 && req.GraduateYear != int32(user.GraduateYear) {
		updates["graduate_year"] = req.GraduateYear
	}
	if req.HukouCity != "" && req.HukouCity != user.HukouCity {
		updates["hukou_city"] = req.HukouCity
	}

	// 4. 执行更新
	if len(updates) > 0 {
		if err := l.svcCtx.DB.Model(&user).Updates(updates).Error; err != nil {
			logx.Errorf("更新用户信息失败: %v", err)
			return &types.BaseResponse{Code: 500, Message: "更新失败"}, nil
		}
	}

	// 5. 重新查询，返回更新后的完整用户信息
	var updatedUser model.User
	l.svcCtx.DB.First(&updatedUser, userId)

	return &types.BaseResponse{
		Code:    0,
		Message: "更新成功",
		Data: &types.UserInfoResponse{
			UserId:       updatedUser.Id,
			UserName:     updatedUser.UserName,
			Phone:        derefString(updatedUser.Phone),
			IdCard:       updatedUser.IdCard,
			BirthDate:    updatedUser.BirthDate,
			Gender:       int32(updatedUser.Gender),
			Education:    int32(updatedUser.Education),
			School:       updatedUser.School,
			GraduateYear: int32(updatedUser.GraduateYear),
			HukouCity:    updatedUser.HukouCity,
			Status:       int32(updatedUser.Status),
			CreatedAt:    updatedUser.CreatedAt.Unix(),
		},
	}, nil
}

func derefString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}
