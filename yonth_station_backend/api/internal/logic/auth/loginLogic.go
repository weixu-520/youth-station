// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1/auth/login
package auth

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"strconv"
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
	// ========== 1. 先查询用户（用于确定 failKey 并复用） ==========
	var user model.User
	isPhone := regexp.MustCompile(`^1[0-9]{10}$`).MatchString(req.Account)
	var queryErr error
	userExists := false

	if isPhone {
		queryErr = l.svcCtx.DB.Where("phone = ?", req.Account).First(&user).Error
		if queryErr == nil {
			userExists = true
		}
	} else {
		queryErr = l.svcCtx.DB.Where("user_name = ?", req.Account).First(&user).Error
		if queryErr == nil {
			userExists = true
		}
	}

	// 处理数据库错误（除记录不存在外）
	if queryErr != nil && !errors.Is(queryErr, gorm.ErrRecordNotFound) {
		logx.Errorf("查询用户失败: %v", queryErr)
		return &types.BaseResponse{Code: 500, Message: "数据库查询错误"}, nil
	}

	// ========== 2. 确定失败计数的 key ==========
	var failKey string
	if userExists {
		// 用户存在，使用手机号作为统一 key（手机号唯一且不变）
		failKey = fmt.Sprintf("login:fail:%s", user.Phone)
	} else {
		// 用户不存在，使用输入的 account 作为 key（只能限制该输入值）
		failKey = fmt.Sprintf("login:fail:%s", req.Account)
	}

	// ========== 3. 检查失败次数 ==========
	failCountStr, _ := l.svcCtx.Redis.GetCtx(l.ctx, failKey)
	if failCountStr != "" {
		failCount, _ := strconv.Atoi(failCountStr)
		if failCount >= 5 {
			ttl, _ := l.svcCtx.Redis.TtlCtx(l.ctx, failKey)
			return &types.BaseResponse{
				Code:    429,
				Message: fmt.Sprintf("登录失败次数过多，请 %d 秒后再试", ttl),
			}, nil
		}
	}

	// ========== 4. 用户不存在，直接返回并记录失败 ==========
	if !userExists {
		_ = l.incrementFailCount(failKey)
		return &types.BaseResponse{Code: 400, Message: "用户不存在"}, nil
	}

	// ========== 5. 校验密码 ==========
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		_ = l.incrementFailCount(failKey)
		return &types.BaseResponse{Code: 400, Message: "密码错误"}, nil
	}

	// ========== 6. 登录成功，清除失败计数 ==========
	_, _ = l.svcCtx.Redis.DelCtx(l.ctx, failKey)

	// 检查账号状态
	if user.Status == 1 {
		return &types.BaseResponse{Code: 403, Message: "账号已被冻结"}, nil
	}

	// 更新最后登录时间
	_ = l.svcCtx.DB.Model(&user).Update("last_login_at", time.Now().Unix()).Error

	// 生成 JWT token
	token, expireAt, err := utils.GenerateToken(user.Id, user.IsAdmin, l.svcCtx.Config.Auth.AccessSecret, l.svcCtx.Config.Auth.AccessExpire)
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
			IsAdmin:   user.IsAdmin,
		},
	}, nil
}

// incrementFailCount 增加失败计数，返回当前失败次数
func (l *LoginLogic) incrementFailCount(key string) error {
	newCount, err := l.svcCtx.Redis.IncrCtx(l.ctx, key)
	if err != nil {
		logx.Errorf("增加失败计数失败: %v", err)
		return err
	}
	// 第一次失败时设置过期时间为 15 分钟
	if newCount == 1 {
		_ = l.svcCtx.Redis.ExpireCtx(l.ctx, key, 900)
	}
	return nil
}
