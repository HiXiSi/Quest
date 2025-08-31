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

// CreateFormSchema 创建表单结构
func CreateFormSchema(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var req struct {
		Name        string             `json:"name" binding:"required"`
		Description string             `json:"description"`
		Fields      []models.FormField `json:"fields" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 验证字段 - 允许创建时没有字段，用户可以后续添加
	// if len(req.Fields) == 0 {
	// 	utils.ErrorResponse(c, 400, "表单必须包含至少一个字段")
	// 	return
	// }

	// 构建Schema数据
	schemaData := models.FormSchemaData{
		Fields: req.Fields,
	}

	schemaJSON, err := json.Marshal(schemaData)
	if err != nil {
		utils.ServerErrorResponse(c, "Schema序列化失败")
		return
	}

	// 创建表单结构
	formSchema := models.FormSchema{
		Name:        req.Name,
		Description: req.Description,
		Schema:      models.JSON(schemaJSON),
		UserID:      userID.(uint),
	}

	if err := config.DB.Create(&formSchema).Error; err != nil {
		utils.ServerErrorResponse(c, "创建表单结构失败")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&formSchema, formSchema.ID)

	utils.SuccessResponse(c, formSchema)
}

// GetFormSchemas 获取表单结构列表
func GetFormSchemas(c *gin.Context) {
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	keyword := c.Query("keyword")

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	var schemas []models.FormSchema
	var total int64

	query := config.DB.Model(&models.FormSchema{})

	// 非管理员只能看到自己的表单
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	// 关键词搜索
	if keyword != "" {
		query = query.Where("name LIKE ? OR description LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}

	// 获取总数
	query.Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := query.Preload("User").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&schemas).Error; err != nil {
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	utils.SuccessResponse(c, gin.H{
		"list":       schemas,
		"total":      total,
		"page":       page,
		"page_size":  pageSize,
		"total_page": (total + int64(pageSize) - 1) / int64(pageSize),
	})
}

// GetFormSchema 获取单个表单结构
func GetFormSchema(c *gin.Context) {
	schemaID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var schema models.FormSchema
	query := config.DB.Where("id = ?", schemaID)

	// 非管理员只能访问自己的表单
	if role != "admin" {
		query = query.Where("user_id = ?", userID)
	}

	if err := query.Preload("User").First(&schema).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "表单结构不存在")
			return
		}
		utils.ServerErrorResponse(c, "查询失败")
		return
	}

	utils.SuccessResponse(c, schema)
}

// UpdateFormSchema 更新表单结构
func UpdateFormSchema(c *gin.Context) {
	schemaID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	var req struct {
		Name        string             `json:"name" binding:"required"`
		Description string             `json:"description"`
		Fields      []models.FormField `json:"fields" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 查找表单结构
	var schema models.FormSchema
	query := config.DB.Where("id = ?", schemaID)

	// 非管理员只能修改自己的表单
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

	// 构建新的Schema数据
	schemaData := models.FormSchemaData{
		Fields: req.Fields,
	}

	schemaJSON, err := json.Marshal(schemaData)
	if err != nil {
		utils.ServerErrorResponse(c, "Schema序列化失败")
		return
	}

	// 更新
	schema.Name = req.Name
	schema.Description = req.Description
	schema.Schema = models.JSON(schemaJSON)

	if err := config.DB.Save(&schema).Error; err != nil {
		utils.ServerErrorResponse(c, "更新失败")
		return
	}

	// 预加载用户信息
	config.DB.Preload("User").First(&schema, schema.ID)

	utils.SuccessResponse(c, schema)
}

// DeleteFormSchema 删除表单结构
func DeleteFormSchema(c *gin.Context) {
	schemaID := c.Param("id")
	userID, _ := c.Get("user_id")
	role, _ := c.Get("role")

	// 查找表单结构
	var schema models.FormSchema
	query := config.DB.Where("id = ?", schemaID)

	// 非管理员只能删除自己的表单
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

	// 检查是否有关联的记录
	var recordCount int64
	config.DB.Model(&models.FormRecord{}).Where("schema_id = ?", schemaID).Count(&recordCount)

	if recordCount > 0 {
		utils.ErrorResponse(c, 400, "无法删除：该表单结构下还有数据记录")
		return
	}

	// 删除表单结构
	if err := config.DB.Delete(&schema).Error; err != nil {
		utils.ServerErrorResponse(c, "删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "删除成功"})
}
