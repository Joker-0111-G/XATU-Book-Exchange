package repository

import (
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"
)

type BookRepo struct{}

func NewBookRepo() *BookRepo {
	return &BookRepo{}
}

type BookListQuery struct {
	Page       int
	PageSize   int
	CategoryID uint
	MinPrice   float64
	MaxPrice   float64
	Condition  int8
	Status     int8
	Sort       string
	Order      string
	Keyword    string
}

// List 图书列表（分页+筛选+排序）
func (r *BookRepo) List(q *BookListQuery) ([]model.Book, int64, error) {
	var books []model.Book
	var total int64

	db := database.DB.Model(&model.Book{}).Where("deleted_at IS NULL")

	// 筛选条件
	if q.Status > 0 {
		db = db.Where("status = ?", q.Status)
	} else {
		db = db.Where("status = 1") // 默认只看在售
	}
	if q.CategoryID > 0 {
		db = db.Where("category_id = ?", q.CategoryID)
	}
	if q.MinPrice > 0 {
		db = db.Where("selling_price >= ?", q.MinPrice)
	}
	if q.MaxPrice > 0 {
		db = db.Where("selling_price <= ?", q.MaxPrice)
	}
	if q.Condition > 0 {
		db = db.Where("`condition` >= ?", q.Condition)
	}
	if q.Keyword != "" {
		db = db.Where("title LIKE ? OR author LIKE ?", "%"+q.Keyword+"%", "%"+q.Keyword+"%")
	}

	// 统计总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序（仅允许 ASC/DESC，防 SQL 注入）
	order := "DESC"
	if q.Order == "ASC" {
		order = "ASC"
	}
	switch q.Sort {
	case "price":
		db = db.Order("selling_price " + order)
	case "created_at":
		db = db.Order("created_at " + order)
	default:
		db = db.Order("created_at DESC")
	}

	// 分页
	offset := (q.Page - 1) * q.PageSize
	if offset < 0 {
		offset = 0
	}
	db = db.Offset(offset).Limit(q.PageSize)

	// 预加载关联
	db = db.Preload("User").Preload("Category")

	err := db.Find(&books).Error
	return books, total, err
}

// FindByID 查找图书详情
func (r *BookRepo) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	err := database.DB.Preload("User").Preload("Category").First(&book, id).Error
	if err != nil {
		return nil, err
	}
	return &book, nil
}

// Create 创建图书
func (r *BookRepo) Create(book *model.Book) error {
	return database.DB.Create(book).Error
}

// Update 更新图书
func (r *BookRepo) Update(book *model.Book) error {
	return database.DB.Save(book).Error
}

// Delete 软删除图书
func (r *BookRepo) Delete(id uint) error {
	return database.DB.Delete(&model.Book{}, id).Error
}

// FindByUserID 查找用户发布的图书
func (r *BookRepo) FindByUserID(userID uint) ([]model.Book, error) {
	var books []model.Book
	err := database.DB.Where("user_id = ? AND deleted_at IS NULL", userID).
		Preload("Category").
		Order("created_at DESC").
		Find(&books).Error
	return books, err
}
