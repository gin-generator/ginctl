package model

import (
  	"github.com/gin-generator/ginctl/package/time"
)

type Roles struct {
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at" validate:"omitempty,datetime"` 
    Id uint32 `json:"id" gorm:"column:id;primaryKey;autoIncrement" validate:"required,numeric"` 
    Name string `json:"name" gorm:"column:name" validate:"required,max=255"`	// 角色名称 
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" validate:"omitempty,datetime"` 
}


func (r *Roles) TableName() string {
	return "roles"
}

func NewRoles() *Roles {
	return &Roles{}
}