package model

import "time"

// Payment 支付记录表（押金支付/退款）
type Payment struct {
	Id            int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                          // 支付记录ID
	ApplicationId int64     `gorm:"column:application_id;type:bigint;not null;index" json:"applicationId"` // 关联的申请ID
	TradeNo       string    `gorm:"column:trade_no;type:varchar(64);uniqueIndex;not null" json:"tradeNo"`  // 支付平台流水号
	Amount        int64     `gorm:"column:amount;type:bigint;not null" json:"amount"`                      // 金额（单位：分）
	PayType       int8      `gorm:"column:pay_type;type:tinyint;default:1" json:"payType"`                 // 类型：-押金支付，2-退款
	Status        int8      `gorm:"column:status;type:tinyint;default:0" json:"status"`                    // 状态：0-待支付，1-成功，-失败，-退款中，-已退款
	PayTime       int64     `gorm:"column:pay_time;type:bigint;default:0" json:"payTime"`                  // 支付成功时间戳
	RefundNo      string    `gorm:"column:refund_no;type:varchar(64);default:''" json:"refundNo"`          // 退款流水号
	RefundTime    int64     `gorm:"column:refund_time;type:bigint;default:0" json:"refundTime"`            // 退款完成时间戳
	CreatedAt     time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                     // 记录创建时间
	UpdatedAt     time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`                     // 最后更新时间戳
}

func (Payment) TableName() string {
	return "payment"
}
