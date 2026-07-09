package rag

import (
	"context"
	"fmt"
	"log"

	"github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// NewRetriever 创建 Milvus 检索器，支持相似度阈值过滤
func NewRetriever(ctx context.Context, addr, collection, arkAPIKey, arkModel string, topK int, scoreThreshold float64, apiType *ark.APIType) (*milvus2.Retriever, error) {
	// 1. 创建 Ark Embedding 模型
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  arkAPIKey,
		Model:   arkModel,
		APIType: apiType,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create ark embedder: %w", err)
	}

	// 2. 配置检索器
	ret, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: addr,
		},
		Collection: collection,
		TopK:       topK,
		SearchMode: search_mode.NewApproximate(milvus2.COSINE), // 余弦相似度
		Embedding:  embedder,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create milvus retriever: %w", err)
	}

	if apiType != nil {
		log.Printf("[Retriever] Created, threshold=%.2f, apiType=%s", scoreThreshold, string(*apiType))
	} else {
		log.Printf("[Retriever] Created, threshold=%.2f, apiType=nil", scoreThreshold)
	}
	return ret, nil
}
