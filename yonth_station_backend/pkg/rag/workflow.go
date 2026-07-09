package rag

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/cloudwego/eino-ext/components/retriever/milvus2"
	"github.com/cloudwego/eino/schema"
)

type RAGWorkflow struct {
	retriever      *milvus2.Retriever
	model          *deepseek.ChatModel
	topK           int
	scoreThreshold float64
}

func NewRAGWorkflow(retriever *milvus2.Retriever, model *deepseek.ChatModel, topK int, scoreThreshold float64) *RAGWorkflow {
	return &RAGWorkflow{retriever: retriever, model: model, topK: topK, scoreThreshold: scoreThreshold}
}

// filterDocs 过滤文档（自动适配 L2 距离 / COSINE 相似度）
func (w *RAGWorkflow) filterDocs(docs []*schema.Document) []*schema.Document {
	if len(docs) == 0 {
		return nil
	}
	// 如果第一个文档的 score < 1.0，视为 L2 距离（越小越好），否则为 COSINE（越大越好）
	useL2 := docs[0].Score() < 1.0
	filtered := make([]*schema.Document, 0, len(docs))
	for _, d := range docs {
		if useL2 {
			if d.Score() <= w.scoreThreshold { filtered = append(filtered, d) }
		} else {
			if d.Score() >= w.scoreThreshold { filtered = append(filtered, d) }
		}
	}
	return filtered
}

func (w *RAGWorkflow) Query(ctx context.Context, userQuestion string) (string, error) {
	docs, err := w.retriever.Retrieve(ctx, userQuestion)
	if err != nil {
		return "", fmt.Errorf("retrieve failed: %w", err)
	}
	log.Printf("[RAG] got %d docs, threshold=%.2f", len(docs), w.scoreThreshold)
	for i, d := range docs {
		log.Printf("[RAG]   [%d] score=%.4f", i, d.Score())
	}

	docs = w.filterDocs(docs)
	if len(docs) == 0 {
		return "抱歉，知识库中没有找到与您问题相关的信息，建议您联系人工客服。", nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()
	log.Printf("[RAG] Calling DeepSeek with %d docs...", len(docs))
	resp, err := w.model.Generate(ctx, buildMessages(docs, userQuestion, ""))
	if err != nil {
		log.Printf("[RAG] DeepSeek error: %v", err)
		return "", fmt.Errorf("model generate failed: %w", err)
	}
	log.Printf("[RAG] DeepSeek response: %s", resp.Content[:min(80, len(resp.Content))])
	return resp.Content, nil
}

func (w *RAGWorkflow) QueryStream(ctx context.Context, userQuestion string) (*schema.StreamReader[*schema.Message], error) {
	docs, err := w.retriever.Retrieve(ctx, userQuestion)
	if err != nil {
		return nil, fmt.Errorf("retrieve failed: %w", err)
	}
	docs = w.filterDocs(docs)
	if len(docs) == 0 {
		return nil, fmt.Errorf("抱歉，知识库中没有找到与您问题相关的信息，建议您联系人工客服")
	}
	log.Printf("[RAG-Stream] Calling DeepSeek with %d docs...", len(docs))
	return w.model.Stream(ctx, buildMessages(docs, userQuestion, ""))
}

func buildMessages(docs []*schema.Document, question string, history string) []*schema.Message {
	var b strings.Builder
	for i, doc := range docs {
		c := doc.Content
		if len([]rune(c)) > 300 {
			c = string([]rune(c)[:300]) + "..."
		}
		b.WriteString(fmt.Sprintf("[%d] %s\n", i+1, c))
	}
	systemPrompt := "你是南昌市青年驿站的智能客服，请结合知识库和对话历史简洁回答用户问题。"
	hist := ""
	if history != "" {
		hist = fmt.Sprintf("对话历史：\n%s\n", history)
	}
	userPrompt := fmt.Sprintf("知识库：\n%s\n%s问题：%s", b.String(), hist, question)
	return []*schema.Message{
		{Role: schema.System, Content: systemPrompt},
		{Role: schema.User, Content: userPrompt},
	}
}

func (w *RAGWorkflow) QueryStreamWithHistory(ctx context.Context, userQuestion string, history string) (*schema.StreamReader[*schema.Message], error) {
	docs, err := w.retriever.Retrieve(ctx, userQuestion)
	if err != nil {
		return nil, fmt.Errorf("retrieve failed: %w", err)
	}
	docs = w.filterDocs(docs)
	if len(docs) == 0 {
		return nil, fmt.Errorf("抱歉，知识库中没有找到与您问题相关的信息，建议您联系人工客服")
	}
	// 用独立 context 避免 HTTP 请求超时打断流式输出
	// 注意：不能 defer cancel()，否则 stream 还没被消费就被杀掉了
	// 120 秒后 context 自动超时，足够 DeepSeek 完成输出
	streamCtx, _ := context.WithTimeout(context.Background(), 120*time.Second)
	log.Printf("[RAG-Stream] Calling DeepSeek with %d docs, history=%d chars", len(docs), len(history))
	stream, err := w.model.Stream(streamCtx, buildMessages(docs, userQuestion, history))
	if err != nil {
		return nil, fmt.Errorf("stream failed: %w", err)
	}
	return stream, nil
}

func min(a, b int) int { if a < b { return a }; return b }
