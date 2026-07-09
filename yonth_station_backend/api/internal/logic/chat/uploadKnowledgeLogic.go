package chat

import (
	"context"
	"fmt"
	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/types"
	"yonth_station_backend/pkg/rag"

	"github.com/cloudwego/eino/schema"
	"github.com/zeromicro/go-zero/core/logx"
)

type UploadKnowledgeLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadKnowledgeLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadKnowledgeLogic {
	return &UploadKnowledgeLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadKnowledgeLogic) UploadKnowledge(req *types.KnowledgeUploadRequest) (resp *types.BaseResponse, err error) {
	// 1. 保存文档元数据到 MySQL（先保存，获取 ID）
	doc := &model.KnowledgeDoc{
		Title:    req.Title,
		Content:  req.Content,
		Category: req.Category,
		Status:   0, // 待索引
	}
	if err := l.svcCtx.DB.Create(doc).Error; err != nil {
		return &types.BaseResponse{Code: 500, Message: "保存文档失败"}, nil
	}

	// 2. 构建 Eino Document 对象
	einoDoc := &schema.Document{
		ID:      fmt.Sprintf("%d", doc.Id),
		Content: req.Content,
		MetaData: map[string]interface{}{
			"title":    req.Title,
			"category": req.Category,
		},
	}

	// 3. 调用 Indexer 向量化并存入 Milvus
	ids, err := rag.AddDocuments(l.ctx, l.svcCtx.Indexer, []*schema.Document{einoDoc})
	if err != nil {
		logx.Errorf("Index failed: %v", err)
		// 失败可回滚或标记状态，此处简单处理
		return &types.BaseResponse{Code: 500, Message: "索引失败"}, nil
	}

	// 4. 更新文档状态为已索引
	l.svcCtx.DB.Model(doc).Update("status", 1)

	return &types.BaseResponse{
		Code:    0,
		Message: "上传成功",
		Data: &types.KnowledgeUploadResponse{
			DocId:   ids[0],
			Message: "文档已上传并索引",
		},
	}, nil
}
