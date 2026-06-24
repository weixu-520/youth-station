package websocket

// Message 定义了 WebSocket 消息的传输格式
// 客户端发送消息时必须符合此结构，服务端广播时也使用此结构
type Message struct {
	// 消息发送者 ID（服务端自动填充，客户端可忽略）
	FromUserId int64 `json:"fromUserId,omitempty"`

	// 消息接收者 ID，0 表示发送给管理员（即客服）
	ToUserId int64 `json:"toUserId"`

	// 消息类型：1-用户发给管理员，2-管理员发给用户
	// 此字段用于标识消息方向，便于服务端处理
	TargetType int8 `json:"targetType"`

	// 消息内容（文本）
	Content string `json:"content"`

	// 消息创建时间（服务端生成，客户端可忽略）
	CreatedAt int64 `json:"createdAt,omitempty"`
}
