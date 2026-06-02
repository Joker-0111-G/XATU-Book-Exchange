package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/service"

	"github.com/gin-gonic/gin"
)

type FavoriteHandler struct {
	svc *service.FavoriteService
}

func NewFavoriteHandler() *FavoriteHandler {
	return &FavoriteHandler{
		svc: service.NewFavoriteService(),
	}
}

// List 收藏列表
func (h *FavoriteHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	favs, err := h.svc.List(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, favs)
}

// Add 添加收藏
func (h *FavoriteHandler) Add(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		BookID uint `json:"book_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	if err := h.svc.Add(userID, req.BookID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// Remove 取消收藏
func (h *FavoriteHandler) Remove(c *gin.Context) {
	userID := c.GetUint("user_id")
	bookID, ok := common.ParseUintParam(c, "bookId")
	if !ok {
		return
	}

	if err := h.svc.Remove(userID, bookID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// Check 检查是否已收藏
func (h *FavoriteHandler) Check(c *gin.Context) {
	userID := c.GetUint("user_id")
	bookID, ok := common.ParseUintParam(c, "bookId")
	if !ok {
		return
	}

	exists, err := h.svc.Check(userID, bookID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.Success(c, gin.H{"favorited": exists})
}
