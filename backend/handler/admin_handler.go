package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"

	"github.com/gin-gonic/gin"
)

type AdminHandler struct{}

func NewAdminHandler() *AdminHandler {
	return &AdminHandler{}
}

// Users 用户管理列表
func (h *AdminHandler) Users(c *gin.Context) {
	var users []model.User
	page := common.ParseIntQuery(c, "page", 1)
	pageSize := common.ParseIntQuery(c, "page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	offset := (page - 1) * pageSize

	var total int64
	if err := database.DB.Model(&model.User{}).Count(&total).Error; err != nil {
		common.SystemError(c)
		return
	}
	if err := database.DB.Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&users).Error; err != nil {
		common.SystemError(c)
		return
	}

	common.SuccessWithPage(c, users, page, pageSize, total)
}

// UpdateUserStatus 启用/禁用用户（检查 RowsAffected）
func (h *AdminHandler) UpdateUserStatus(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	res := database.DB.Model(&model.User{}).Where("id = ?", id).Update("status", req.Status)
	if res.Error != nil {
		common.SystemError(c)
		return
	}
	if res.RowsAffected == 0 {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, nil)
}

// Books 图书管理列表
func (h *AdminHandler) Books(c *gin.Context) {
	var books []model.Book
	page := common.ParseIntQuery(c, "page", 1)
	pageSize := common.ParseIntQuery(c, "page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	status := common.ParseIntQuery(c, "status", 0)
	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Book{}).Where("deleted_at IS NULL")
	if status > 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		common.SystemError(c)
		return
	}
	if err := db.Preload("User").Preload("Category").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&books).Error; err != nil {
		common.SystemError(c)
		return
	}

	common.SuccessWithPage(c, books, page, pageSize, total)
}

// UpdateBookStatus 下架违规图书（检查 RowsAffected）
func (h *AdminHandler) UpdateBookStatus(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	var req struct {
		Status int8 `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	res := database.DB.Model(&model.Book{}).Where("id = ? AND deleted_at IS NULL", id).Update("status", req.Status)
	if res.Error != nil {
		common.SystemError(c)
		return
	}
	if res.RowsAffected == 0 {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, nil)
}

// Orders 全部订单管理
func (h *AdminHandler) Orders(c *gin.Context) {
	var orders []model.Order
	page := common.ParseIntQuery(c, "page", 1)
	pageSize := common.ParseIntQuery(c, "page_size", 20)
	if pageSize > 100 {
		pageSize = 100
	}
	status := common.ParseIntQuery(c, "status", -1)
	offset := (page - 1) * pageSize

	db := database.DB.Model(&model.Order{})
	if status >= 0 {
		db = db.Where("status = ?", status)
	}

	var total int64
	if err := db.Count(&total).Error; err != nil {
		common.SystemError(c)
		return
	}
	if err := db.Preload("Seller").Preload("Buyer").Preload("Book").Offset(offset).Limit(pageSize).Order("created_at DESC").Find(&orders).Error; err != nil {
		common.SystemError(c)
		return
	}

	common.SuccessWithPage(c, orders, page, pageSize, total)
}

// Categories 分类管理列表
func (h *AdminHandler) Categories(c *gin.Context) {
	var categories []model.Category
	if err := database.DB.Order("sort_order ASC, id ASC").Find(&categories).Error; err != nil {
		common.SystemError(c)
		return
	}
	common.Success(c, categories)
}

// CreateCategory 添加分类
func (h *AdminHandler) CreateCategory(c *gin.Context) {
	var req struct {
		Name      string `json:"name" binding:"required"`
		ParentID  uint   `json:"parent_id"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	cat := &model.Category{
		Name:      req.Name,
		ParentID:  req.ParentID,
		Icon:      req.Icon,
		SortOrder: req.SortOrder,
	}
	if err := database.DB.Create(cat).Error; err != nil {
		common.Error(c, common.CodeNotAllowed)
		return
	}

	common.Success(c, cat)
}

// UpdateCategory 编辑分类（检查 RowsAffected）
func (h *AdminHandler) UpdateCategory(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	var req struct {
		Name      string `json:"name"`
		ParentID  uint   `json:"parent_id"`
		Icon      string `json:"icon"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	updates := map[string]interface{}{}
	if req.Name != "" {
		updates["name"] = req.Name
	}
	if req.Icon != "" {
		updates["icon"] = req.Icon
	}
	updates["parent_id"] = req.ParentID
	updates["sort_order"] = req.SortOrder

	res := database.DB.Model(&model.Category{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		common.SystemError(c)
		return
	}
	if res.RowsAffected == 0 {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, nil)
}

// DeleteCategory 删除分类（先检查是否有子分类或关联图书）
func (h *AdminHandler) DeleteCategory(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	// 先检查是否有子分类
	var childCount int64
	database.DB.Model(&model.Category{}).Where("parent_id = ?", id).Count(&childCount)
	if childCount > 0 {
		common.ErrorWithMsg(c, common.CodeNotAllowed, "该分类下有子分类，无法删除")
		return
	}

	// 检查是否有关联图书
	var bookCount int64
	database.DB.Model(&model.Book{}).Where("category_id = ? AND deleted_at IS NULL", id).Count(&bookCount)
	if bookCount > 0 {
		common.ErrorWithMsg(c, common.CodeNotAllowed, "该分类下有关联图书，无法删除")
		return
	}

	if err := database.DB.Delete(&model.Category{}, id).Error; err != nil {
		common.SystemError(c)
		return
	}
	common.Success(c, nil)
}

// CreateBanner 添加轮播图
func (h *AdminHandler) CreateBanner(c *gin.Context) {
	var req struct {
		Title     string `json:"title"`
		ImageURL  string `json:"image_url" binding:"required"`
		LinkURL   string `json:"link_url"`
		SortOrder int    `json:"sort_order"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	banner := &model.Banner{
		Title:     req.Title,
		ImageURL:  req.ImageURL,
		LinkURL:   req.LinkURL,
		SortOrder: req.SortOrder,
		IsActive:  1,
	}
	if err := database.DB.Create(banner).Error; err != nil {
		common.Error(c, common.CodeNotAllowed)
		return
	}

	common.Success(c, banner)
}

// UpdateBanner 编辑轮播图（检查 RowsAffected）
func (h *AdminHandler) UpdateBanner(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	var req struct {
		Title     string `json:"title"`
		ImageURL  string `json:"image_url"`
		LinkURL   string `json:"link_url"`
		SortOrder int    `json:"sort_order"`
		IsActive  int8   `json:"is_active"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	updates := map[string]interface{}{}
	if req.Title != "" {
		updates["title"] = req.Title
	}
	if req.ImageURL != "" {
		updates["image_url"] = req.ImageURL
	}
	if req.LinkURL != "" {
		updates["link_url"] = req.LinkURL
	}
	updates["sort_order"] = req.SortOrder
	updates["is_active"] = req.IsActive

	res := database.DB.Model(&model.Banner{}).Where("id = ?", id).Updates(updates)
	if res.Error != nil {
		common.SystemError(c)
		return
	}
	if res.RowsAffected == 0 {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, nil)
}

// DeleteBanner 删除轮播图
func (h *AdminHandler) DeleteBanner(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}
	if err := database.DB.Delete(&model.Banner{}, id).Error; err != nil {
		common.Error(c, common.CodeNotFound)
		return
	}
	common.Success(c, nil)
}

// Statistics 数据统计
func (h *AdminHandler) Statistics(c *gin.Context) {
	var (
		userCount  int64
		bookCount  int64
		orderCount int64
		doneCount  int64
	)

	database.DB.Model(&model.User{}).Count(&userCount)
	database.DB.Model(&model.Book{}).Where("deleted_at IS NULL").Count(&bookCount)
	database.DB.Model(&model.Order{}).Count(&orderCount)
	database.DB.Model(&model.Order{}).Where("status = 2").Count(&doneCount)

	common.Success(c, gin.H{
		"user_count":       userCount,
		"book_count":       bookCount,
		"order_count":      orderCount,
		"completed_orders": doneCount,
	})
}
