package controllers

import (
	"io"
	"material-platform/utils"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
)

// AssetUploadResponse 资源上传响应结构
type AssetUploadResponse struct {
	URL      string `json:"url"`      // 相对路径URL
	FileName string `json:"filename"` // 文件名
	FileSize int64  `json:"filesize"` // 文件大小
}

// UploadAsset 通用资源上传接口
// 用于表单字段的文件上传，与文件管理模块完全独立
func UploadAsset(c *gin.Context) {
	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, "文件上传失败: "+err.Error())
		return
	}
	defer file.Close()

	// 验证文件大小（限制50MB，比文件管理模块的限制小一些）
	if header.Size > 50*1024*1024 {
		utils.ErrorResponse(c, 400, "文件大小不能超过50MB")
		return
	}

	// 生成文件名和路径
	fileName := utils.GenerateFileName(header.Filename)

	// 创建资源上传目录（与文件管理模块的uploads目录分开）
	uploadDir := filepath.Join("../assets", time.Now().Format("2006/01/02"))
	if err := utils.EnsureDir(uploadDir); err != nil {
		utils.ServerErrorResponse(c, "创建上传目录失败")
		return
	}

	// 保存文件
	filePath := filepath.Join(uploadDir, fileName)
	dst, err := os.Create(filePath)
	if err != nil {
		utils.ServerErrorResponse(c, "文件保存失败")
		return
	}
	defer dst.Close()

	// 复制文件内容
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.ServerErrorResponse(c, "文件保存失败")
		return
	}

	// 生成相对路径URL（用于前端访问）
	relativeURL := "/assets/" + time.Now().Format("2006/01/02") + "/" + fileName

	// 返回文件信息
	utils.SuccessResponse(c, AssetUploadResponse{
		URL:      relativeURL,
		FileName: fileName,
		FileSize: header.Size,
	})
}
