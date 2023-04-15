package model

import (
	"github.com/gogoclouds/gogo/_examples/internal/app/common"
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
	Username   string             `json:"name" binding:"required"`
	Password   string             `json:"password" binding:"required"`
	RePassword string             `json:"rePassword" binding:"required,eqfield=Password"`
	Nickname   string             `json:"nickname"`
	Age        uint8              `json:"age" binding:"gte=1,lte=130"`
	Email      string             `json:"email" binding:"required,email"`
	Phone      string             `json:"phone" binding:"required,number,startswith=1,len=11"`
	Locality   []*common.Locality `json:"locality" binding:"required,gte=1"`
}