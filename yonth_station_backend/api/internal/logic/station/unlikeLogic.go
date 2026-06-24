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

type UnlikeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnlikeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlikeLogic {
	return &UnlikeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnlikeLogic) Unlike(req *types.LikeRequest) (resp *types.BaseResponse, err error) {
	userId, _ := utils.GetUserId(l.ctx)
	stationId := req.StationId
	userKey := fmt.Sprintf("station:like:user:%d", stationId)
	countKey := fmt.Sprintf("station:like:count:%d", stationId)

	script := `
        if redis.call('SISMEMBER', KEYS[1], ARGV[1]) == 0 then
            return -1
        end
        redis.call('SREM', KEYS[1], ARGV[1])
        return redis.call('DECR', KEYS[2])
    `
	result, err := l.svcCtx.Redis.EvalCtx(l.ctx, script, []string{userKey, countKey}, userId)
	if err != nil {
		logx.Errorf("Redis Eval error: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}
	newCount, _ := result.(int64)
	if newCount == -1 {
		return &types.BaseResponse{Code: 400, Message: "未点赞"}, nil
	}

	go func() {
		_ = l.svcCtx.DB.Where("user_id = ? AND station_id = ?", userId, stationId).Delete(&model.StationLike{}).Error
	}()
	return &types.BaseResponse{
		Code:    0,
		Message: "已取消点赞",
		Data:    &types.LikeResponse{Count: newCount},
	}, nil
}
