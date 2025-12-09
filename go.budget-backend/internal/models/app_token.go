package models

import "time"

type AppTokenModel struct {
	BaseModel
	TargetID  uint      `json:"-" gorm:"not null" `
	Type      string    `json:"-" gorm:"type:varchar(100);not null"`
	Used      bool      `json:"-" gorm:"default:false;not null"`
	Token     string    `json:"token" gorm:"type:varchar(255);not null;uniqueIndex"`
	ExpiredAt time.Time `json:"-" gorm:"not null"`
}

func (AppTokenModel) TableName() string {
	return "app_tokens"
}
