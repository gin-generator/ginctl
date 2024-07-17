package model

import (
	"github.com/gin-generator/ginctl/package/time"
)

type Resources struct {
	Id        uint32    `json:"id" gorm:"column:id;primaryKey;autoIncrement" validate:"required,numeric"`
	Name      string    `json:"name" gorm:"column:name" validate:"required,max=255"` // 资源名称
	Type      uint8     `json:"type" gorm:"column:type" validate:"required,numeric"` // 资源类型
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at" validate:"omitempty,datetime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" validate:"omitempty,datetime"`
}

func (r *Resources) TableName() string {
	return "resources"
}

func NewResources() *Resources {
	return &Resources{}
}
