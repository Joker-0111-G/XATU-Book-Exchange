package model

import "time"

type Message struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	FromUserID uint      `gorm:"index:idx_from_to;not null" json:"from_user_id"`
	ToUserID   uint      `gorm:"index:idx_from_to;index;not null" json:"to_user_id"`
	BookID     uint      `gorm:"default:0" json:"book_id"`
	Content    string    `gorm:"type:text;not null" json:"content"`
	ContentType int8     `gorm:"default:0" json:"content_type"`
	IsRead     int8      `gorm:"default:0;index" json:"is_read"`
	CreatedAt  time.Time `json:"created_at"`

	// 关联
	FromUser User `gorm:"foreignKey:FromUserID" json:"from_user,omitempty"`
	ToUser   User `gorm:"foreignKey:ToUserID" json:"to_user,omitempty"`
}

func (Message) TableName() string {
	return "messages"
}
