package models

import (
	"time"
	"gorm.io/gorm"
)

// Tag 标签模型
type Tag struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"unique;not null;size:50" json:"name"`
	Color       string         `gorm:"size:20" json:"color"`
	Description string         `gorm:"size:500" json:"description"`
	UsageCount  int            `gorm:"default:0" json:"usage_count"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联关系
	Files       []File         `gorm:"many2many:file_tags;" json:"files,omitempty"`
}

// TableName 指定表名
func (Tag) TableName() string {
	return "tags"
}