package model

import "time"

type Category struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"size:50;not null" json:"name"`
	ParentID  uint      `gorm:"default:0;index" json:"parent_id"`
	Icon      string    `gorm:"size:255;default:''" json:"icon"`
	SortOrder int       `gorm:"default:0" json:"sort_order"`
	CreatedAt time.Time `json:"created_at"`
}

func (Category) TableName() string {
	return "categories"
}
