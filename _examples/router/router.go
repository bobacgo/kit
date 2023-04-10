package router

import (
	"github.com/gin-gonic/gin"
	v1 "github.com/gogoclouds/gogo/_examples/api/module_one/v1"
)

func LoadRouter(e *gin.Engine) {
	g := e.Group("v1")

	// sys user
	r := g.Group("user")
	r.POST("create", v1.SysUserApi.Create)
	r.PUT("update", v1.SysUserApi.Update)
	r.DELETE("delete", v1.SysUserApi.Delete)
	r.POST("pageList", v1.SysUserApi.PageList)

	// ....
}
