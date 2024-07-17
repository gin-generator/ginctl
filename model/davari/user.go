package model

import (
	"github.com/gin-generator/ginctl/package/time"
)

type User struct {
	Id        int32     `json:"id" gorm:"column:id;primaryKey;autoIncrement" validate:"required,numeric"`
	Username  string    `json:"username" gorm:"column:username" validate:"required,max=100"` // 用户名称
	Password  string    `json:"password" gorm:"column:password" validate:"required,max=60"`  // 密码
	Avatar    string    `json:"avatar" gorm:"column:avatar" validate:"required,max=255"`     // 头像
	Phone     string    `json:"phone" gorm:"column:phone" validate:"required,max=11"`        // 电话号码
	Email     string    `json:"email" gorm:"column:email" validate:"required,max=255"`       // 邮箱
	Gender    uint8     `json:"gender" gorm:"column:gender" validate:"required,numeric"`     // 性别：0-未知、1-男、2-女
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at" validate:"omitempty,datetime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"column:updated_at" validate:"omitempty,datetime"`
}

func (u *User) TableName() string {
	return "user"
}

func NewUser() *User {
	return &User{}
}
