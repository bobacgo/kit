package model

import "github.com/gogoclouds/gogo/g"

type SysUser struct {
	Username string `json:"username"`
	Passcode string `json:"-"`
	Nickname string `json:"nickname"`
}

type ReqPageQuery struct {
	g.ReqPageInfo
	Username string `json:"username"`
	Nickname string `json:"nickname"`
}
