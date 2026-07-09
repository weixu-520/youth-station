package chat

import (
	"context"
	"fmt"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
)

// NewDeepSeekModel 创建 DeepSeek 聊天模型
func NewDeepSeekModel(ctx context.Context, apiKey, model, baseURL string) (*deepseek.ChatModel, error) {
	cm, err := deepseek.NewChatModel(ctx, &deepseek.ChatModelConfig{
		APIKey:  apiKey,
		Model:   model,
		BaseURL: baseURL,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create DeepSeek model: %w", err)
	}
	return cm, nil
}
