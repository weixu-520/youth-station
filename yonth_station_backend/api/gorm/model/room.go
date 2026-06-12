package model

import "time"

// Room 房间表（每个驿站下的具体房间）
type Room struct {
	Id         int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`                   // 房间ID
	StationId  int64     `gorm:"column:station_id;type:bigint;not null;index" json:"stationId"`  // 所属驿站ID
	RoomNumber string    `gorm:"column:room_number;type:varchar(20);not null" json:"roomNumber"` // 房间号，如"301"
	Status     int8      `gorm:"column:status;type:tinyint;default:0" json:"status"`             // 状态：0-空闲，1-占用
	CreatedAt  time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`              // 创建时间
	UpdatedAt  time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`              // 更新时间
}

func (Room) TableName() string {
	return "room"
}
