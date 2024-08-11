package model

import (
  	"github.com/gin-generator/ginctl/package/time"
)

type Policies struct {
    Action string `json:"action" gorm:"column:action" validate:"required,max=6"`	// 操作 
    CreatedAt time.Time `json:"created_at" gorm:"column:created_at" validate:"omitempty,datetime"` 
    Effect string `json:"effect" gorm:"column:effect" validate:"required,max=5"`	// 影响：ALLOW-允许、DENY-拒绝 
    Id uint32 `json:"id" gorm:"column:id;primaryKey;autoIncrement" validate:"required,numeric"` 
    ResourceId uint32 `json:"resource_id" gorm:"column:resource_id" validate:"required,numeric"`	// 资源id 
    RoleId uint32 `json:"role_id" gorm:"column:role_id" validate:"required,numeric"`	// 角色id 
    UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" validate:"omitempty,datetime"` 
}


func (p *Policies) TableName() string {
	return "policies"
}

func NewPolicies() *Policies {
	return &Policies{}
}