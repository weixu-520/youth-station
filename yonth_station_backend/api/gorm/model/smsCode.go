package model

import "time"

// SmsCode 短信验证码记录表（用于记录历史，实际校验建议使用Redis）
type SmsCode struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                  // 记录ID
	Phone     string    `gorm:"column:phone;type:varchar(11);not null;index" json:"phone"`     // 手机号?
	Code      string    `gorm:"column:code;type:varchar(6);not null" json:"code"`              // 6位数字验证码
	Used      int8      `gorm:"column:used;type:tinyint;default:0" json:"used"`                // 是否已使用：0-未使用，1-已使�?
	ExpiredAt int64     `gorm:"column:expired_at;type:bigint;not null;index" json:"expiredAt"` // 过期时间戳（Unix秒）
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`             // 创建时间
}

func (SmsCode) TableName() string {
	return "sms_code"
}
