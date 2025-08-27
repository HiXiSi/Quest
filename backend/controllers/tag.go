package controllers

import (
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetTags 获取标签列表
func GetTags(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var tags []models.Tag
	var total int64

	query := config.DB.Model(&models.Tag{})

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}

	// 计算总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Order("usage_count DESC, created_at DESC").Offset(offset).Limit(pageSize).Find(&tags).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	utils.PageResponse(c, tags, total, page, pageSize)
}

// CreateTag 创建标签
func CreateTag(c *gin.Context) {
	var tag models.Tag

	if err := c.ShouldBindJSON(&tag); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查标签名称是否已存在
	var existingTag models.Tag
	if err := config.DB.Where("name = ?", tag.Name).First(&existingTag).Error; err == nil {
		utils.ErrorResponse(c, 400, "标签名称已存在")
		return
	}

	// 初始化使用次数
	tag.UsageCount = 0

	if err := config.DB.Create(&tag).Error; err != nil {
		utils.ServerErrorResponse(c, "标签创建失败")
		return
	}

	utils.SuccessResponse(c, tag)
}

// UpdateTag 更新标签
func UpdateTag(c *gin.Context) {
	tagID := c.Param("id")

	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "标签不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	var updateData models.Tag
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查标签名称是否已存在（排除自己）
	if updateData.Name != "" && updateData.Name != tag.Name {
		var existingTag models.Tag
		if err := config.DB.Where("name = ? AND id != ?", updateData.Name, tag.ID).First(&existingTag).Error; err == nil {
			utils.ErrorResponse(c, 400, "标签名称已存在")
			return
		}
		tag.Name = updateData.Name
	}

	// 更新其他字段
	if updateData.Color != "" {
		tag.Color = updateData.Color
	}
	if updateData.Description != "" {
		tag.Description = updateData.Description
	}

	if err := config.DB.Save(&tag).Error; err != nil {
		utils.ServerErrorResponse(c, "标签更新失败")
		return
	}

	utils.SuccessResponse(c, tag)
}

// DeleteTag 删除标签
func DeleteTag(c *gin.Context) {
	tagID := c.Param("id")

	var tag models.Tag
	if err := config.DB.First(&tag, tagID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "标签不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 检查是否有文件使用此标签
	var fileTagCount int64
	config.DB.Model(&models.FileTag{}).Where("tag_id = ?", tag.ID).Count(&fileTagCount)
	if fileTagCount > 0 {
		utils.ErrorResponse(c, 400, "存在文件使用此标签，无法删除")
		return
	}

	if err := config.DB.Delete(&tag).Error; err != nil {
		utils.ServerErrorResponse(c, "标签删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "标签删除成功"})
}
