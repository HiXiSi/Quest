package controllers

import (
	"fmt"
	"io"
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// UploadFile 文件上传
func UploadFile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	// 获取上传的文件
	file, header, err := c.Request.FormFile("file")
	if err != nil {
		utils.ErrorResponse(c, 400, "文件上传失败: "+err.Error())
		return
	}
	defer file.Close()

	// 获取可选参数
	categoryID := c.PostForm("category_id")
	description := c.PostForm("description")
	tagIDs := c.PostForm("tag_ids") // 逗号分隔的标签ID

	// 验证文件大小（限制100MB）
	if header.Size > 100*1024*1024 {
		utils.ErrorResponse(c, 400, "文件大小不能超过100MB")
		return
	}

	// 计算文件哈希
	md5Hash, sha256Hash, err := utils.GetFileHash(file)
	if err != nil {
		utils.ServerErrorResponse(c, "文件哈希计算失败")
		return
	}

	// 检查文件是否已存在（通过MD5判断）
	var existingFile models.File
	if err := config.DB.Where("md5_hash = ? AND user_id = ?", md5Hash, userID).First(&existingFile).Error; err == nil {
		utils.ErrorResponse(c, 400, "文件已存在")
		return
	}

	// 生成文件名和路径
	fileName := utils.GenerateFileName(header.Filename)
	mimeType := utils.GetMimeType(header.Filename)
	fileType := utils.GetFileType(mimeType)

	// 创建上传目录
	uploadDir := filepath.Join("../uploads", time.Now().Format("2006/01/02"))
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
	file.Seek(0, 0) // 重置文件指针
	_, err = io.Copy(dst, file)
	if err != nil {
		utils.ServerErrorResponse(c, "文件保存失败")
		return
	}

	// 验证分类ID
	var categoryIDPtr *uint
	if categoryID != "" {
		catID, err := strconv.ParseUint(categoryID, 10, 32)
		if err == nil {
			var category models.Category
			if err := config.DB.First(&category, catID).Error; err == nil {
				categoryIDUint := uint(catID)
				categoryIDPtr = &categoryIDUint
			}
		}
	}

	// 创建文件记录
	fileRecord := models.File{
		OriginalName: header.Filename,
		FileName:     fileName,
		FilePath:     filePath,
		FileSize:     header.Size,
		FileType:     fileType,
		MimeType:     mimeType,
		MD5Hash:      md5Hash,
		SHA256Hash:   sha256Hash,
		Description:  description,
		UserID:       userID.(uint),
		CategoryID:   categoryIDPtr,
	}

	if err := config.DB.Create(&fileRecord).Error; err != nil {
		// 删除已保存的文件
		os.Remove(filePath)
		utils.ServerErrorResponse(c, "文件记录保存失败")
		return
	}

	// 处理标签关联
	if tagIDs != "" {
		tagIDList := strings.Split(tagIDs, ",")
		for _, tagIDStr := range tagIDList {
			tagID, err := strconv.ParseUint(strings.TrimSpace(tagIDStr), 10, 32)
			if err != nil {
				continue
			}

			// 验证标签存在
			var tag models.Tag
			if err := config.DB.First(&tag, tagID).Error; err != nil {
				continue
			}

			// 创建文件标签关联
			fileTag := models.FileTag{
				FileID: fileRecord.ID,
				TagID:  uint(tagID),
			}
			config.DB.Create(&fileTag)

			// 更新标签使用次数
			config.DB.Model(&tag).UpdateColumn("usage_count", gorm.Expr("usage_count + ?", 1))
		}
	}

	// 预加载关联数据
	config.DB.Preload("User").Preload("Category").Preload("Tags").First(&fileRecord, fileRecord.ID)

	utils.SuccessResponse(c, models.FileUploadResponse{
		FileID:   fileRecord.ID,
		FileName: fileRecord.FileName,
		FileSize: fileRecord.FileSize,
		FilePath: fileRecord.FilePath,
	})
}

// GetFiles 获取文件列表
func GetFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")
	categoryID := c.Query("category_id")
	tagID := c.Query("tag_id")
	fileType := c.Query("file_type")
	sortBy := c.DefaultQuery("sort_by", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var files []models.File
	var total int64

	query := config.DB.Model(&models.File{}).Where("is_deleted = ?", false)

	// 非管理员只能看到自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("original_name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 分类筛选
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}

	// 文件类型筛选
	if fileType != "" {
		query = query.Where("file_type = ?", fileType)
	}

	// 标签筛选
	if tagID != "" {
		query = query.Joins("JOIN file_tags ON files.id = file_tags.file_id").
			Where("file_tags.tag_id = ?", tagID)
	}

	// 计算总数
	query.Count(&total)

	// 排序
	orderClause := fmt.Sprintf("%s %s", sortBy, strings.ToUpper(sortOrder))
	query = query.Order(orderClause)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Category").Preload("Tags").
		Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	utils.PageResponse(c, files, total, page, pageSize)
}

// GetFile 获取单个文件信息
func GetFile(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, false)

	// 非管理员只能看到自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Preload("User").Preload("Category").Preload("Tags").First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 增加查看次数
	config.DB.Model(&file).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1))

	utils.SuccessResponse(c, file)
}

// UpdateFile 更新文件信息
func UpdateFile(c *gin.Context) {
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

	var updateData struct {
		OriginalName string `json:"original_name"`
		Description  string `json:"description"`
		CategoryID   *uint  `json:"category_id"`
		TagIDs       []uint `json:"tag_ids"`
		IsPublic     *bool  `json:"is_public"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 更新文件信息
	if updateData.OriginalName != "" {
		file.OriginalName = updateData.OriginalName
	}
	if updateData.Description != "" {
		file.Description = updateData.Description
	}
	if updateData.CategoryID != nil {
		// 验证分类存在
		var category models.Category
		if err := config.DB.First(&category, *updateData.CategoryID).Error; err == nil {
			file.CategoryID = updateData.CategoryID
		}
	}
	if updateData.IsPublic != nil {
		file.IsPublic = *updateData.IsPublic
	}

	// 保存文件更新
	if err := config.DB.Save(&file).Error; err != nil {
		utils.ServerErrorResponse(c, "文件信息更新失败")
		return
	}

	// 更新标签关联
	if updateData.TagIDs != nil {
		// 删除现有标签关联
		config.DB.Where("file_id = ?", file.ID).Delete(&models.FileTag{})

		// 创建新的标签关联
		for _, tagID := range updateData.TagIDs {
			var tag models.Tag
			if err := config.DB.First(&tag, tagID).Error; err != nil {
				continue
			}

			fileTag := models.FileTag{
				FileID: file.ID,
				TagID:  tagID,
			}
			config.DB.Create(&fileTag)
		}
	}

	// 预加载关联数据
	config.DB.Preload("User").Preload("Category").Preload("Tags").First(&file, file.ID)

	utils.SuccessResponse(c, file)
}

// DeleteFile 删除文件（移到回收站）
func DeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, false)

	// 非管理员只能删除自己的文件
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

	// 标记为已删除
	now := time.Now()
	file.IsDeleted = true
	file.DeletedAt = &now

	if err := config.DB.Save(&file).Error; err != nil {
		utils.ServerErrorResponse(c, "文件删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "文件已移到回收站"})
}

// RestoreFile 恢复文件
func RestoreFile(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, true)

	// 非管理员只能恢复自己的文件
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

	// 恢复文件
	file.IsDeleted = false
	file.DeletedAt = nil

	if err := config.DB.Save(&file).Error; err != nil {
		utils.ServerErrorResponse(c, "文件恢复失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "文件恢复成功"})
}

// BatchDeleteFiles 批量删除文件
func BatchDeleteFiles(c *gin.Context) {
	var req struct {
		FileIDs []uint `json:"file_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	query := config.DB.Model(&models.File{}).Where("id IN ? AND is_deleted = ?", req.FileIDs, false)

	// 非管理员只能删除自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	now := time.Now()
	updates := map[string]interface{}{
		"is_deleted": true,
		"deleted_at": now,
	}

	if err := query.Updates(updates).Error; err != nil {
		utils.ServerErrorResponse(c, "批量删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "批量删除成功"})
}

// BatchRestoreFiles 批量恢复文件
func BatchRestoreFiles(c *gin.Context) {
	var req struct {
		FileIDs []uint `json:"file_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	query := config.DB.Model(&models.File{}).Where("id IN ? AND is_deleted = ?", req.FileIDs, true)

	// 非管理员只能恢复自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	updates := map[string]interface{}{
		"is_deleted": false,
		"deleted_at": nil,
	}

	if err := query.Updates(updates).Error; err != nil {
		utils.ServerErrorResponse(c, "批量恢复失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "批量恢复成功"})
}

// GetDeletedFiles 获取回收站文件列表
func GetDeletedFiles(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var files []models.File
	var total int64

	query := config.DB.Model(&models.File{}).Where("is_deleted = ?", true)

	// 非管理员只能看到自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").Preload("Category").Preload("Tags").
		Order("deleted_at DESC").Offset(offset).Limit(pageSize).Find(&files).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	utils.PageResponse(c, files, total, page, pageSize)
}

// EmptyRecycleBin 清空回收站
func EmptyRecycleBin(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var files []models.File
	query := config.DB.Where("is_deleted = ?", true)

	// 非管理员只能清空自己的回收站
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&files).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 删除物理文件和数据库记录
	for _, file := range files {
		// 删除物理文件
		if err := os.Remove(file.FilePath); err != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("删除物理文件失败: %s, 错误: %v\n", file.FilePath, err)
		}

		// 删除文件标签关联
		config.DB.Where("file_id = ?", file.ID).Delete(&models.FileTag{})

		// 删除数据库记录
		config.DB.Unscoped().Delete(&file)
	}

	utils.SuccessResponse(c, gin.H{"message": "回收站清空成功"})
}

// PermanentDeleteFile 彻底删除单个文件
func PermanentDeleteFile(c *gin.Context) {
	fileID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var file models.File
	query := config.DB.Where("id = ? AND is_deleted = ?", fileID, true)

	// 非管理员只能删除自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&file).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "文件不存在或未在回收站中")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 删除物理文件
	if err := os.Remove(file.FilePath); err != nil {
		// 记录错误但继续删除数据库记录
		fmt.Printf("删除物理文件失败: %s, 错误: %v\n", file.FilePath, err)
	}

	// 删除文件标签关联
	config.DB.Where("file_id = ?", file.ID).Delete(&models.FileTag{})

	// 彻底删除数据库记录
	if err := config.DB.Unscoped().Delete(&file).Error; err != nil {
		utils.ServerErrorResponse(c, "文件删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "文件已彻底删除"})
}

// BatchPermanentDeleteFiles 批量彻底删除文件
func BatchPermanentDeleteFiles(c *gin.Context) {
	var req struct {
		FileIDs []uint `json:"file_ids" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var files []models.File
	query := config.DB.Where("id IN ? AND is_deleted = ?", req.FileIDs, true)

	// 非管理员只能删除自己的文件
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Find(&files).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	if len(files) == 0 {
		utils.ErrorResponse(c, 400, "没有找到可删除的文件")
		return
	}

	// 删除物理文件和数据库记录
	var deletedCount int
	for _, file := range files {
		// 删除物理文件
		if err := os.Remove(file.FilePath); err != nil {
			// 记录错误但继续处理其他文件
			fmt.Printf("删除物理文件失败: %s, 错误: %v\n", file.FilePath, err)
		}

		// 删除文件标签关联
		config.DB.Where("file_id = ?", file.ID).Delete(&models.FileTag{})

		// 彻底删除数据库记录
		if err := config.DB.Unscoped().Delete(&file).Error; err == nil {
			deletedCount++
		}
	}

	utils.SuccessResponse(c, gin.H{
		"message":       fmt.Sprintf("成功彻底删除 %d 个文件", deletedCount),
		"deleted_count": deletedCount,
	})
}
