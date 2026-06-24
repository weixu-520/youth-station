package websocket

import (
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userId int64
}

// readPump 从 WebSocket 读取消息，并交给 Hub 处理
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		var m Message
		if err := json.Unmarshal(msg, &m); err != nil {
			continue // 解析错误则忽略
		}
		m.FromUserId = c.userId // 自动填充发送者 ID
		c.hub.HandleMessage(m)  // 交给 Hub 处理
	}
}

// writePump 从 send channel 取消息写入 WebSocket
func (c *Client) writePump() {
	ticker := time.NewTicker(30 * time.Second) // 心跳 Ping
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case msg, ok := <-c.send:
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			c.conn.WriteMessage(websocket.TextMessage, msg)
		case <-ticker.C:
			// 发送 Ping 保持连接活跃
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
