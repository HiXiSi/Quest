package controllers

import (
	"encoding/json"
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateFormRecord 创建表单数据记录
func CreateFormRecord(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req struct {
		SchemaID uint                   `json:"schema_id" binding:"required"`
		Data     map[string]interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证表单结构是否存在
	var schema models.FormSchema
	query := config.DB.Where("id = ?", req.SchemaID)

	// 非管理员只能在自己的表单中创建记录
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&schema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "表单结构不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	// 验证数据格式
	if err := validateFormData(schema.Schema, req.Data); err != nil {
		utils.ErrorResponse(c, 400, "数据验证失败: "+err.Error())
		return
	}

	// 序列化数据
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		utils.ServerErrorResponse(c, "数据序列化失败")
		return
	}

	// 创建记录
	record := models.FormRecord{
		SchemaID: req.SchemaID,
		Data:     models.JSON(dataJSON),
		UserID:   userID.(uint),
	}

	if err := config.DB.Create(&record).Error; err != nil {
		utils.ServerErrorResponse(c, "创建记录失败")
		return
	}

	// 预加载关联数据
	config.DB.Preload("Schema").Preload("User").First(&record, record.ID)

	utils.SuccessResponse(c, record)
}

// GetFormRecords 获取表单数据记录列表
func GetFormRecords(c *gin.Context) {
	schemaID := c.Param("id")
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

	// 验证表单结构权限
	var schema models.FormSchema
	query := config.DB.Where("id = ?", schemaID)

	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&schema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "表单结构不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	var records []models.FormRecord
	var total int64

	recordQuery := config.DB.Model(&models.FormRecord{}).Where("schema_id = ?", schemaID)

	// 获取总数
	recordQuery.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := recordQuery.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&records).Error; err != nil {
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":       records,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
		"schema":     schema,
	})
}

// GetFormRecord 获取单个表单数据记录
func GetFormRecord(c *gin.Context) {
	recordID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var record models.FormRecord
	query := config.DB.Where("id = ?", recordID)

	// 非管理员只能访问自己创建的记录或自己表单下的记录
	if role != "admin" {
		query = query.Where("user_id = ? OR schema_id IN (?)",
			userID,
			config.DB.Model(&models.FormSchema{}).Select("id").Where("user_id = ?", userID))
	}

	if err := query.Preload("Schema").Preload("User").First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "记录不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	utils.SuccessResponse(c, record)
}

// UpdateFormRecord 更新表单数据记录
func UpdateFormRecord(c *gin.Context) {
	recordID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req struct {
		Data map[string]interface{} `json:"data" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 查找记录
	var record models.FormRecord
	query := config.DB.Where("id = ?", recordID)

	// 非管理员只能修改自己创建的记录
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Preload("Schema").First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "记录不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	// 验证数据格式
	if err := validateFormData(record.Schema.Schema, req.Data); err != nil {
		utils.ErrorResponse(c, 400, "数据验证失败: "+err.Error())
		return
	}

	// 序列化数据
	dataJSON, err := json.Marshal(req.Data)
	if err != nil {
		utils.ServerErrorResponse(c, "数据序列化失败")
		return
	}

	// 更新记录
	record.Data = models.JSON(dataJSON)

	if err := config.DB.Save(&record).Error; err != nil {
		utils.ServerErrorResponse(c, "更新失败")
		return
	}

	// 重新加载记录
	config.DB.Preload("Schema").Preload("User").First(&record, record.ID)

	utils.SuccessResponse(c, record)
}

// DeleteFormRecord 删除表单数据记录
func DeleteFormRecord(c *gin.Context) {
	recordID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// 查找记录
	var record models.FormRecord
	query := config.DB.Where("id = ?", recordID)

	// 非管理员只能删除自己创建的记录
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.First(&record).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "记录不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	// 删除记录
	if err := config.DB.Delete(&record).Error; err != nil {
		utils.ServerErrorResponse(c, "删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "删除成功"})
}

// validateFormData 验证表单数据
func validateFormData(schema models.JSON, data map[string]interface{}) error {
	// 解析Schema
	var schemaData models.FormSchemaData
	if err := json.Unmarshal(schema, &schemaData); err != nil {
		return err
	}

	// 验证必填字段
	for _, field := range schemaData.Fields {
		if field.Required {
			if value, exists := data[field.ID]; !exists || value == nil || value == "" {
				return gin.Error{
					Err:  nil,
					Type: gin.ErrorTypePublic,
					Meta: "字段 '" + field.Label + "' 是必填的",
				}
			}
		}
	}

	return nil
}
