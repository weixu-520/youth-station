// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2
// 访问路径：0.0.0.0:8888/api/v1/user/info
package user

import (
	"context"
	"encoding/json"
	"yonth_station_backend/api/gorm/model"

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
	// todo: add your logic here and delete this line
	//1. 从上下文中获取 userId（由 JWT 中间件注入）
	// go-zero JWT 解析使用 jwt.WithJSONNumber()，数字类型是 json.Number
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return &types.BaseResponse{Code: 401, Message: "未登录或登录已过期"}, nil
	}
	userId, err := userIdVal.(json.Number).Int64()
	if err != nil || userId == 0 {
		return &types.BaseResponse{
			Code:    401,
			Message: "未登录或登录已过期",
		}, nil
	}
	// 2. 查询用户
	var user model.User
	if err := l.svcCtx.DB.Where("id = ?", userId).First(&user).Error; err != nil {
		logx.Errorf("查询用户失败: %v", err)
		return &types.BaseResponse{
			Code:    404,
			Message: "用户不存在",
		}, nil
	}
	// 3. 脱敏处理
	// 手机号脱敏：138****8000
	phoneStr := ""
	if user.Phone != nil && *user.Phone != "" {
		phone := *user.Phone
		if len(phone) == 11 {
			phoneStr = phone[:3] + "****" + phone[7:]
		} else {
			phoneStr = phone // 非标准长度不脱敏
		}
	}
	// 身份证号脱敏：110101******1234 (保留前6后4，中间8位星号)
	idCardStr := ""
	if len(user.IdCard) >= 18 {
		idCardStr = user.IdCard[:6] + "********" + user.IdCard[14:]
	} else if user.IdCard != "" {
		idCardStr = "****"
	}
	return &types.BaseResponse{
		Code:    0,
		Message: "success",
		Data: &types.UserInfoResponse{
			UserId:       user.Id,
			UserName:     user.UserName,
			Phone:        phoneStr,
			IdCard:       idCardStr,
			BirthDate:    user.BirthDate,
			Gender:       int32(user.Gender),
			Education:    int32(user.Education),
			School:       user.School,
			GraduateYear: int32(user.GraduateYear),
			HukouCity:    user.HukouCity,
			Status:       int32(user.Status),
			CreatedAt:    user.CreatedAt.Unix(),
		},
	}, nil
}
