package common

// 业务错误码
const (
	CodeSuccess = 0

	// 通用错误 1000-1999
	CodeParamError    = 1001
	CodeUnauthorized  = 1002
	CodeForbidden     = 1003
	CodeNotFound      = 1004
	CodeAlreadyExists = 1005
	CodeNotAllowed    = 1006

	// 用户错误 2000-2999
	CodeUserCredential = 2001
	CodeUserDisabled   = 2002

	// 图书错误 3000-3999
	CodeBookSoldOut   = 3001
	CodeBookSelfBuy   = 3002

	// 订单错误 4000-4999
	CodeOrderStatus   = 4001

	// 上传错误 5000-5999
	CodeUploadTooLarge = 5001
	CodeUploadType     = 5002

	// 系统错误
	CodeSystemError = 9999
)

var codeMsgMap = map[int]string{
	CodeSuccess:        "success",
	CodeParamError:     "参数错误",
	CodeUnauthorized:   "未认证",
	CodeForbidden:      "权限不足",
	CodeNotFound:       "资源不存在",
	CodeAlreadyExists:  "资源已存在",
	CodeNotAllowed:     "操作不允许",
	CodeUserCredential: "手机号或密码错误",
	CodeUserDisabled:   "用户已被禁用",
	CodeBookSoldOut:    "图书已售出",
	CodeBookSelfBuy:    "不能购买自己的图书",
	CodeOrderStatus:    "订单状态不允许当前操作",
	CodeUploadTooLarge: "上传文件过大",
	CodeUploadType:     "上传文件类型不支持",
	CodeSystemError:    "系统内部错误",
}

func GetMsg(code int) string {
	msg, ok := codeMsgMap[code]
	if !ok {
		return "未知错误"
	}
	return msg
}
