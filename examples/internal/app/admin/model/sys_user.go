package model

import (
	"github.com/bobacgo/kit/web/orm"
	"github.com/bobacgo/kit/web/r/page"
)

type SysUser struct {
	Username string `json:"username"`
	Passcode string `json:"-"`
	Nickname string `json:"nickname"`
}

type PageQuery struct {
	page.Query
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}

type CreateUserReq struct {
	Username   string        `json:"name" binding:"required"`
	Password   string        `json:"password" binding:"required"`
	RePassword string        `json:"rePassword" binding:"required,eqfield=Password"`
	Nickname   string        `json:"nickname"`
	Birthday   orm.LocalTime `json:"birthday"`
	Gender     uint8         `json:"gender" binding:"lte=3"` // 0-未知|1-女|2-男
	Email      string        `json:"email" binding:"required,email"`
	Phone      string        `json:"phone" binding:"required,number,startswith=1,len=11"`
}