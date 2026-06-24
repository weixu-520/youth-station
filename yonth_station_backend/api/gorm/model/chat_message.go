package model

import "time"

type ChatMessage struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	FromUserId int64     `gorm:"column:from_user_id;not null;index" json:"fromUserId"`
	ToUserId   int64     `gorm:"column:to_user_id;default:0" json:"toUserId"`
	TargetType int8      `gorm:"column:target_type;not null" json:"targetType"` // 1-用户→管理员，2-管理员→用户
	Content    string    `gorm:"column:content;type:text;not null" json:"content"`
	IsRead     int8      `gorm:"column:is_read;default:0" json:"isRead"`
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (ChatMessage) TableName() string {
	return "chat_message"
}
