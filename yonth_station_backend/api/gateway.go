// Code scaffolded by goctl. Safe to edit.
// goctl 1.9.2

package main

import (
	"flag"
	"fmt"
	"net/http"

	"yonth_station_backend/api/internal/config"
	"yonth_station_backend/api/internal/handler"
	"yonth_station_backend/api/internal/middleware"
	"yonth_station_backend/api/internal/svc"
	"yonth_station_backend/api/internal/websocket"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/gateway-api.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	ctx := svc.NewServiceContext(c)
	defer server.Stop()
	// 注册限流中间件（根据配置自动选择模式）
	server.Use(middleware.RateLimit(
		ctx.Redis,
		c.RateLimit.Mode,
		c.RateLimit.Rate,
		c.RateLimit.Burst,
	))
	defer func() {
		if ctx.KafkaLogWriter != nil {
			_ = ctx.KafkaLogWriter.Close()
		}
	}()

	// 添加 WebSocket 路由
	server.AddRoute(rest.Route{
		Method:  http.MethodGet,
		Path:    "/api/v1/ws",
		Handler: websocket.ServeWs(ctx.Hub, ctx.Config.Auth.AccessSecret),
	})
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
