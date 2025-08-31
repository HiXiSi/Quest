package models

import (
	"time"
)

// FormRecord 表单数据记录
type FormRecord struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	SchemaID  uint      `json:"schema_id" gorm:"not null"`
	Data      JSON      `json:"data" gorm:"type:json;not null"`
	UserID    uint      `json:"user_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// 关联
	Schema FormSchema `json:"schema,omitempty" gorm:"foreignKey:SchemaID"`
	User   User       `json:"user,omitempty" gorm:"foreignKey:UserID"`
}

// TableName 指定表名
func (FormRecord) TableName() string {
	return "form_records"
}

// FormRecordResponse 表单记录响应结构
type FormRecordResponse struct {
	ID        uint                   `json:"id"`
	SchemaID  uint                   `json:"schema_id"`
	Data      map[string]interface{} `json:"data"`
	UserID    uint                   `json:"user_id"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
	Schema    *FormSchema            `json:"schema,omitempty"`
}
