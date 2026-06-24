package chat

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetHistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetHistoryLogic {
	return &GetHistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetHistoryLogic) GetHistory() (resp *types.BaseResponse, err error) {
	userId, _ := utils.GetUserId(l.ctx)
	isAdmin, _ := utils.GetIsAdmin(l.ctx)

	var msgs []model.ChatMessage
	query := l.svcCtx.DB.Order("created_at ASC").Limit(200)

	if isAdmin {
		// 管理员：拉取所有消息
		query.Find(&msgs)
	} else {
		// 普通用户：拉取与自己相关的消息
		query.Where("from_user_id = ? OR to_user_id = ?", userId, userId).Find(&msgs)
	}

	list := make([]map[string]interface{}, 0, len(msgs))
	for _, m := range msgs {
		list = append(list, map[string]interface{}{
			"fromUserId": m.FromUserId,
			"toUserId":   m.ToUserId,
			"targetType": m.TargetType,
			"content":    m.Content,
			"createdAt":  m.CreatedAt.Unix(),
		})
	}

	return &types.BaseResponse{Code: 0, Message: "success", Data: list}, nil
}
