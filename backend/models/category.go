package models

import (
	"time"
	"gorm.io/gorm"
)

// Category 分类模型
type Category struct {
	ID          uint           `gorm:"primaryKey" json:"id"`
	Name        string         `gorm:"not null;size:100" json:"name"`
	Description string         `gorm:"size:500" json:"description"`
	ParentID    *uint          `gorm:"index" json:"parent_id"`
	Icon        string         `gorm:"size:100" json:"icon"`
	Color       string         `gorm:"size:20" json:"color"`
	Sort        int            `gorm:"default:0" json:"sort"`
	IsActive    bool           `gorm:"default:true" json:"is_active"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	
	// 关联关系
	Parent      *Category      `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Children    []Category     `gorm:"foreignKey:ParentID" json:"children,omitempty"`
	Files       []File         `gorm:"foreignKey:CategoryID" json:"files,omitempty"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}