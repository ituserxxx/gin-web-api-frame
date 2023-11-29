package model

import (
	"gorm.io/gorm"
)

type Query struct {
	Page  int `form:"page" binding:"required,gte=1"`
	Limit int `form:"limit" binding:"required,oneof=10 20 30"`
}
type CreateUser struct {
	gorm.Model
	ID       uint
	Name     string
	Phone    string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Password string `json:"password"`
}

// CreateUserReq 创建用户
type CreateUserReq struct {
	ID     uint
	Name   string
	Phone  string `valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Passwd string `json:"password"`
}

type DeleteUserReq struct {
	ID uint
}

// LoginReq 用户登录
type LoginReq struct {
	Phone  string `json:"phone" valid:"matches(^1[3-9]{1}\\d{9}$)"`
	Passwd string `json:"password" binding:"required"`
}

// 更改密码
type ChangePwdReq struct {
	OldPasswd string `json:"oldPassword" binding:"required"`
	NewPasswd string `json:"newPassword" binding:"required"`
}
