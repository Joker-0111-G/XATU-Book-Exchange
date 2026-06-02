package model

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	Phone        string         `gorm:"uniqueIndex;size:20;not null" json:"phone"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Nickname     string         `gorm:"size:50;not null" json:"nickname"`
	Avatar       string         `gorm:"size:255;default:''" json:"avatar"`
	Email        string         `gorm:"size:100;default:''" json:"email"`
	Major        string         `gorm:"size:100;default:''" json:"major"`
	Wechat       string         `gorm:"size:50;default:''" json:"wechat"`
	Status       int8           `gorm:"default:1;index" json:"status"`
	IsAdmin      int8           `gorm:"default:0" json:"is_admin"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

func (User) TableName() string {
	return "users"
}
