package handler

import (
	"xatu-book-exchange/common"
	"xatu-book-exchange/service"
	"xatu-book-exchange/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	svc *service.UserService
}

func NewUserHandler() *UserHandler {
	return &UserHandler{
		svc: service.NewUserService(),
	}
}

// Register 用户注册
func (h *UserHandler) Register(c *gin.Context) {
	var req service.RegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	user, err := h.svc.Register(&req)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeParamError, err.Error())
		return
	}

	common.Success(c, user)
}

// Login 用户登录
func (h *UserHandler) Login(c *gin.Context) {
	var req service.LoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	resp, err := h.svc.Login(&req)
	if err != nil {
		if err == service.ErrUserDisabled {
			common.ErrorWithMsg(c, common.CodeUserDisabled, err.Error())
		} else {
			common.ErrorWithMsg(c, common.CodeUserCredential, err.Error())
		}
		return
	}

	common.Success(c, resp)
}

// RefreshToken 刷新 Token
func (h *UserHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	resp, err := h.svc.RefreshToken(req.RefreshToken)
	if err != nil {
		common.ErrorWithMsg(c, common.CodeUnauthorized, err.Error())
		return
	}

	common.Success(c, resp)
}

// Logout 用户登出（将 access_token 加入 Redis 黑名单）
func (h *UserHandler) Logout(c *gin.Context) {
	claims, exists := c.Get("token_claims")
	if !exists {
		common.Success(c, nil)
		return
	}

	// 将当前 access_token 加入黑名单
	if err := utils.AddToBlacklist(claims.(*utils.Claims)); err != nil {
		common.SystemError(c)
		return
	}

	common.Success(c, nil)
}

// Profile 获取个人信息
func (h *UserHandler) Profile(c *gin.Context) {
	userID := c.GetUint("user_id")
	user, err := h.svc.Profile(userID)
	if err != nil {
		common.Error(c, common.CodeNotFound)
		return
	}

	common.Success(c, user)
}

// UpdateProfile 更新个人信息
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req service.UpdateProfileReq
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	if err := h.svc.UpdateProfile(userID, &req); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}

// ChangePassword 修改密码
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userID := c.GetUint("user_id")
	var req struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	if err := h.svc.ChangePassword(userID, req.OldPassword, req.NewPassword); err != nil {
		common.ErrorWithMsg(c, common.CodeNotAllowed, err.Error())
		return
	}

	common.Success(c, nil)
}
