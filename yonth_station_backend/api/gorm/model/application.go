package model

import "time"

// Application 申请记录表
type Application struct {
	Id                    int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                                              // 申请ID
	UserId                int64     `gorm:"column:user_id;type:bigint;not null;index" json:"userId"`                                   // 申请人用户ID
	StationId             int64     `gorm:"column:station_id;type:bigint;not null;index" json:"stationId"`                             // 申请的驿站ID
	RoomId                int64     `gorm:"column:room_id;type:bigint;default:0" json:"roomId"`                                        // 实际分配的房间ID（审核通过后分配）
	CheckinDate           string    `gorm:"column:checkin_date;type:date;not null" json:"checkinDate"`                                 // 计划入住日期（YYYY-MM-DD）
	CheckoutDate          string    `gorm:"column:checkout_date;type:date;not null" json:"checkoutDate"`                               // 计划退房日期（YYYY-MM-DD）
	Status                int8      `gorm:"column:status;type:tinyint;default:0;index" json:"status"`                                  // 申请状态：0-待审核，1-通过）2-拒绝）3-已取消，4-已入住，5-已退房
	VisitPurpose          int8      `gorm:"column:visit_purpose;type:tinyint;default:0" json:"visitPurpose"`                           // 来访目的）-求职）-创业）-研学
	InterviewProofType    int8      `gorm:"column:interview_proof_type;type:tinyint;default:0" json:"interviewProofType"`              // 面试证明类型）-邮件）-截图）-公告）-函件
	InterviewProofContent string    `gorm:"column:interview_proof_content;type:varchar(512);default:''" json:"interviewProofContent"`  // 证明内容文字描述
	InterviewProofFileUrl string    `gorm:"column:interview_proof_file_url;type:varchar(256);default:''" json:"interviewProofFileUrl"` // 证明文件URL（图片或PDF）
	BusinessPlan          string    `gorm:"column:business_plan;type:text" json:"businessPlan"`                                        // 创业计划简介（当目的为创业时）
	Remark                string    `gorm:"column:remark;type:varchar(256);default:''" json:"remark"`                                  // 用户备注
	RejectReason          string    `gorm:"column:reject_reason;type:varchar(256);default:''" json:"rejectReason"`                     // 审核拒绝原因
	DepositAmount         int64     `gorm:"column:deposit_amount;type:bigint;default:0" json:"depositAmount"`                          // 押金金额（单位：分）
	DepositStatus         int8      `gorm:"column:deposit_status;type:tinyint;default:0" json:"depositStatus"`                         // 押金状态：0-未缴纳，1-已缴纳，2-已退款
	CheckinAt             int64     `gorm:"column:checkin_at;type:bigint;default:0" json:"checkinAt"`                                  // 实际入住时间戳（Unix秒）
	CheckoutAt            int64     `gorm:"column:checkout_at;type:bigint;default:0" json:"checkoutAt"`                                // 实际退房时间戳
	AuditBy               string    `gorm:"column:audit_by;type:varchar(64);default:''" json:"auditBy"`                                // 审核人（管理员账号）
	AuditAt               int64     `gorm:"column:audit_at;type:bigint;default:0" json:"auditAt"`                                      // 审核时间?
	AppliedAt             int64     `gorm:"column:applied_at;type:bigint;not null" json:"appliedAt"`                                   // 申请提交时间?
	UpdatedAt             int64     `gorm:"column:updated_at;type:bigint;not null" json:"updatedAt"`                                   // 记录更新时间戳（Unix秒）
	CreatedAt             time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`                                         // 创建时间（datetime）
}

func (Application) TableName() string {
	return "application"
}
