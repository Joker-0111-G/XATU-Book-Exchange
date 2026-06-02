package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/service"

	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	svc *service.MessageService
}

func NewMessageHandler() *MessageHandler {
	return &MessageHandler{
		svc: service.NewMessageService(),
	}
}

// Conversations 会话列表
func (h *MessageHandler) Conversations(c *gin.Context) {
	userID := c.GetUint("user_id")
	msgs, err := h.svc.Conversations(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, msgs)
}

// History 聊天记录
func (h *MessageHandler) History(c *gin.Context) {
	userID := c.GetUint("user_id")
	otherUserID, ok := common.ParseUintParam(c, "userId")
	if !ok {
		return
	}

	msgs, err := h.svc.History(userID, otherUserID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}
	common.Success(c, msgs)
}

// Send 发送消息
func (h *MessageHandler) Send(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.SendMessageReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	msg, err := h.svc.Send(userID, &req)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, msg)
}

// MarkRead 标记已读
func (h *MessageHandler) MarkRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		FromUserID uint `json:"from_user_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	if err := h.svc.MarkRead(req.FromUserID, userID); err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.Success(c, nil)
}

// UnreadCount 未读消息数
func (h *MessageHandler) UnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.svc.UnreadCount(userID)
	if err != nil {
		common.Error(c, common.CodeSystemError)
		return
	}

	common.Success(c, gin.H{"unread_count": count})
}
