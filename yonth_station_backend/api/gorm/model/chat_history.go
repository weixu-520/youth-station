package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

// ChatHistory 对话历史表，记录用户与智能客服的每一次对话
type ChatHistory struct {
	// Id 对话记录唯一标识，自增主键
	Id int64 `gorm:"column:id;primaryKey;autoIncrement" json:"id"`

	// SessionId 会话ID，用于关联同一用户的多轮对话
	// 每次新对话会生成一个新的 SessionId，或由前端传入已有 SessionId
	SessionId string `gorm:"column:session_id;type:varchar(64);not null;index" json:"sessionId"`

	// UserId 用户ID，关联用户表，用于追溯对话者身份
	UserId int64 `gorm:"column:user_id;type:bigint;not null;index" json:"userId"`

	// Question 用户提出的问题（原文）
	Question string `gorm:"column:question;type:text;not null" json:"question"`

	// Answer 智能客服返回的回答（原文）
	Answer string `gorm:"column:answer;type:text;not null" json:"answer"`

	// Sources 回答所引用的知识库来源（JSON 数组）
	// 存储格式：["doc_id_1", "doc_id_2", ...]
	// 便于后续分析哪些文档被高频引用
	Sources JSONArray `gorm:"column:sources;type:json;default:null" json:"sources,omitempty"`

	// CreatedAt 对话发生时间（自动填充）
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

// TableName 返回数据库表名
func (ChatHistory) TableName() string {
	return "chat_history"
}

// JSONArray 自定义类型，用于处理 GORM 对 JSON 数组的序列化/反序列化
type JSONArray []string

// Scan 实现 sql.Scanner 接口，从数据库读取 JSON 数据
func (j *JSONArray) Scan(value interface{}) error {
	if value == nil {
		*j = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}
	return json.Unmarshal(bytes, j)
}

// Value 实现 driver.Valuer 接口，将 Go 数据写入数据库
func (j JSONArray) Value() (driver.Value, error) {
	if len(j) == 0 {
		return "[]", nil
	}
	return json.Marshal(j)
}
