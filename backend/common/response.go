package common

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type PageMeta struct {
	Page     int   `json:"page"`
	PageSize int   `json:"page_size"`
	Total    int64 `json:"total"`
}

type PageResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    PageMeta    `json:"meta"`
}

// HTTPStatusFromCode 将业务错误码映射为 HTTP 状态码
func HTTPStatusFromCode(code int) int {
	switch code {
	case CodeSuccess:
		return http.StatusOK
	case CodeParamError, CodeUploadTooLarge, CodeUploadType:
		return http.StatusBadRequest
	case CodeUnauthorized, CodeUserCredential:
		return http.StatusUnauthorized
	case CodeForbidden, CodeUserDisabled:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	case CodeAlreadyExists, CodeNotAllowed, CodeBookSoldOut, CodeBookSelfBuy, CodeOrderStatus:
		return http.StatusConflict
	case CodeSystemError:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    CodeSuccess,
		Message: GetMsg(CodeSuccess),
		Data:    data,
	})
}

func SuccessWithPage(c *gin.Context, data interface{}, page, pageSize int, total int64) {
	c.JSON(http.StatusOK, PageResponse{
		Code:    CodeSuccess,
		Message: GetMsg(CodeSuccess),
		Data:    data,
		Meta: PageMeta{
			Page:     page,
			PageSize: pageSize,
			Total:    total,
		},
	})
}

func Error(c *gin.Context, code int) {
	c.JSON(HTTPStatusFromCode(code), Response{
		Code:    code,
		Message: GetMsg(code),
		Data:    nil,
	})
}

func ErrorWithMsg(c *gin.Context, code int, msg string) {
	c.JSON(HTTPStatusFromCode(code), Response{
		Code:    code,
		Message: msg,
		Data:    nil,
	})
}

func SystemError(c *gin.Context) {
	c.JSON(http.StatusInternalServerError, Response{
		Code:    CodeSystemError,
		Message: GetMsg(CodeSystemError),
		Data:    nil,
	})
}
