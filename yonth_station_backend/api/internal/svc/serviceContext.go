// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"context"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/config"
	"yonth_station_backend/api/internal/websocket"
	"yonth_station_backend/pkg/chat"
	"yonth_station_backend/pkg/logger"
	"yonth_station_backend/pkg/rabbitmq"
	"yonth_station_backend/pkg/rag"

	ark "github.com/cloudwego/eino-ext/components/embedding/ark"
	"github.com/cloudwego/eino-ext/components/model/deepseek"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	milvusIndexer "github.com/cloudwego/eino-ext/components/indexer/milvus2"
	milvusRetriever "github.com/cloudwego/eino-ext/components/retriever/milvus2"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	DB             *gorm.DB
	KafkaLogWriter *logger.KafkaLogWriter
	RabbitMQ       *rabbitmq.RabbitMQClient
	Hub            *websocket.Hub
	Indexer        *milvusIndexer.Indexer
	Retriever      *milvusRetriever.Retriever
	ChatModel      *deepseek.ChatModel
	RAGWorkflow    *rag.RAGWorkflow
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接MySQL
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Errorf("failed to connect database: %v", err)
		panic(err)
	}

	//自动迁移表结构（仅添加缺失列，不删除已有列/数据）
	if err := db.AutoMigrate(&model.User{}, &model.Station{}, &model.Room{}, &model.Application{}, &model.Payment{}, &model.ChatMessage{}, &model.StationComment{}, &model.StationLike{}, &model.KnowledgeDoc{}, &model.ChatHistory{}); err != nil {
		logx.Errorf("failed to auto migrate: %v", err)
		panic(err)
	}

	//连接redis
	rds := redis.MustNewRedis(c.Redis)

	// 初始化 RabbitMQ
	mqClient, err := rabbitmq.NewRabbitMQClient(
		c.RabbitMQ.Host,
		c.RabbitMQ.Port,
		c.RabbitMQ.User,
		c.RabbitMQ.Password,
		c.RabbitMQ.Exchange,
	)
	if err != nil {
		logx.Errorf("RabbitMQ init error: %v", err)
		// 可 panic 或降级，此处先记录错误，继续运行（但 WebSocket 广播会失效）
	}

	// 创建 Hub 并启动
	hub := websocket.NewHub(db, mqClient)
	go hub.Run()

	ctx := context.Background()

	// 根据配置决定 API 类型（多模态模型需用 APITypeMultiModal）
	// 使用多模态 API 类型（火山引擎 Ark 需此参数）
	t := ark.APITypeMultiModal
	apiType := &t

	// 1. 创建 DeepSeek 模型
	chatModel, err := chat.NewDeepSeekModel(ctx, c.Chat.DeepSeek.APIKey, c.Chat.DeepSeek.Model, c.Chat.DeepSeek.BaseURL)
	if err != nil {
		logx.Errorf("DeepSeek init error: %v", err)
	}

	// 2. 创建 Indexer（写入）
	indexer, err := rag.NewIndexer(ctx, c.Chat.Milvus.Address, c.Chat.Milvus.Collection, c.Chat.Embedding.APIKey, c.Chat.Embedding.Model, apiType)
	if err != nil {
		logx.Errorf("Indexer init error: %v", err)
	}

	// 3. 创建 Retriever（读取）
	retriever, err := rag.NewRetriever(
		ctx,
		c.Chat.Milvus.Address,
		c.Chat.Milvus.Collection,
		c.Chat.Embedding.APIKey,
		c.Chat.Embedding.Model,
		c.Chat.RAG.TopK,
		c.Chat.RAG.ScoreThreshold,
		apiType,
	)
	if err != nil {
		logx.Errorf("Retriever init error: %v", err)
	}
	ragWorkflow := rag.NewRAGWorkflow(retriever, chatModel, c.Chat.RAG.TopK, c.Chat.RAG.ScoreThreshold)

	ctxSer := &ServiceContext{
		Config:      c,
		DB:          db,
		Redis:       rds,
		RabbitMQ:    mqClient,
		Hub:         hub,
		ChatModel:   chatModel,
		Indexer:     indexer,
		Retriever:   retriever,
		RAGWorkflow: ragWorkflow,
	}

	// 如果启用 Kafka，初始化日志写入器并替换默认日志输出
	if c.Kafka.Enabled {
		writer := logger.NewKafkaLogWriter(c.Kafka.Brokers, c.Kafka.Topic)
		ctxSer.KafkaLogWriter = writer
		logx.SetWriter(writer) // 将所有日志输出到 Kafka（控制台不再输出）
	}

	return ctxSer
}
