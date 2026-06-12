// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1//auth/register
package auth

import (
	"context"
	"errors"

	model "yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	// 1. 检查用户名是否已存在
	var existUser model.User
	if err := l.svcCtx.DB.Where("user_name = ?", req.UserName).First(&existUser).Error; err == nil {
		return &types.BaseResponse{Code: 400, Message: "用户名已存在"}, nil
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "数据库查询错误"}, nil
	}

	//2.如果填写了手机号，检查手机号是否已存在
	var phoneToStore *string = nil //存入数据库时使用
	if req.Phone != "" {
		var existPhone model.User
		if err := l.svcCtx.DB.Where("phone = ?", req.Phone).First(&existPhone).Error; err == nil {
			return &types.BaseResponse{Code: 400, Message: "手机号已存在"}, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logx.Errorf("查询手机号失败: %v", err)
			return &types.BaseResponse{Code: 500, Message: "数据库查询错误"}, nil
		}
	}
	phoneToStore = &req.Phone

	//3. 密码加密
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		logx.Errorf("密码加密失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "密码加密错误"}, nil
	}
	//4.创建用户
	user := &model.User{
		UserName: req.UserName,
		Password: hashedPassword,
		Status:   0,
		Phone:    phoneToStore, // 若 phoneToStore 为 nil，则数据库存储 NULL
	}
	if err := l.svcCtx.DB.Create(user).Error; err != nil {
		logx.Errorf("创建用户失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "注册失败"}, nil
	}
	//5.生成JWT token
	token, expireAt, err := utils.GenerateToken(user.Id, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "注册成功",
		Data: &types.RegisterResponse{
			Token:     token,
			ExpiresAt: expireAt,
			UserId:    user.Id,
			UserName:  user.UserName,
		},
	}, nil

}
