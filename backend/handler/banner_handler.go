package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/database"
	"xatu-book-exchange/model"

	"github.com/gin-gonic/gin"
)

type BannerHandler struct{}

func NewBannerHandler() *BannerHandler {
	return &BannerHandler{}
}

// List 获取启用的轮播图
func (h *BannerHandler) List(c *gin.Context) {
	var banners []model.Banner
	if err := database.DB.Where("is_active = 1").Order("sort_order ASC, id ASC").Find(&banners).Error; err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, banners)
}

// Statistics 公开统计信息（无需登录）
func (h *BannerHandler) Statistics(c *gin.Context) {
	var (
		userCount    int64
		bookCount    int64
		orderCount   int64
		doneCount    int64
	)

	database.DB.Model(&model.User{}).Count(&userCount)
	database.DB.Model(&model.Book{}).Where("deleted_at IS NULL AND status = 1").Count(&bookCount)
	database.DB.Model(&model.Order{}).Count(&orderCount)
	database.DB.Model(&model.Order{}).Where("status = 2").Count(&doneCount)

	common.Success(c, gin.H{
		"user_count":       userCount,
		"book_count":       bookCount,
		"order_count":      orderCount,
		"completed_orders": doneCount,
	})
}
