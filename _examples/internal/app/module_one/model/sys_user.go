package model

import (
	"github.com/gogoclouds/gogo/web/r"
)

type SysUser struct {
	Username string `json:"username"`
	Passcode string `json:"-"`
	Nickname string `json:"nickname"`
}

type PageQuery struct {
	r.PageInfo
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type CreateUserReq struct {
	Username   string `json:"name" binding:"required"`
	Password   string `json:"password" binding:"required"`
	Nickname   string `json:"nickname"`
	Age        uint8  `json:"age" binding:"gte=1,lte=130"`
	Email      string `json:"email" binding:"required,email"`
	RePassword string `json:"rePassword" binding:"required,eqfield=Password"`
}