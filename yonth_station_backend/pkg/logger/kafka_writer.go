package logger

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/zeromicro/go-zero/core/logx"
)

type KafkaLogWriter struct {
	writer *kafka.Writer
}

type LogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
	Service   string `json:"service"`
	TraceID   string `json:"trace_id,omitempty"`
}

func NewKafkaLogWriter(brokers []string, topic string) *KafkaLogWriter {
	w := &kafka.Writer{
		Addr:                   kafka.TCP(brokers...),
		Topic:                  topic,
		Balancer:               &kafka.LeastBytes{},
		RequiredAcks:           kafka.RequireOne,
		AllowAutoTopicCreation: true,
		Async:                  true,
		BatchSize:              100,
		BatchTimeout:           5 * time.Second,
	}
	return &KafkaLogWriter{writer: w}
}

// 实现 io.Writer 接口
func (w *KafkaLogWriter) Write(p []byte) (n int, err error) {
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     "INFO",
		Message:   string(p),
		Service:   "gateway-api",
	}
	data, err := json.Marshal(entry)
	if err != nil {
		return 0, err
	}
	msg := kafka.Message{
		Key:   []byte("log"),
		Value: data,
	}
	ctx := context.Background()
	err = w.writer.WriteMessages(ctx, msg)
	if err != nil {
		return 0, err
	}
	return len(p), nil
}

func (w *KafkaLogWriter) writeLevel(level string, v interface{}) {
	// 用 fmt.Sprint 兼容所有类型（string / fmt.Stringer / struct 等）
	msg := fmt.Sprint(v)
	entry := LogEntry{
		Timestamp: time.Now().Format(time.RFC3339),
		Level:     level,
		Message:   msg,
		Service:   "gateway-api",
	}
	data, _ := json.Marshal(entry)
	if err := w.writer.WriteMessages(context.Background(), kafka.Message{Key: []byte("log"), Value: data}); err != nil {
		fmt.Fprintf(os.Stderr, "[Kafka] write error: %v\n", err)
	}
}

func (w *KafkaLogWriter) Alert(v interface{})                     { w.writeLevel("ALERT", v) }
func (w *KafkaLogWriter) Close() error                            { return w.writer.Close() }
func (w *KafkaLogWriter) Debug(v interface{}, _ ...logx.LogField) { w.writeLevel("DEBUG", v) }
func (w *KafkaLogWriter) Error(v interface{}, _ ...logx.LogField) { w.writeLevel("ERROR", v) }
func (w *KafkaLogWriter) Info(v interface{}, _ ...logx.LogField)  { w.writeLevel("INFO", v) }
func (w *KafkaLogWriter) Severe(v interface{})                    { w.writeLevel("SEVERE", v) }
func (w *KafkaLogWriter) Slow(v interface{}, _ ...logx.LogField)  { w.writeLevel("SLOW", v) }
func (w *KafkaLogWriter) Stack(v interface{})                     { w.writeLevel("STACK", v) }
func (w *KafkaLogWriter) Stat(v interface{}, _ ...logx.LogField)  { w.writeLevel("STAT", v) }
