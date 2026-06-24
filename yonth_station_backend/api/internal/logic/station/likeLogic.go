// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"context"
	"fmt"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type LikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LikeLogic {
	return &LikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LikeLogic) Like(req *types.LikeRequest) (resp *types.BaseResponse, err error) {
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录"}, nil
	}

	stationId := req.StationId
	userKey := fmt.Sprintf("station:like:user:%d", stationId)
	countKey := fmt.Sprintf("station:like:count:%d", stationId)

	// Lua 脚本原子操作
	script := `
        if redis.call('SISMEMBER', KEYS[1], ARGV[1]) == 1 then
            return -1
        end
        redis.call('SADD', KEYS[1], ARGV[1])
        return redis.call('INCR', KEYS[2])
    `

	result, err := l.svcCtx.Redis.EvalCtx(l.ctx, script, []string{userKey, countKey}, userId)

	if err != nil {
		logx.Errorf("Redis Eval error: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	newCount, _ := result.(int64)
	if newCount == -1 {
		return &types.BaseResponse{Code: 400, Message: "已点赞"}, nil
	}

	// 异步持久化到 MySQL
	go func() {
		like := &model.StationLike{UserId: userId, StationId: stationId}
		_ = l.svcCtx.DB.Create(like).Error
	}()

	return &types.BaseResponse{
		Code:    0,
		Message: "点赞成功",
		Data:    &types.LikeResponse{Count: newCount},
	}, nil
}
