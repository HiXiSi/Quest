package models

import (
	"time"
)

// FileTag 文件标签关联模型
type FileTag struct {
	FileID    uint      `gorm:"primaryKey" json:"file_id"`
	TagID     uint      `gorm:"primaryKey" json:"tag_id"`
	CreatedAt time.Time `json:"created_at"`
	
	// 关联关系
	File      File      `gorm:"foreignKey:FileID" json:"file,omitempty"`
	Tag       Tag       `gorm:"foreignKey:TagID" json:"tag,omitempty"`
}

// TableName 指定表名
func (FileTag) TableName() string {
	return "file_tags"
}