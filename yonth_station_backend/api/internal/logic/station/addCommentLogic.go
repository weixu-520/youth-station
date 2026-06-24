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

type AddCommentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewAddCommentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddCommentLogic {
	return &AddCommentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *AddCommentLogic) AddComment(req *types.AddCommentRequest) (resp *types.BaseResponse, err error) {
	userId, err := utils.GetUserId(l.ctx)
	if err != nil {
		return &types.BaseResponse{Code: 401, Message: "未登录"}, nil
	}
	comment := &model.StationComment{
		UserId:    userId,
		StationId: req.StationId,
		Content:   req.Content,
		ParentId:  req.ParentId,
	}
	if err := l.svcCtx.DB.Create(comment).Error; err != nil {
		logx.Errorf("Create comment error: %v", err)
		return &types.BaseResponse{Code: 500, Message: "发表失败"}, nil
	}
	// 同步删除该驿站的所有评论缓存，避免前端立即刷新时拿到旧缓存
	pattern := fmt.Sprintf("station:comments:%d:*", req.StationId)
	cursor := uint64(0)
	for {
		keys, nextCursor, err := l.svcCtx.Redis.ScanCtx(l.ctx, cursor, pattern, 100)
		if err != nil {
			logx.Errorf("SCAN error: %v", err)
			break
		}
		if len(keys) > 0 {
			if _, err := l.svcCtx.Redis.DelCtx(l.ctx, keys...); err != nil {
				logx.Errorf("DEL error: %v", err)
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}

	return &types.BaseResponse{
		Code:    0,
		Message: "评论成功",
		Data:    &types.AddCommentResponse{CommentId: comment.Id},
	}, nil

}
