// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"context"
	"fmt"
	"strconv"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetLikeCountLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetLikeCountLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetLikeCountLogic {
	return &GetLikeCountLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetLikeCountLogic) GetLikeCount(req *types.GetLikeCountRequest) (resp *types.BaseResponse, err error) {
	stationId := req.StationId
	key := fmt.Sprintf("station:like:count:%d", stationId)
	val, err := l.svcCtx.Redis.GetCtx(l.ctx, key)
	var count int64
	if err != nil || val == "" {
		// 从 MySQL 查询
		var cnt int64
		l.svcCtx.DB.Model(&model.StationLike{}).Where("station_id = ?", stationId).Count(&cnt)
		count = cnt
		// 用 Setnx 避免覆盖并发的 INCR（like/unlike 的 Lua 脚本可能已在 Redis 写了新值）
		if _, err := l.svcCtx.Redis.SetnxExCtx(l.ctx, key, strconv.FormatInt(cnt, 10), 3600); err != nil {
			// Setnx 失败说明 key 已存在（有并发写入），重新读取
			if v, e := l.svcCtx.Redis.GetCtx(l.ctx, key); e == nil && v != "" {
				count, _ = strconv.ParseInt(v, 10, 64)
			}
		}
	} else {
		count, _ = strconv.ParseInt(val, 10, 64)
	}
	return &types.BaseResponse{Code: 0, Data: &types.LikeResponse{Count: count}}, nil
}
