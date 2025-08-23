package models

import (
	"time"
	"gorm.io/gorm"
)

// User 用户模型
type User struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Username    string         `gorm:"unique;not null;size:50" json:"username"`
	Email       string         `gorm:"unique;not null;size:100" json:"email"`
	Password    string         `gorm:"not null;size:255" json:"-"`
	Role        string         `gorm:"default:user;size:20" json:"role"` // admin, user
	Avatar      string         `gorm:"size:255" json:"avatar"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联关系
	Files       []File         `gorm:"foreignKey:UserID" json:"files,omitempty"`
}

// TableName 指定表名
func (User) TableName() string {
	return "users"
}