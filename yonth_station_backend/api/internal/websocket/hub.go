package websocket

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"yonth_station_backend/api/gorm/model"
	"yonth_station_backend/pkg/rabbitmq"

	"gorm.io/gorm"
)

// Hub 管理所有 WebSocket 连接，处理消息路由和广播
type Hub struct {
	clients    map[int64]*Client        // 在线用户映射：userId -> Client
	register   chan *Client             // 新连接注册通道
	unregister chan *Client             // 连接关闭注销通道
	db         *gorm.DB                 // 数据库（持久化消息）
	mqClient   *rabbitmq.RabbitMQClient // RabbitMQ 客户端
	mu         sync.RWMutex             // 保护 clients 并发访问
	seenMsgs   map[string]time.Time     // 消息去重：hash → 过期时间
	seenMu     sync.Mutex               // 保护 seenMsgs
}

// NewHub 创建 Hub 实例
func NewHub(db *gorm.DB, mqClient *rabbitmq.RabbitMQClient) *Hub {
	h := &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		db:         db,
		mqClient:   mqClient,
		seenMsgs:   make(map[string]time.Time),
	}
	// 定期清理过期的去重记录
	go func() {
		for {
			time.Sleep(30 * time.Second)
			h.seenMu.Lock()
			now := time.Now()
			for k, t := range h.seenMsgs {
				if now.Sub(t) > 60*time.Second {
					delete(h.seenMsgs, k)
				}
			}
			h.seenMu.Unlock()
		}
	}()
	return h
}

// Run 启动 Hub 主循环，监听注册/注销事件和 RabbitMQ 广播消息
func (h *Hub) Run() {
	// 启动 RabbitMQ 消费者（持续接收广播消息）
	go h.consumeRabbitMQ()

	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.userId] = client
			h.mu.Unlock()
			log.Printf("[WebSocket] User %d connected", client.userId)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.userId]; ok {
				delete(h.clients, client.userId)
				close(client.send)
				log.Printf("[WebSocket] User %d disconnected", client.userId)
			}
			h.mu.Unlock()
		}
	}
}

// msgHash 计算消息的去重哈希（fromUserId + content + createdAt）//保证消息的幂等性
func msgHash(m *Message) string {
	s := fmt.Sprintf("%d|%s|%d", m.FromUserId, m.Content, m.CreatedAt)
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// isDuplicate 检查消息是否已在短时间内处理过
func (h *Hub) isDuplicate(hash string) bool {
	h.seenMu.Lock()
	defer h.seenMu.Unlock()
	if _, ok := h.seenMsgs[hash]; ok {
		return true
	}
	h.seenMsgs[hash] = time.Now()
	return false
}

// consumeRabbitMQ 持续消费 RabbitMQ 广播消息
// 当其他实例发布消息时，会通过此函数接收并转发给本地用户
func (h *Hub) consumeRabbitMQ() {
	if h.mqClient == nil {
		return // RabbitMQ 未配置，跳过
	}
	deliveries, err := h.mqClient.Consume()
	if err != nil {
		log.Printf("[RabbitMQ] Failed to start consuming: %v", err)
		return
	}
	for d := range deliveries {
		var msg Message
		if err := json.Unmarshal(d.Body, &msg); err != nil {
			continue
		}
		// 去重：HandleMessage 已经本地投递过了，跳过
		hash := msgHash(&msg)
		if h.isDuplicate(hash) {
			continue
		}
		// 转发给本地在线用户（来自其他实例的消息）
		if msg.ToUserId != 0 {
			h.sendToUser(msg.ToUserId, d.Body)
		} else {
			h.sendToAdmins(d.Body)
		}
	}
}

// sendToUser 发送消息给指定用户（若在线）
func (h *Hub) sendToUser(userId int64, data []byte) {
	h.mu.RLock()
	client, ok := h.clients[userId]
	h.mu.RUnlock()
	if ok {
		select {
		case client.send <- data:
		default:
			// 如果发送 channel 已满，丢弃消息（可记录日志）
			log.Printf("[WebSocket] User %d send channel full", userId)
		}
	}
}

// sendToAdmins 广播消息给所有在线管理员
func (h *Hub) sendToAdmins(data []byte) {
	// 先获取所有在线管理员 ID（从数据库查询角色为 admin 的用户）
	// 为减少数据库压力，可缓存管理员列表
	adminIds := h.getAdminIds()
	h.mu.RLock()
	defer h.mu.RUnlock()
	for _, uid := range adminIds {
		if client, ok := h.clients[uid]; ok {
			select {
			case client.send <- data:
			default:
			}
		}
	}
}

// getAdminIds 获取所有管理员 ID（可缓存）
func (h *Hub) getAdminIds() []int64 {
	var users []model.User
	// 从数据库查询 is_admin = 1 的用户
	h.db.Select("id").Where("is_admin = ?", true).Find(&users)
	ids := make([]int64, len(users))
	for i, u := range users {
		ids[i] = u.Id
	}
	return ids
}

// broadcastToOtherInstances 将消息发布到 RabbitMQ，通知其他实例
func (h *Hub) broadcastToOtherInstances(m Message) {
	if err := h.mqClient.PublishMessage(m); err != nil {
		log.Printf("[RabbitMQ] Publish error: %v", err)
	}
}

func (h *Hub) HandleMessage(m Message) {
	// 1. 持久化到 MySQL
	msg := &model.ChatMessage{
		FromUserId: m.FromUserId,
		ToUserId:   m.ToUserId,
		TargetType: m.TargetType,
		Content:    m.Content,
	}
	if err := h.db.Create(msg).Error; err != nil {
		log.Printf("[Database] Save message error: %v", err)
		return
	}

	// 2. 登记去重（先登记，这样 RabbitMQ 回显时自动跳过）
	hash := msgHash(&m)
	h.isDuplicate(hash)

	data, _ := json.Marshal(m)

	// 3. 本地直接投递（核心通道）
	if m.ToUserId != 0 {
		h.sendToUser(m.ToUserId, data)
	} else {
		h.sendToAdmins(data)
	}

	// 4. 广播到 RabbitMQ（多实例场景，单实例可跳过）
	if h.mqClient != nil {
		if err := h.mqClient.PublishMessage(m); err != nil {
			log.Printf("[RabbitMQ] Publish error: %v", err)
		}
	}
}
