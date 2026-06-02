package utils

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var PhoneRegex = regexp.MustCompile(`^1[3-9]\d{9}$`)

// ValidatePhone 自定义手机号校验
func ValidatePhone(fl validator.FieldLevel) bool {
	phone := fl.Field().String()
	return PhoneRegex.MatchString(phone)
}

// ValidatePassword 密码强度校验（6-32位）
func ValidatePassword(fl validator.FieldLevel) bool {
	pwd := fl.Field().String()
	if len(pwd) < 6 || len(pwd) > 32 {
		return false
	}
	return true
}
