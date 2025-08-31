package config

import (
	"log"
	"material-platform/models"
	"material-platform/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitDB 初始化数据库连接
func InitDB() {
	var err error

	// 连接SQLite数据库
	DB, err = gorm.Open(sqlite.Open("../static/database.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("数据库连接失败:", err)
	}

	// 自动迁移数据表
	err = DB.AutoMigrate(
		&models.User{},
		&models.Category{},
		&models.Tag{},
		&models.File{},
		&models.FileTag{},
		&models.FormSchema{},
		&models.FormRecord{},
	)
	if err != nil {
		log.Fatal("数据表迁移失败:", err)
	}

	// 创建默认管理员账号
	createDefaultAdmin()

	log.Println("数据库初始化成功")
}

// createDefaultAdmin 创建默认管理员账号
func createDefaultAdmin() {
	// 检查是否已存在管理员账号
	var adminCount int64
	DB.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount)

	if adminCount > 0 {
		return // 已存在管理员，无需创建
	}

	// 创建默认管理员
	hashedPassword, err := utils.HashPassword("admin123")
	if err != nil {
		log.Printf("管理员密码加密失败: %v", err)
		return
	}

	admin := models.User{
		Username: "admin",
		Email:    "admin@material-platform.com",
		Password: hashedPassword,
		Role:     "admin",
	}

	if err := DB.Create(&admin).Error; err != nil {
		log.Printf("创建默认管理员失败: %v", err)
		return
	}

	log.Println("默认管理员账号创建成功")
	log.Println("管理员用户名: admin")
	log.Println("管理员密码: admin123")
	log.Println("请及时修改默认密码！")
}
