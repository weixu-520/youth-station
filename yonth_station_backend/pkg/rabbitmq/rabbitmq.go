package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQClient 封装 RabbitMQ 连接和操作
type RabbitMQClient struct {
	conn         *amqp.Connection
	channel      *amqp.Channel
	exchangeName string
	mu           sync.Mutex // 保护 Publish 并发
}

// NewRabbitMQClient 创建并初始化 RabbitMQ 客户端
// host: RabbitMQ 服务地址, port: 端口, user/password: 认证, exchange: 交换机名称（fanout 类型）
func NewRabbitMQClient(host string, port int, user, password, exchange string) (*RabbitMQClient, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%d/", user, password, host, port)
	conn, err := amqp.Dial(url)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, fmt.Errorf("failed to open channel: %w", err)
	}

	// 声明 fanout 类型交换机（广播模式）
	err = channel.ExchangeDeclare(
		exchange, // 交换机名称
		"fanout", // 类型：广播
		true,     // durable：持久化（重启保留）
		false,    // auto-deleted：不自动删除
		false,    // internal：内部使用
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		conn.Close()
		channel.Close()
		return nil, fmt.Errorf("failed to declare exchange: %w", err)
	}

	return &RabbitMQClient{
		conn:         conn,
		channel:      channel,
		exchangeName: exchange,
	}, nil
}

// Publish 发布消息到交换机（广播给所有绑定队列）
func (c *RabbitMQClient) Publish(body []byte) error {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.channel.PublishWithContext(
		context.Background(),
		c.exchangeName, // exchange
		"",             // routing key（fanout 忽略）
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

// Consume 订阅消息，返回一个接收通道（deliveries）
// 每次调用会创建一个临时独占队列，绑定到交换机，连接断开时队列自动删除
func (c *RabbitMQClient) Consume() (<-chan amqp.Delivery, error) {
	// 声明临时队列（名称留空则随机生成，exclusive=true 表示连接断开时删除）
	q, err := c.channel.QueueDeclare(
		"",    // name（空字符串自动生成）
		false, // durable
		false, // delete when unused
		true,  // exclusive（连接断开时自动删除）
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return nil, fmt.Errorf("failed to declare queue: %w", err)
	}

	// 将队列绑定到交换机
	err = c.channel.QueueBind(
		q.Name,         // queue name
		"",             // routing key（fanout 忽略）
		c.exchangeName, // exchange
		false,
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to bind queue: %w", err)
	}

	// 开始消费（auto-ack: true 自动确认，适用于广播场景）
	return c.channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack（自动确认）
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
}

// Close 关闭连接和通道
func (c *RabbitMQClient) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// PublishMessage 泛型发布，自动序列化
func (c *RabbitMQClient) PublishMessage(msg interface{}) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}
	return c.Publish(data)
}
