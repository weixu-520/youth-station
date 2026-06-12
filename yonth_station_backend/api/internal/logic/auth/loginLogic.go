// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1/auth/login
package auth

import (
	"context"
	"errors"
	"regexp"
	"time"
	model "yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/gorm"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginRequest) (resp *types.BaseResponse, err error) {
	// todo: add your logic here and delete this line
	var user model.User
	// 判断 account 是手机号还是用户名
	isPhone := regexp.MustCompile(`^1[0-9]{10}$`).MatchString(req.Account)
	if isPhone {
		// 按手机号查询（注意 phone 字段为 NULL 的记录不会匹配）
		err = l.svcCtx.DB.Where("phone = ?", req.Account).First(&user).Error
	} else {
		// 按用户名查询
		err = l.svcCtx.DB.Where("user_name = ?", req.Account).First(&user).Error
	}
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return &types.BaseResponse{Code: 400, Message: "用户不存在"}, nil
		}
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "数据库查询错误"}, nil
	}

	// 校验密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		return &types.BaseResponse{Code: 400, Message: "密码错误"}, nil
	}

	// 检查账号状态
	if user.Status == 1 {
		return &types.BaseResponse{Code: 403, Message: "账号已被冻结"}, nil
	}

	// 更新最后登录时间
	_ = l.svcCtx.DB.Model(&user).Update("last_login_at", time.Now().Unix()).Error

	// 生成 JWT token
	token, expireAt, err := utils.GenerateToken(user.Id, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
	if err != nil {
		logx.Errorf("生成token失败: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "登录成功",
		Data: &types.LoginResponse{
			Token:     token,
			ExpiresAt: expireAt,
			UserId:    user.Id,
			UserName:  user.UserName,
		},
	}, nil
}
