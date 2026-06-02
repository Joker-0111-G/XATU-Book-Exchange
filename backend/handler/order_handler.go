package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/service"

	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	svc *service.OrderService
}

func NewOrderHandler() *OrderHandler {
	return &OrderHandler{
		svc: service.NewOrderService(),
	}
}

// Create 创建订单
func (h *OrderHandler) Create(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.CreateOrderReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	order, err := h.svc.Create(userID, &req)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, order)
}

// List 买家订单列表
func (h *OrderHandler) List(c *gin.Context) {
	userID := c.GetUint("user_id")
	orders, err := h.svc.List(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, orders)
}

// SalesList 卖家订单列表
func (h *OrderHandler) SalesList(c *gin.Context) {
	userID := c.GetUint("user_id")
	orders, err := h.svc.SalesList(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, orders)
}

// UserOrders 我买的订单（买家视角）
func (h *OrderHandler) UserOrders(c *gin.Context) {
	h.List(c)
}

// UserSales 我的售出（卖家视角）
func (h *OrderHandler) UserSales(c *gin.Context) {
	h.SalesList(c)
}

// Get 订单详情
func (h *OrderHandler) Get(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	order, err := h.svc.Get(id, userID)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, order)
}

// Confirm 卖家确认订单
func (h *OrderHandler) Confirm(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Confirm(id, userID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// Complete 确认完成
func (h *OrderHandler) Complete(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Complete(id, userID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// Cancel 取消订单
func (h *OrderHandler) Cancel(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, ok := common.ParseUintParam(c, "id")
	if !ok {
		return
	}

	if err := h.svc.Cancel(id, userID); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}
