// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/api/internal/config"
	"yonth_station_backend/api/internal/websocket"
	"yonth_station_backend/pkg/logger"
	"yonth_station_backend/pkg/rabbitmq"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config         config.Config
	Redis          *redis.Redis
	DB             *gorm.DB
	KafkaLogWriter *logger.KafkaLogWriter
	RabbitMQ       *rabbitmq.RabbitMQClient
	Hub            *websocket.Hub
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接MySQL
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Errorf("failed to connect database: %v", err)
		panic(err)
	}

	//自动迁移表结构（仅添加缺失列，不删除已有列/数据）
	if err := db.AutoMigrate(&model.User{}, &model.Station{}, &model.Room{}, &model.Application{}, &model.Payment{}, &model.ChatMessage{}, &model.StationComment{}, &model.StationLike{}); err != nil {
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
	ctx := &ServiceContext{
		Config:   c,
		DB:       db,
		Redis:    rds,
		RabbitMQ: mqClient,
		Hub:      hub,
	}

	// 如果启用 Kafka，初始化日志写入器并替换默认日志输出
	if c.Kafka.Enabled {
		writer := logger.NewKafkaLogWriter(c.Kafka.Brokers, c.Kafka.Topic)
		ctx.KafkaLogWriter = writer
		logx.SetWriter(writer) // 将所有日志输出到 Kafka（控制台不再输出）
	}

	return ctx
}
