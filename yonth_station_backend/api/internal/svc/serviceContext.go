// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package svc

import (
	"yonth_station_backend/api/internal/config"

	"github.com/go-redis/redis"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ServiceContext struct {
	Config config.Config
	RDB    *redis.Client
	DB     *gorm.DB
}

func NewServiceContext(c config.Config) *ServiceContext {
	//连接MySQL
	db, err := gorm.Open(mysql.Open(c.Mysql.DataSource), &gorm.Config{})
	if err != nil {
		logx.Errorf("failed to connect database: %v", err)
		panic(err)
	}
	//连接redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     c.Redis.Addr,     // Redis服务器地址
		Password: c.Redis.Password, // 没有密码，默认值
		DB:       c.Redis.DB,       // 默认DB 0
	})
	return &ServiceContext{
		Config: c,
		RDB:    rdb,
		DB:     db,
	}
}
