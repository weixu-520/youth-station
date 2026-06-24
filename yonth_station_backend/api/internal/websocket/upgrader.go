package websocket

import (
	"net/http"
	"yonth_station_backend/pkg/utils"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // 允许跨域，生产环境应限制
}

// ServeWs 升级 HTTP 请求为 WebSocket，并注册 Client
func ServeWs(hub *Hub, secret string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 从 URL 参数获取 token
		token := r.URL.Query().Get("token")
		if token == "" {
			http.Error(w, "missing token", http.StatusUnauthorized)
			return
		}
		// 解析 JWT 获取 userId
		userId, err := utils.ParseToken(token, secret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			http.Error(w, "upgrade failed", http.StatusInternalServerError)
			return
		}

		client := &Client{
			hub:    hub,
			conn:   conn,
			send:   make(chan []byte, 256),
			userId: userId,
		}
		hub.register <- client

		// 启动读写协程
		go client.writePump()
		go client.readPump()
	}
}
