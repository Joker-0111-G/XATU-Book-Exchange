package service

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"time"

	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
	"xatu-book-exchange/repository"

	"gorm.io/gorm"
)

type OrderService struct {
	orderRepo *repository.OrderRepo
	bookRepo  *repository.BookRepo
}

func NewOrderService() *OrderService {
	return &OrderService{
		orderRepo: repository.NewOrderRepo(),
		bookRepo:  repository.NewBookRepo(),
	}
}

type CreateOrderReq struct {
	BookID       uint   `json:"book_id" binding:"required"`
	ContactPhone string `json:"contact_phone"`
	ContactWechat string `json:"contact_wechat"`
	Note         string `json:"note"`
}
// generateOrderNo 生成唯一订单号: XATU + 时间戳(13位毫秒) + 4位随机数
func generateOrderNo() string {
	n, _ := rand.Int(rand.Reader, big.NewInt(10000))
	return fmt.Sprintf("XATU%d%04d", time.Now().UnixMilli(), n.Int64())
}

// Create 创建订单（使用 GORM 事务，带并发安全）
func (s *OrderService) Create(buyerID uint, req *CreateOrderReq) (*model.Order, error) {
	book, err := s.bookRepo.FindByID(req.BookID)
	if err != nil {
		return nil, ErrNotFound
	}
	if book.UserID == buyerID {
		return nil, fmt.Errorf("不能购买自己的图书")
	}

	orderNo := generateOrderNo()

	order := &model.Order{
		OrderNo:       orderNo,
		SellerID:      book.UserID,
		BuyerID:       buyerID,
		BookID:        req.BookID,
		Price:         book.SellingPrice,
		Status:        0,
		ContactPhone:  req.ContactPhone,
		ContactWechat: req.ContactWechat,
		Note:          req.Note,
	}

	// 事务内创建订单 + 标记图书已售
	// WHERE status=1 防止并发重复下单，RowsAffected 检查确保只有一个人成功
	err = database.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(order).Error; err != nil {
			return err
		}
		res := tx.Model(&model.Book{}).Where("id = ? AND status = 1", req.BookID).
			Update("status", 2)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			// 图书状态已被其他事务修改（已售），回滚订单
			return fmt.Errorf("图书已售出")
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return order, nil
}

// List 买家订单列表
func (s *OrderService) List(buyerID uint) ([]model.Order, error) {
	return s.orderRepo.FindByBuyerID(buyerID)
}

// SalesList 卖家订单列表
func (s *OrderService) SalesList(sellerID uint) ([]model.Order, error) {
	return s.orderRepo.FindBySellerID(sellerID)
}

// Get 订单详情
func (s *OrderService) Get(id, userID uint) (*model.Order, error) {
	order, err := s.orderRepo.FindByID(id)
	if err != nil {
		return nil, ErrNotFound
	}
	if order.BuyerID != userID && order.SellerID != userID {
		return nil, ErrNoPermission
	}
	return order, nil
}

// Confirm 卖家确认订单（事务保护状态检查+更新）
func (s *OrderService) Confirm(orderID, sellerID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			return ErrNotFound
		}
		if order.SellerID != sellerID {
			return ErrNoPermission
		}
		res := tx.Model(&model.Order{}).Where("id = ? AND status = 0", orderID).
			Update("status", 1)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrOrderStatus
		}
		return nil
	})
}

// Complete 由买家确认完成交易（事务保护）
func (s *OrderService) Complete(orderID, userID uint) error {
	now := time.Now()
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			return ErrNotFound
		}
		// 只允许买家确认完成
		if order.BuyerID != userID {
			return ErrNoPermission
		}
		res := tx.Model(&model.Order{}).Where("id = ? AND status = 1", orderID).
			Updates(map[string]interface{}{"status": 2, "completed_at": &now})
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrOrderStatus
		}
		return nil
	})
}

// Cancel 取消订单（含事务，带 RowsAffected 检查）
func (s *OrderService) Cancel(orderID, userID uint) error {
	return database.DB.Transaction(func(tx *gorm.DB) error {
		var order model.Order
		if err := tx.First(&order, orderID).Error; err != nil {
			return ErrNotFound
		}
		if order.BuyerID != userID {
			return ErrNoPermission
		}
		// 订单状态 0->3，同时图书恢复在售
		res := tx.Model(&model.Order{}).Where("id = ? AND status = 0", orderID).
			Update("status", 3)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected == 0 {
			return ErrOrderStatus
		}
		tx.Model(&model.Book{}).Where("id = ?", order.BookID).
			Update("status", 1)
		return nil
	})
}
