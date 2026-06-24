package model

import "time"

type StationComment struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;not null;index" json:"userId"`
	StationId int64     `gorm:"column:station_id;not null;index" json:"stationId"`
	Content   string    `gorm:"column:content;type:varchar(500);not null" json:"content"`
	ParentId  int64     `gorm:"column:parent_id;default:0" json:"parentId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (StationComment) TableName() string {
	return "station_comment"
}
