package core

import (
	"time"

	"gorm.io/gorm"
)

// Model 基础模型
type Model struct {
	ID        uint64         `json:"id" gorm:"column:id;primaryKey;not_null;autoIncrement;comment:ID"`
	DeletedAt gorm.DeletedAt `json:"deleted_at" gorm:"index;comment:是否已删除[null:否,非null:是]"`
	CreatedAt DateTime       `json:"created_at" gorm:"type:datetime;comment:创建日期"`
	UpdatedAt DateTime       `json:"updated_at" gorm:"type:datetime;comment:更新日期"`
}

// BeforeCreate 创建添加时间
func (b *Model) BeforeCreate(tx *gorm.DB) (err error) {
	b.CreatedAt = DateTime(time.Now())
	b.UpdatedAt = b.CreatedAt
	return
}

// BeforeUpdate 更新时间
func (b *Model) BeforeUpdate(tx *gorm.DB) (err error) {
	b.UpdatedAt = DateTime(time.Now())
	return
}
