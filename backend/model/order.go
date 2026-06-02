package model

import (
	"time"
)

type Order struct {
	ID           uint       `gorm:"primaryKey" json:"id"`
	OrderNo      string     `gorm:"uniqueIndex;size:32;not null" json:"order_no"`
	SellerID     uint       `gorm:"index;not null" json:"seller_id"`
	BuyerID      uint       `gorm:"index;not null" json:"buyer_id"`
	BookID       uint       `gorm:"index;not null" json:"book_id"`
	Price        float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Status       int8       `gorm:"default:0;index" json:"status"`
	ContactPhone string     `gorm:"size:20;default:''" json:"contact_phone"`
	ContactWechat string   `gorm:"size:50;default:''" json:"contact_wechat"`
	Note         string     `gorm:"size:500;default:''" json:"note"`
	CompletedAt  *time.Time `json:"completed_at"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`

	// 关联
	Seller User `gorm:"foreignKey:SellerID" json:"seller,omitempty"`
	Buyer  User `gorm:"foreignKey:BuyerID" json:"buyer,omitempty"`
	Book   Book `gorm:"foreignKey:BookID" json:"book,omitempty"`
}

func (Order) TableName() string {
	return "orders"
}
