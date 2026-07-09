package rag

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/schema"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// NewIndexer 创建 Milvus 索引器，用于将文档向量化并存入知识库
// apiType: 模型 API 类型，传 nil 使用默认，多模态模型需传 &ark.APITypeMultiModal
func NewIndexer(ctx context.Context, addr, collection, arkAPIKey, arkModel string, apiType *ark.APIType) (*milvus2.Indexer, error) {
	// 1. 创建火山引擎 Ark Embedding 模型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  arkAPIKey,
		Model:   arkModel,
		APIType: apiType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ark embedder: %w", err)
	}

	// 2. 配置 Milvus 索引器（COSINE 与 retriever 保持一致）
	idx, err := milvus2.NewIndexer(ctx, &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: addr,
		},
		Collection: collection,
		Vector:     &milvus2.VectorConfig{Dimension: 2048, VectorField: "vector", MetricType: milvus2.COSINE},
		Embedding:  embedder,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus indexer: %w", err)
	}

	log.Printf("[Indexer] Created successfully for collection: %s", collection)
	return idx, nil
}

// AddDocuments 向知识库添加文档
func AddDocuments(ctx context.Context, idx *milvus2.Indexer, docs []*schema.Document) ([]string, error) {
	ids, err := idx.Store(ctx, docs)
	if err != nil {
		return nil, fmt.Errorf("failed to store documents: %w", err)
	}
	log.Printf("[Indexer] Successfully stored %d documents", len(ids))
	return ids, nil
}
