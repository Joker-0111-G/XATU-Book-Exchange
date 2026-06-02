package common

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

// ParseUintParam 从路径参数中解析 uint，失败则返回错误响应
func ParseUintParam(c *gin.Context, name string) (uint, bool) {
	val := c.Param(name)
	id, err := strconv.ParseUint(val, 10, 64)
	if err != nil {
		Error(c, CodeParamError)
		return 0, false
	}
	return uint(id), true
}

// ParseIntQuery 从查询参数中解析 int，失败返回默认值
func ParseIntQuery(c *gin.Context, name string, defaultVal int) int {
	val := c.DefaultQuery(name, "")
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}
