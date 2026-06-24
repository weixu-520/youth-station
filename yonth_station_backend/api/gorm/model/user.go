package model

import "time"

type User struct {
	Id           int64     `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UserName     string    `gorm:"column:user_name;type:varchar(64);uniqueIndex;not null" json:"userName"`
	Phone        *string   `gorm:"column:phone;type:varchar(11);uniqueIndex" json:"phone,omitempty"` // 修改为指针，允许NULL
	Password     string    `gorm:"column:password;type:varchar(128);not null" json:"-"`
	IdCard       string    `gorm:"column:id_card;type:varchar(18);default:''" json:"idCard"`
	BirthDate    string    `gorm:"column:birth_date;type:date;default:null" json:"birthDate"`
	Gender       int8      `gorm:"column:gender;type:tinyint;default:0" json:"gender"`
	Education    int8      `gorm:"column:education;type:tinyint;default:0" json:"education"`
	School       string    `gorm:"column:school;type:varchar(128);default:''" json:"school"`
	GraduateYear int       `gorm:"column:graduate_year;type:int;default:0" json:"graduateYear"`
	HukouCity    string    `gorm:"column:hukou_city;type:varchar(64);default:''" json:"hukouCity"` // 户籍城市，格式：省+市
	Status       int8      `gorm:"column:status;type:tinyint;default:0" json:"status"`             // 账号状态：0-正常，1-冻结
	IsAdmin      bool      `gorm:"column:is_admin;type:tinyint(1);default:0" json:"isAdmin"`
	LastLoginAt  int64     `gorm:"column:last_login_at;type:bigint;default:0" json:"lastLoginAt"`
	CreatedAt    time.Time `gorm:"column:created_at;autoCreateTime" json:"createdAt"`
	UpdatedAt    time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updatedAt"`
}

func (User) TableName() string {
	return "user"
}
