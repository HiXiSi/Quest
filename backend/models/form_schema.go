package models

import (
	"encoding/json"
	"time"
)

// FormSchema 表单结构定义
type FormSchema struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Name        string    `json:"name" gorm:"not null;size:255"`
	Description string    `json:"description"`
	Schema      JSON      `json:"schema" gorm:"type:json;not null"`
	UserID      uint      `json:"user_id" gorm:"not null"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`

	// 关联
	User    User         `json:"user" gorm:"foreignKey:UserID"`
	Records []FormRecord `json:"records,omitempty" gorm:"foreignKey:SchemaID"`
}

// JSON 自定义JSON类型，用于处理JSON字段
type JSON json.RawMessage

// Scan 实现scanner接口
func (j *JSON) Scan(value interface{}) error {
	if value == nil {
		*j = nil
		return nil
	}

	switch v := value.(type) {
	case []byte:
		*j = JSON(v)
	case string:
		*j = JSON(v)
	default:
		return nil
	}
	return nil
}

// Value 实现driver.Valuer接口
func (j JSON) Value() (interface{}, error) {
	if len(j) == 0 {
		return nil, nil
	}
	return string(j), nil
}

// MarshalJSON 实现json.Marshaler接口
func (j JSON) MarshalJSON() ([]byte, error) {
	if len(j) == 0 {
		return []byte("null"), nil
	}
	return j, nil
}

// UnmarshalJSON 实现json.Unmarshaler接口
func (j *JSON) UnmarshalJSON(data []byte) error {
	if j == nil {
		return nil
	}
	*j = JSON(data)
	return nil
}

// FormField 表单字段定义结构
type FormField struct {
	// 基础属性
	ID           string `json:"id"`
	Name         string `json:"name"`                    // 字段名称
	Label        string `json:"label"`                   // 字段标签
	Type         string `json:"type"`                    // 字段类型
	Required     bool   `json:"required"`                // 是否必填
	Placeholder  string `json:"placeholder,omitempty"`   // 占位符
	DefaultValue string `json:"default_value,omitempty"` // 默认值
	DbType       string `json:"dbType,omitempty"`        // 数据库类型

	// 格式属性
	// 唯一ID类型专用
	IdType string `json:"id_type,omitempty"` // auto_increment 或 uuid

	// 字符串类型专用
	Format    string `json:"format,omitempty"`     // text, email, phone, url, password
	MinLength *int   `json:"min_length,omitempty"` // 最小长度
	MaxLength *int   `json:"max_length,omitempty"` // 最大长度

	// 数值类型专用
	MinValue  *float64 `json:"min_value,omitempty"` // 最小值
	MaxValue  *float64 `json:"max_value,omitempty"` // 最大值
	Precision *int     `json:"precision,omitempty"` // 小数位数(仅浮点数)

	// 时间类型专用
	TimeFormat string `json:"time_format,omitempty"` // date_object, datetime, date, time

	// 枚举类型专用
	EnumOptions []FormFieldOption `json:"enum_options,omitempty"` // 枚举选项

	// 表单形式
	InputType    string `json:"input_type,omitempty"`    // 输入方式
	TextareaRows *int   `json:"textarea_rows,omitempty"` // 多行文本行数

	// 兼容旧字段
	Options    []FormFieldOption      `json:"options,omitempty"`
	Validation map[string]interface{} `json:"validation,omitempty"`
	SortOrder  int                    `json:"sort_order"`
}

// FormFieldOption 表单字段选项（用于select、radio、checkbox等）
type FormFieldOption struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// FormSchemaData 完整的表单结构数据
type FormSchemaData struct {
	Fields []FormField `json:"fields"`
}

// TableName 指定表名
func (FormSchema) TableName() string {
	return "form_schemas"
}
