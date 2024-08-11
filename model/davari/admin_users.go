package model

import (
  	"github.com/gin-generator/ginctl/package/time"
)

type AdminUsers struct {
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at" validate:"omitempty,datetime"` 
    Id int32 `json:"id" gorm:"column:id;primaryKey;autoIncrement" validate:"required,numeric"` 
    Password string `json:"password" gorm:"column:password" validate:"required,max=255"`	// 密码 
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" validate:"omitempty,datetime"` 
    Username string `json:"username" gorm:"column:username" validate:"required,max=255"`	// 用户名称 
}


func (a *AdminUsers) TableName() string {
	return "admin_users"
}

func NewAdminUsers() *AdminUsers {
	return &AdminUsers{}
}