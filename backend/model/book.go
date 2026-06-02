package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// StringArray 用于 JSON 类型字段
type StringArray []string

func (s *StringArray) Scan(val interface{}) error {
	if val == nil {
		*s = StringArray{}
		return nil
	}
	bytes, ok := val.([]byte)
	if !ok {
		return errors.New("类型断言失败")
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	if s == nil {
		return json.Marshal([]string{})
	}
	return json.Marshal(s)
}

type Book struct {
	ID           uint           `gorm:"primaryKey" json:"id"`
	UserID       uint           `gorm:"index;not null" json:"user_id"`
	CategoryID   uint           `gorm:"index;not null" json:"category_id"`
	Title        string         `gorm:"size:200;not null;index" json:"title"`
	Author       string         `gorm:"size:100;default:''" json:"author"`
	SellingPrice float64        `gorm:"type:decimal(10,2);not null" json:"selling_price"`
	Condition    int8           `gorm:"not null" json:"condition"`
	Description  string         `gorm:"type:text" json:"description"`
	Images       StringArray    `gorm:"type:json" json:"images"`
	Status       int8           `gorm:"default:1;index" json:"status"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"deleted_at"`

	// 关联
	User     User     `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
}

func (Book) TableName() string {
	return "books"
}
