package main

import (
	"context"
	"fmt"
	"log"

	"yonth_station_backend/api/internal/config"

	ark "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2/search_mode"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
	"github.com/zeromicro/go-zero/core/conf"
)

func main() {
	var c config.Config
	conf.MustLoad("../../etc/gateway-api.yaml", &c)

	ctx := context.Background()

	apiType := ark.APITypeMultiModal
	embedder, err := ark.NewEmbedder(ctx, &ark.EmbeddingConfig{
		APIKey:  c.Chat.Embedding.APIKey,
		Model:   c.Chat.Embedding.Model,
		APIType: &apiType,
	})
	if err != nil {
		log.Fatalf("embedder: %v", err)
	}

	ret, err := milvus2.NewRetriever(ctx, &milvus2.RetrieverConfig{
		ClientConfig: &milvusclient.ClientConfig{Address: c.Chat.Milvus.Address},
		Collection:   c.Chat.Milvus.Collection,
		TopK:         5,
		SearchMode:   search_mode.NewApproximate(milvus2.COSINE),
		Embedding:    embedder,
	})
	if err != nil {
		log.Fatalf("retriever: %v", err)
	}

	questions := []string{"如何申请入住", "押金怎么退", "你好"}
	for _, q := range questions {
		docs, err := ret.Retrieve(ctx, q)
		if err != nil {
			log.Printf("[%s] error: %v", q, err)
			continue
		}
		fmt.Printf("\n=== 问题: %s ===\n", q)
		fmt.Printf("返回 %d 条结果:\n", len(docs))
		for i, d := range docs {
			content := d.Content
			if len(content) > 60 {
				content = content[:60] + "..."
			}
			fmt.Printf("  [%d] score=%.4f  title=%s  content=%s\n", i+1, d.Score(), d.MetaData["title"], content)
		}
	}
}
