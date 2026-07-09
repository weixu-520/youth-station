// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package config

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf
	Auth struct {
		AccessSecret string
		AccessExpire int64
	}

	Redis redis.RedisConf

	RateLimit struct {
		Mode  string `json:",default=global"` // global 或 ip
		Rate  int    `json:",default=10"`
		Burst int    `json:",default=20"`
	}

	Kafka struct {
		Brokers []string
		Topic   string
		Enabled bool
	}

	Mysql struct {
		DataSource string
	}

	RabbitMQ struct {
		Host     string
		Port     int
		User     string
		Password string
		Exchange string
	}
	//智能客服配置M配置
	Chat struct {
		//LLM配置
		DeepSeek struct {
			APIKey  string
			Model   string
			BaseURL string
		}
		//向量数据库配置
		Milvus struct {
			Address    string
			Collection string
		}
		//向量模型配置
		Embedding struct {
			APIKey        string
			Model         string
			UseMultiModal bool `json:",default=true"` // 多模态模型需开启
		}
		//增强检索配置
		RAG struct {
			TopK           int
			ScoreThreshold float64
		}
	}
}
