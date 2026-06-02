package model

import "time"

type Banner struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Title     string    `gorm:"size:100;default:''" json:"title"`
	ImageURL  string    `gorm:"size:255;not null" json:"image_url"`
	LinkURL   string    `gorm:"size:255;default:''" json:"link_url"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	IsActive  int8      `gorm:"default:1" json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

func (Banner) TableName() string {
	return "banners"
}
