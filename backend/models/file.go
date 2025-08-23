package models

import (
	"time"
)

// File 文件模型
type File struct {
	ID           uint   `gorm:"primaryKey" json:"id"`
	OriginalName string `gorm:"not null;size:255" json:"original_name"`
	FileName     string `gorm:"not null;size:255" json:"file_name"`
	FilePath     string `gorm:"not null;size:500" json:"file_path"`
	FileSize     int64  `gorm:"not null" json:"file_size"`
	FileType     string `gorm:"not null;size:100" json:"file_type"`
	MimeType     string `gorm:"not null;size:100" json:"mime_type"`
	MD5Hash      string `gorm:"size:32;index" json:"md5_hash"`
	SHA256Hash   string `gorm:"size:64;index" json:"sha256_hash"`

	// 文件元数据
	Width    int     `json:"width,omitempty"`
	Height   int     `json:"height,omitempty"`
	Duration float64 `json:"duration,omitempty"`

	// 业务字段
	Description   string `gorm:"size:1000" json:"description"`
	IsPublic      bool   `gorm:"default:false" json:"is_public"`
	DownloadCount int    `gorm:"default:0" json:"download_count"`
	ViewCount     int    `gorm:"default:0" json:"view_count"`

	// 回收站
	IsDeleted bool       `gorm:"default:false;index" json:"is_deleted"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`

	// 外键
	UserID     uint  `gorm:"not null;index" json:"user_id"`
	CategoryID *uint `gorm:"index" json:"category_id"`

	// 时间戳
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联关系
	User     User      `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Category *Category `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	Tags     []Tag     `gorm:"many2many:file_tags;" json:"tags,omitempty"`
}

// TableName 指定表名
func (File) TableName() string {
	return "files"
}
