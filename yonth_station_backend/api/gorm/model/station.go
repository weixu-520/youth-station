package model

import "time"

type Station struct {
	Id             int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	StationName    string    `gorm:"column:station_name;type:varchar(128);not null" json:"stationName"`
	District       string    `gorm:"column:district;type:varchar(32);not null;index" json:"district"` // 区域，如：东湖区
	Address        string    `gorm:"column:address;type:varchar(256);not null" json:"address"`
	Latitude       float64   `gorm:"column:latitude;type:decimal(10,7);not null" json:"latitude"`
	Longitude      float64   `gorm:"column:longitude;type:decimal(10,7);not null" json:"longitude"`
	ContactPhone   string    `gorm:"column:contact_phone;type:varchar(20);not null" json:"contactPhone"`
	BusinessHours  string    `gorm:"column:business_hours;type:varchar(64);not null" json:"businessHours"`
	TotalRooms     int       `gorm:"column:total_rooms;type:int;default:0" json:"totalRooms"`
	AvailableRooms int       `gorm:"column:available_rooms;type:int;default:0" json:"availableRooms"`
	Status         int8      `gorm:"column:status;type:tinyint;default:1" json:"status"` // 0-关闭，1-运营中
	Description    string    `gorm:"column:description;type:text" json:"description"`
	Amenities      string    `gorm:"column:amenities;type:varchar(512);default:''" json:"amenities"` // JSON字符串
	NearbyMetro    string    `gorm:"column:nearby_metro;type:varchar(128);default:''" json:"nearbyMetro"`
	ImageUrl       string    `gorm:"column:image_url;type:varchar(256);default:''" json:"imageUrl"`
	WeeklyQuota    int       `gorm:"column:weekly_quota;type:int;default:0" json:"weeklyQuota"`       // 每周配额
	RemainingQuota int       `gorm:"column:remaining_quota;type:int;default:0" json:"remainingQuota"` // 本周剩余
	AvgRating      float32   `gorm:"column:avg_rating;type:decimal(2,1);default:0.0" json:"avgRating"`
	TotalReviews   int       `gorm:"column:total_reviews;type:int;default:0" json:"totalReviews"`
	CreatedAt      time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt      time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (Station) TableName() string {
	return "station"
}
