package controllers

import (
	"material-platform/config"
	"material-platform/models"
	"material-platform/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// Register 用户注册
func Register(c *gin.Context) {
	var req models.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已存在
	var existingUser models.User
	if err := config.DB.Where("username = ?", req.Username).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, 400, "用户名已存在")
		return
	}

	// 检查邮箱是否已存在
	if err := config.DB.Where("email = ?", req.Email).First(&existingUser).Error; err == nil {
		utils.ErrorResponse(c, 400, "邮箱已存在")
		return
	}

	// 加密密码
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		utils.ServerErrorResponse(c, "密码加密失败")
		return
	}

	// 创建用户
	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.ServerErrorResponse(c, "用户创建失败")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.ServerErrorResponse(c, "Token生成失败")
		return
	}

	utils.SuccessResponse(c, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// Login 用户登录
func Login(c *gin.Context) {
	var req models.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 查找用户
	var user models.User
	if err := config.DB.Where("username = ? OR email = ?", req.Username, req.Username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.ErrorResponse(c, 400, "用户名或密码错误")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 验证密码
	if !utils.CheckPasswordHash(req.Password, user.Password) {
		utils.ErrorResponse(c, 400, "用户名或密码错误")
		return
	}

	// 生成token
	token, err := utils.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		utils.ServerErrorResponse(c, "Token生成失败")
		return
	}

	utils.SuccessResponse(c, models.LoginResponse{
		Token: token,
		User:  user,
	})
}

// GetUserProfile 获取用户个人信息
func GetUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "用户不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	utils.SuccessResponse(c, user)
}

// UpdateUserProfile 更新用户个人信息
func UpdateUserProfile(c *gin.Context) {
	userID, _ := c.Get("user_id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "用户不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 绑定更新数据
	var updateData struct {
		Email  string `json:"email"`
		Avatar string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查邮箱是否已被其他用户使用
	if updateData.Email != "" && updateData.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", updateData.Email, user.ID).First(&existingUser).Error; err == nil {
			utils.ErrorResponse(c, 400, "邮箱已被使用")
			return
		}
		user.Email = updateData.Email
	}

	if updateData.Avatar != "" {
		user.Avatar = updateData.Avatar
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ServerErrorResponse(c, "用户信息更新失败")
		return
	}

	utils.SuccessResponse(c, user)
}

// GetAllUsers 获取所有用户（管理员功能）
func GetAllUsers(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var users []models.User
	var total int64

	// 计算总数
	config.DB.Model(&models.User{}).Count(&total)

	// 分页查询
	offset := (page - 1) * pageSize
	if err := config.DB.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	utils.PageResponse(c, users, total, page, pageSize)
}

// UpdateUserByAdmin 管理员更新用户信息
func UpdateUserByAdmin(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "用户不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	// 绑定更新数据
	var updateData struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Role     string `json:"role"`
		Avatar   string `json:"avatar"`
	}

	if err := c.ShouldBindJSON(&updateData); err != nil {
		utils.ErrorResponse(c, 400, "请求参数错误: "+err.Error())
		return
	}

	// 检查用户名是否已被使用
	if updateData.Username != "" && updateData.Username != user.Username {
		var existingUser models.User
		if err := config.DB.Where("username = ? AND id != ?", updateData.Username, user.ID).First(&existingUser).Error; err == nil {
			utils.ErrorResponse(c, 400, "用户名已被使用")
			return
		}
		user.Username = updateData.Username
	}

	// 检查邮箱是否已被使用
	if updateData.Email != "" && updateData.Email != user.Email {
		var existingUser models.User
		if err := config.DB.Where("email = ? AND id != ?", updateData.Email, user.ID).First(&existingUser).Error; err == nil {
			utils.ErrorResponse(c, 400, "邮箱已被使用")
			return
		}
		user.Email = updateData.Email
	}

	if updateData.Role != "" {
		user.Role = updateData.Role
	}

	if updateData.Avatar != "" {
		user.Avatar = updateData.Avatar
	}

	if err := config.DB.Save(&user).Error; err != nil {
		utils.ServerErrorResponse(c, "用户信息更新失败")
		return
	}

	utils.SuccessResponse(c, user)
}

// DeleteUser 删除用户（管理员功能）
func DeleteUser(c *gin.Context) {
	userID := c.Param("id")

	var user models.User
	if err := config.DB.First(&user, userID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			utils.NotFoundResponse(c, "用户不存在")
			return
		}
		utils.ServerErrorResponse(c, "数据库查询失败")
		return
	}

	if err := config.DB.Delete(&user).Error; err != nil {
		utils.ServerErrorResponse(c, "用户删除失败")
		return
	}

	utils.SuccessResponse(c, gin.H{"message": "用户删除成功"})
}

// GetSystemStats 获取系统统计信息（管理员功能）
func GetSystemStats(c *gin.Context) {
	var userCount, fileCount, categoryCount, tagCount int64

	config.DB.Model(&models.User{}).Count(&userCount)
	config.DB.Model(&models.File{}).Where("is_deleted = ?", false).Count(&fileCount)
	config.DB.Model(&models.Category{}).Count(&categoryCount)
	config.DB.Model(&models.Tag{}).Count(&tagCount)

	// 计算文件总大小
	var totalSize int64
	config.DB.Model(&models.File{}).Where("is_deleted = ?", false).Select("COALESCE(SUM(file_size), 0)").Scan(&totalSize)

	stats := gin.H{
		"user_count":           userCount,
		"file_count":           fileCount,
		"category_count":       categoryCount,
		"tag_count":            tagCount,
		"total_size":           totalSize,
		"total_size_formatted": utils.FormatFileSize(totalSize),
	}

	utils.SuccessResponse(c, stats)
}
