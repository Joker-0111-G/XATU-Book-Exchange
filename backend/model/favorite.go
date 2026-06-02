package model

import "time"

type Favorite struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	UserID    uint      `gorm:"uniqueIndex:idx_user_book;not null" json:"user_id"`
	BookID    uint      `gorm:"uniqueIndex:idx_user_book;not null" json:"book_id"`
	CreatedAt time.Time `json:"created_at"`

	// 关联
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Book Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

func (Favorite) TableName() string {
	return "favorites"
}
