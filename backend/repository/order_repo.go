package repository

import (
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
)

type OrderRepo struct{}

func NewOrderRepo() *OrderRepo {
	return &OrderRepo{}
}

func (r *OrderRepo) Create(order *model.Order) error {
	return database.DB.Create(order).Error
}

func (r *OrderRepo) FindByID(id uint) (*model.Order, error) {
	var order model.Order
	err := database.DB.Preload("Seller").Preload("Buyer").Preload("Book").
		First(&order, id).Error
	return &order, err
}

func (r *OrderRepo) FindByOrderNo(orderNo string) (*model.Order, error) {
	var order model.Order
	err := database.DB.Where("order_no = ?", orderNo).First(&order).Error
	return &order, err
}

func (r *OrderRepo) FindByBuyerID(buyerID uint) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Where("buyer_id = ?", buyerID).
		Preload("Seller").Preload("Book").
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepo) FindBySellerID(sellerID uint) ([]model.Order, error) {
	var orders []model.Order
	err := database.DB.Where("seller_id = ?", sellerID).
		Preload("Buyer").Preload("Book").
		Order("created_at DESC").Find(&orders).Error
	return orders, err
}

func (r *OrderRepo) Update(order *model.Order) error {
	return database.DB.Save(order).Error
}
