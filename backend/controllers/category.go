package controllers

import (
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetCategories 获取分类列表
func GetCategories(c *gin.Context) {
	var categories []models.Category

	// 支持树形结构和平铺结构
	treeMode := c.Query("tree") == "true"

	if treeMode {
		// 获取根分类及其子分类
		if err := config.DB.Where("parent_id IS NULL").Preload("Children").Find(&categories).Error; err != nil {
			utils.ServerErrorResponse(c, "数据库查询失败")
			return
		}
	} else {
		// 获取所有分类（平铺）
		if err := config.DB.Preload("Parent").Find(&categories).Error; err != nil {
			utils.ServerErrorResponse(c, "数据库查询失败")
			return
		}
	}

	utils.SuccessResponse(c, categories)
}

// CreateCategory 创建分类
func CreateCategory(c *gin.Context) {
	var category models.Category

	if err := c.ShouldBindJSON(&category); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查分类名称是否已存在（同一父级下）
	var existingCategory models.Category
	query := config.DB.Where("name = ?", category.Name)

	if category.ParentID != nil {
		query = query.Where("parent_id = ?", *category.ParentID)
	} else {
		query = query.Where("parent_id IS NULL")
	}

	if err := query.First(&existingCategory).Error; err == nil {
		utils.ErrorResponse(c, 400, "同一层级下分类名称已存在")
		return
	}

	// 如果指定了父分类，验证父分类是否存在
	if category.ParentID != nil {
		var parentCategory models.Category
		if err := config.DB.First(&parentCategory, *category.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.ErrorResponse(c, 400, "父分类不存在")
				return
			}
			utils.ServerErrorResponse(c, "数据库查询失败")
			return
		}
	}

	if err := config.DB.Create(&category).Error; err != nil {
		utils.ServerErrorResponse(c, "分类创建失败")
		return
	}

	// 预加载父分类信息
	config.DB.Preload("Parent").First(&category, category.ID)

	utils.SuccessResponse(c, category)
}

// UpdateCategory 更新分类
func UpdateCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "分类不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	var updateData models.Category
	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查分类名称是否已存在（同一父级下，排除自己）
	if updateData.Name != "" && updateData.Name != category.Name {
		var existingCategory models.Category
		query := config.DB.Where("name = ? AND id != ?", updateData.Name, category.ID)

		if updateData.ParentID != nil {
			query = query.Where("parent_id = ?", *updateData.ParentID)
		} else {
			query = query.Where("parent_id IS NULL")
		}

		if err := query.First(&existingCategory).Error; err == nil {
			utils.ErrorResponse(c, 400, "同一层级下分类名称已存在")
			return
		}
		category.Name = updateData.Name
	}

	// 如果要更改父分类，需要验证
	if updateData.ParentID != nil {
		// 防止设置自己为父分类
		if *updateData.ParentID == category.ID {
			utils.ErrorResponse(c, 400, "不能将自己设置为父分类")
			return
		}

		// 验证父分类是否存在
		var parentCategory models.Category
		if err := config.DB.First(&parentCategory, *updateData.ParentID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				utils.ErrorResponse(c, 400, "父分类不存在")
				return
			}
			utils.ServerErrorResponse(c, "数据库查询失败")
			return
		}

		// 防止循环引用（检查父分类的父分类链中是否包含当前分类）
		if isDescendant(config.DB, category.ID, *updateData.ParentID) {
			utils.ErrorResponse(c, 400, "不能设置子分类为父分类，会造成循环引用")
			return
		}

		category.ParentID = updateData.ParentID
	}

	// 更新其他字段
	if updateData.Description != "" {
		category.Description = updateData.Description
	}
	if updateData.Icon != "" {
		category.Icon = updateData.Icon
	}
	if updateData.Color != "" {
		category.Color = updateData.Color
	}
	category.Sort = updateData.Sort
	category.IsActive = updateData.IsActive

	if err := config.DB.Save(&category).Error; err != nil {
		utils.ServerErrorResponse(c, "分类更新失败")
		return
	}

	// 预加载父分类信息
	config.DB.Preload("Parent").First(&category, category.ID)

	utils.SuccessResponse(c, category)
}

// DeleteCategory 删除分类
func DeleteCategory(c *gin.Context) {
	categoryID := c.Param("id")

	var category models.Category
	if err := config.DB.First(&category, categoryID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "分类不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 检查是否有子分类
	var childCount int64
	config.DB.Model(&models.Category{}).Where("parent_id = ?", category.ID).Count(&childCount)
	if childCount > 0 {
		utils.ErrorResponse(c, 400, "存在子分类，无法删除")
		return
	}

	// 检查是否有文件使用此分类
	var fileCount int64
	config.DB.Model(&models.File{}).Where("category_id = ? AND is_deleted = ?", category.ID, false).Count(&fileCount)
	if fileCount > 0 {
		utils.ErrorResponse(c, 400, "存在文件使用此分类，无法删除")
		return
	}

	if err := config.DB.Delete(&category).Error; err != nil {
		utils.ServerErrorResponse(c, "分类删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "分类删除成功"})
}

// isDescendant 检查targetID是否是parentID的后代
func isDescendant(db *gorm.DB, targetID, parentID uint) bool {
	var parent models.Category
	if err := db.First(&parent, parentID).Error; err != nil {
		return false
	}

	// 如果父分类的ParentID为空，说明是根分类，不存在循环
	if parent.ParentID == nil {
		return false
	}

	// 如果父分类的ParentID等于目标ID，说明存在循环
	if *parent.ParentID == targetID {
		return true
	}

	// 递归检查上级分类
	return isDescendant(db, targetID, *parent.ParentID)
}
