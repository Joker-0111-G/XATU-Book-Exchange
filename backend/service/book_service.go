package service

import (
	"xatu-book-exchange/model"
	"xatu-book-exchange/repository"
)

type BookService struct {
	bookRepo *repository.BookRepo
}

func NewBookService() *BookService {
	return &BookService{
		bookRepo: repository.NewBookRepo(),
	}
}

type CreateBookReq struct {
	Title        string   `json:"title" binding:"required"`
	Author       string   `json:"author"`
	CategoryID   uint     `json:"category_id" binding:"required"`
	SellingPrice float64  `json:"selling_price" binding:"required"`
	Condition    int8     `json:"condition" binding:"required"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
}

type UpdateBookReq struct {
	Title        string   `json:"title"`
	Author       string   `json:"author"`
	CategoryID   uint     `json:"category_id"`
	SellingPrice float64  `json:"selling_price"`
	Condition    int8     `json:"condition"`
	Description  string   `json:"description"`
	Images       []string `json:"images"`
}

type BookListReq struct {
	Page       int     `form:"page"`
	PageSize   int     `form:"page_size"`
	CategoryID uint    `form:"category_id"`
	MinPrice   float64 `form:"min_price"`
	MaxPrice   float64 `form:"max_price"`
	Condition  int8    `form:"condition"`
	Sort       string  `form:"sort"`
	Order      string  `form:"order"`
	Keyword    string  `form:"keyword"`
}

// Create 发布图书
func (s *BookService) Create(userID uint, req *CreateBookReq) (*model.Book, error) {
	book := &model.Book{
		UserID:       userID,
		CategoryID:   req.CategoryID,
		Title:        req.Title,
		Author:       req.Author,
		SellingPrice: req.SellingPrice,
		Condition:    req.Condition,
		Description:  req.Description,
		Images:       req.Images,
		Status:       1,
	}
	if err := s.bookRepo.Create(book); err != nil {
		return nil, err
	}
	return book, nil
}

// List 图书列表
func (s *BookService) List(req *BookListReq) ([]model.Book, int64, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}
	if req.Order == "" {
		req.Order = "DESC"
	}

	q := &repository.BookListQuery{
		Page:       req.Page,
		PageSize:   req.PageSize,
		CategoryID: req.CategoryID,
		MinPrice:   req.MinPrice,
		MaxPrice:   req.MaxPrice,
		Condition:  req.Condition,
		Sort:       req.Sort,
		Order:      req.Order,
		Keyword:    req.Keyword,
		Status:     1, // 默认只看在售
	}

	return s.bookRepo.List(q)
}

// Get 图书详情
func (s *BookService) Get(id uint) (*model.Book, error) {
	return s.bookRepo.FindByID(id)
}

// Update 更新图书
func (s *BookService) Update(bookID, userID uint, req *UpdateBookReq) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book.UserID != userID {
		return ErrNoPermission
	}

	if req.Title != "" {
		book.Title = req.Title
	}
	if req.Author != "" {
		book.Author = req.Author
	}
	if req.CategoryID > 0 {
		book.CategoryID = req.CategoryID
	}
	if req.SellingPrice > 0 {
		book.SellingPrice = req.SellingPrice
	}
	if req.Condition > 0 {
		book.Condition = req.Condition
	}
	if req.Description != "" {
		book.Description = req.Description
	}
	if req.Images != nil {
		book.Images = req.Images
	}

	return s.bookRepo.Update(book)
}

// UpdateStatus 更新图书状态
func (s *BookService) UpdateStatus(bookID, userID uint, status int8) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book.UserID != userID {
		return ErrNoPermission
	}
	if book.Status == 2 {
		return ErrBookSoldOut
	}
	book.Status = status
	return s.bookRepo.Update(book)
}

// Delete 删除图书
func (s *BookService) Delete(bookID, userID uint) error {
	book, err := s.bookRepo.FindByID(bookID)
	if err != nil {
		return err
	}
	if book.UserID != userID {
		return ErrNoPermission
	}
	return s.bookRepo.Delete(bookID)
}

// UserBooks 获取用户发布的图书
func (s *BookService) UserBooks(userID uint) ([]model.Book, error) {
	return s.bookRepo.FindByUserID(userID)
}

// Search 搜索图书
func (s *BookService) Search(req *BookListReq) ([]model.Book, int64, error) {
	return s.List(req)
}
