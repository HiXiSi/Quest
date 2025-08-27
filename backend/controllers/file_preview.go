package controllers

import (
	"io"
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// DownloadFile 文件下载
func DownloadFile(c *gin.Context) {
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Where("id = ? AND is_deleted = ?", fileID, false).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "文件不存在")
		return
	}

	// 增加下载次数
	config.DB.Model(&file).UpdateColumn("download_count", gorm.Expr("download_count + ?", 1))

	// 设置响应头
	c.Header("Content-Disposition", "attachment; filename=\""+file.OriginalName+"\"")
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))

	// 发送文件
	c.File(file.FilePath)
}

// PreviewFile 文件预览
func PreviewFile(c *gin.Context) {
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Where("id = ? AND is_deleted = ?", fileID, false).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "文件不存在")
		return
	}

	// 检查文件类型是否支持预览
	if !isPreviewSupported(file.FileType, file.MimeType) {
		utils.ErrorResponse(c, 400, "文件类型不支持预览")
		return
	}

	// 增加查看次数
	config.DB.Model(&file).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	// 对于文本文件，可以直接返回内容
	if file.FileType == "text" {
		content, err := os.ReadFile(file.FilePath)
		if err != nil {
			utils.ServerErrorResponse(c, "读取文件失败")
			return
		}

		// 限制文本文件大小（1MB）
		if len(content) > 1024*1024 {
			utils.ErrorResponse(c, 400, "文本文件过大，无法预览")
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"type":    "text",
			"content": string(content),
		})
		return
	}

	// 对于图片和其他文件，设置适当的Content-Type并返回文件流
	c.Header("Content-Type", file.MimeType)
	c.Header("Content-Length", strconv.FormatInt(file.FileSize, 10))

	// 对于图片，添加缓存控制
	if file.FileType == "image" {
		c.Header("Cache-Control", "public, max-age=3600")
	}

	// 发送文件
	c.File(file.FilePath)
}

// GetFileContent 获取文件内容（用于在线编辑）
func GetFileContent(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, false)

	// 非管理员只能访问自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 只有文本文件支持在线编辑
	if file.FileType != "text" {
		utils.ErrorResponse(c, 400, "只有文本文件支持在线编辑")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "文件不存在")
		return
	}

	// 读取文件内容
	content, err := os.ReadFile(file.FilePath)
	if err != nil {
		utils.ServerErrorResponse(c, "读取文件失败")
		return
	}

	// 限制文本文件大小（1MB）
	if len(content) > 1024*1024 {
		utils.ErrorResponse(c, 400, "文件过大，无法编辑")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"file":    file,
		"content": string(content),
	})
}

// UpdateFileContent 更新文件内容（用于在线编辑）
func UpdateFileContent(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, false)

	// 非管理员只能修改自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 只有文本文件支持在线编辑
	if file.FileType != "text" {
		utils.ErrorResponse(c, 400, "只有文本文件支持在线编辑")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 备份原文件
	backupPath := file.FilePath + ".backup"
	if err := copyFile(file.FilePath, backupPath); err != nil {
		utils.ServerErrorResponse(c, "备份文件失败")
		return
	}

	// 写入新内容
	if err := os.WriteFile(file.FilePath, []byte(req.Content), 0644); err != nil {
		// 恢复备份
		copyFile(backupPath, file.FilePath)
		utils.ServerErrorResponse(c, "保存文件失败")
		return
	}

	// 更新文件信息
	fileInfo, _ := os.Stat(file.FilePath)
	file.FileSize = fileInfo.Size()

	if err := config.DB.Save(&file).Error; err != nil {
		utils.ServerErrorResponse(c, "更新文件信息失败")
		return
	}

	// 删除备份文件
	os.Remove(backupPath)

	utils.SuccessResponse(c, gin.H{"message": "文件保存成功"})
}

// GetFileThumbnail 获取文件缩略图
func GetFileThumbnail(c *gin.Context) {
	fileID := c.Param("id")

	var file models.File
	if err := config.DB.Where("id = ? AND is_deleted = ?", fileID, false).First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 只有图片文件支持缩略图
	if file.FileType != "image" {
		utils.ErrorResponse(c, 400, "只有图片文件支持缩略图")
		return
	}

	// 检查文件是否存在
	if _, err := os.Stat(file.FilePath); os.IsNotExist(err) {
		utils.NotFoundResponse(c, "文件不存在")
		return
	}

	// 简单实现：直接返回原图片（在实际应用中，这里应该生成缩略图）
	c.Header("Content-Type", file.MimeType)
	c.Header("Cache-Control", "public, max-age=3600")
	c.File(file.FilePath)
}

// isPreviewSupported 检查文件类型是否支持预览
func isPreviewSupported(fileType, mimeType string) bool {
	supportedTypes := []string{"image", "text", "pdf", "video"}
	for _, t := range supportedTypes {
		if fileType == t {
			return true
		}
	}

	// 特定MIME类型支持
	supportedMimes := []string{
		"application/json",
		"application/xml",
		"text/html",
		"text/css",
		"text/javascript",
	}
	for _, mime := range supportedMimes {
		if strings.Contains(mimeType, mime) {
			return true
		}
	}

	return false
}

// copyFile 复制文件
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
