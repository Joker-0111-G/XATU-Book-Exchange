package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"xatu-book-exchange/common"
	"xatu-book-exchange/config"

	"github.com/gin-gonic/gin"
)

type UploadHandler struct{}

func NewUploadHandler() *UploadHandler {
	return &UploadHandler{}
}

// UploadImage 上传图片
func (h *UploadHandler) UploadImage(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	// 检查大小
	maxSize := config.AppConfig.Upload.MaxSize
	if file.Size > int64(maxSize)*1024*1024 {
		common.Error(c, common.CodeUploadTooLarge)
		return
	}

	// 检查文件扩展名
	ext := filepath.Ext(file.Filename)
	allowed := false
	for _, t := range config.AppConfig.Upload.AllowTypes {
		if ext == t {
			allowed = true
			break
		}
	}
	if !allowed {
		common.Error(c, common.CodeUploadType)
		return
	}

	// 检查 MIME 类型（读取文件头）
	f, err := file.Open()
	if err != nil {
		common.SystemError(c)
		return
	}
	defer f.Close()

	buf := make([]byte, 512)
	if _, err := f.Read(buf); err != nil {
		common.SystemError(c)
		return
	}
	contentType := http.DetectContentType(buf)
	if contentType != "image/jpeg" && contentType != "image/png" && contentType != "image/webp" {
		common.ErrorWithMsg(c, common.CodeUploadType, "仅支持 JPEG/PNG/WebP 格式")
		return
	}

	// 生成安全的存储路径
	now := time.Now()
	dir := fmt.Sprintf("%s/%d/%02d/%02d", config.AppConfig.Upload.Dir, now.Year(), now.Month(), now.Day())
	if err := os.MkdirAll(dir, 0755); err != nil {
		common.SystemError(c)
		return
	}

	filename := fmt.Sprintf("%d%06d%s", time.Now().UnixNano(), os.Getpid()%1000000, ext)
	dst := filepath.Join(dir, filename)

	// 确保路径在 upload 目录内（防目录遍历）
	absUploadDir, _ := filepath.Abs(config.AppConfig.Upload.Dir)
	absDst, _ := filepath.Abs(filepath.Dir(dst))
	if !strings.HasPrefix(absDst, absUploadDir) {
		common.SystemError(c)
		return
	}

	if err := c.SaveUploadedFile(file, dst); err != nil {
		common.SystemError(c)
		return
	}

	common.Success(c, gin.H{
		"url": "/uploads/images/" + now.Format("2006/01/02") + "/" + filename,
	})
}

// DeleteImage 删除图片（只能删除 upload 目录内的文件）
func (h *UploadHandler) DeleteImage(c *gin.Context) {
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		common.Error(c, common.CodeParamError)
		return
	}

	// 从 URL 提取文件路径
	path := req.URL
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// 安全检查：路径必须在 upload 目录内
	absPath, _ := filepath.Abs(path)
	absUploadDir, _ := filepath.Abs(config.AppConfig.Upload.Dir)
	if !strings.HasPrefix(absPath, absUploadDir) {
		common.ErrorWithMsg(c, common.CodeForbidden, "不允许删除此路径的文件")
		return
	}

	if err := os.Remove(absPath); err != nil {
		// 文件不存在不报错
		common.Success(c, nil)
		return
	}

	common.Success(c, nil)
}
