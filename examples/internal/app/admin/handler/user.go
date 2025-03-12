package handler

import (
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	// svc biz.SysUser
}

func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// PageList 获取用户分页列表
// @Summary 用户管理
// @Description 获取用户分页列表
// @Tags 用户
// @Accept application/json
// @Produce application/json
// @Success 200 {object} []v1.UserPageListResp "查询成功"
// @Router /v1/user/pageList [post]
func (u *UserHandler) PageList(c *gin.Context) {
	panic("implement me")
}

func (u *UserHandler) Get(c *gin.Context) {
}

func (u *UserHandler) Create(c *gin.Context) {
}

func (u *UserHandler) Update(c *gin.Context) {

}

func (u *UserHandler) Delete(c *gin.Context) {

}