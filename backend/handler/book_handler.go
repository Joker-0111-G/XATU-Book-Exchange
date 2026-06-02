package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/service"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	svc *service.BookService
}

func NewBookHandler() *BookHandler {
	return &BookHandler{
		svc: service.NewBookService(),
	}
}

// Create 发布图书
func (h *BookHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	book, err := h.svc.Create(userID, &req)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, book)
}

// List 图书列表（公开）
func (h *BookHandler) List(c *gin.Context) {
	var req service.BookListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	books, total, err := h.svc.List(&req)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.SuccessWithPage(c, books, req.Page, req.PageSize, total)
}

// Get 图书详情（公开）
func (h *BookHandler) Get(c *gin.Context) {
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	book, err := h.svc.Get(id)
	if err != nil {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, book)
}

// Search 搜索图书（公开）
func (h *BookHandler) Search(c *gin.Context) {
	var req service.BookListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	books, total, err := h.svc.Search(&req)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.SuccessWithPage(c, books, req.Page, req.PageSize, total)
}

// Update 编辑图书
func (h *BookHandler) Update(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	var req service.UpdateBookReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	if err := h.svc.Update(id, userID, &req); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// UpdateStatus 上架/下架
func (h *BookHandler) UpdateStatus(c *gin.Context) {
	userID := c.GetUint("user_id")
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

	if err := h.svc.UpdateStatus(id, userID, req.Status); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// Delete 删除图书
func (h *BookHandler) Delete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Delete(id, userID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// UserBooks 我发布的图书
func (h *BookHandler) UserBooks(c *gin.Context) {
	userID := c.GetUint("user_id")
	books, err := h.svc.UserBooks(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.Success(c, books)
}
