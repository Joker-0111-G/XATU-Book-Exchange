package service

import "errors"

var (
	ErrNoPermission      = errors.New("无权操作")
	ErrBookSoldOut       = errors.New("图书已售出")
	ErrAlreadyExist      = errors.New("资源已存在")
	ErrNotFound          = errors.New("资源不存在")
	ErrOrderStatus       = errors.New("订单状态不允许当前操作")
	ErrInvalidCredential = errors.New("手机号或密码错误")
	ErrUserDisabled      = errors.New("账号已被禁用")
)
