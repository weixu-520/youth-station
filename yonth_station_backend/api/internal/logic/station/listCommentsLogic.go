// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package station

import (
	"context"
	"encoding/json"
	"fmt"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCommentsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCommentsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCommentsLogic {
	return &ListCommentsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCommentsLogic) ListComments(req *types.CommentListRequest) (resp *types.BaseResponse, err error) {
	cacheKey := fmt.Sprintf("station:comments:%d:%d:%d", req.StationId, req.Page, req.PageSize)
	cached, _ := l.svcCtx.Redis.GetCtx(l.ctx, cacheKey)
	if cached != "" {
		var result types.CommentListResponse
		if err := json.Unmarshal([]byte(cached), &result); err == nil {
			return &types.BaseResponse{Code: 0, Data: &result}, nil
		}
	}

	// 1. 查询评论总数
	var total int64
	query := l.svcCtx.DB.Model(&model.StationComment{}).Where("station_id = ?", req.StationId)
	if err := query.Count(&total).Error; err != nil {
		logx.Errorf("Count comments error: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 2. 分页查询评论列表
	var comments []model.StationComment
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("created_at DESC").Offset(int(offset)).Limit(int(req.PageSize)).Find(&comments).Error; err != nil {
		logx.Errorf("Find comments error: %v", err)
		return &types.BaseResponse{Code: 500, Message: "系统错误"}, nil
	}

	// 3. 如果无评论，直接返回
	if len(comments) == 0 {
		result := types.CommentListResponse{Total: 0, List: []types.CommentItem{}}
		return &types.BaseResponse{Code: 0, Data: &result}, nil
	}

	// 4. 收集所有用户ID
	userIds := make([]int64, 0, len(comments))
	for _, c := range comments {
		userIds = append(userIds, c.UserId)
	}

	// 5. 批量查询用户信息（仅查询需要的字段）
	var users []model.User
	userMap := make(map[int64]string, len(userIds))
	if err := l.svcCtx.DB.Select("id, user_name").Where("id IN ?", userIds).Find(&users).Error; err != nil {
		logx.Errorf("Batch query users error: %v", err)
		// 即使查询用户失败，也继续返回评论（用户名置空）
	} else {
		for _, u := range users {
			userMap[u.Id] = u.UserName
		}
	}

	// 6. 组装响应数据
	list := make([]types.CommentItem, 0, len(comments))
	for _, c := range comments {
		userName := userMap[c.UserId] // 如果没查到，则为空字符串
		list = append(list, types.CommentItem{
			Id:        c.Id,
			UserId:    c.UserId,
			UserName:  userName,
			Content:   c.Content,
			ParentId:  c.ParentId,
			CreatedAt: c.CreatedAt.Unix(),
		})
	}

	result := types.CommentListResponse{
		Total: total,
		List:  list,
	}
	// 7. 写入缓存（5分钟过期）
	data, _ := json.Marshal(result)
	_ = l.svcCtx.Redis.SetexCtx(l.ctx, cacheKey, string(data), 300)

	return &types.BaseResponse{Code: 0, Data: &result}, nil

}
