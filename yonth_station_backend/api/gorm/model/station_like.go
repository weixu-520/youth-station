package model

import "time"

type StationLike struct {
	Id        int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserId    int64     `gorm:"column:user_id;not null;index:idx_user_station,unique" json:"userId"`
	StationId int64     `gorm:"column:station_id;not null;index:idx_user_station,unique" json:"stationId"`
	CreatedAt time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
}

func (StationLike) TableName() string {
	return "station_like"
}
