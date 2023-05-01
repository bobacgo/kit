package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/_examples/internal/app/system/biz"
	"github.com/gogoclouds/gogo/_examples/internal/app/system/model"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/gin/valid"
	"github.com/gogoclouds/gogo/web/r"
)

type sysUserApi struct{}

func (api *sysUserApi) PageList(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.PageQuery](ctx)
	if !ok {
		return
	}
	pageResp, err := biz.SysUser.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessData(ctx, pageResp)
}

func (api *sysUserApi) Create(ctx *gin.Context) {
	req, ok := valid.ShouldBind[model.CreateUserReq](ctx)
	if !ok {
		return
	}
	_ = req
	reply.SuccessMsg(ctx, r.OKCreate)
}

func (api *sysUserApi) Update(ctx *gin.Context) {
	logger.Error("update fail")
	reply.FailMsg(ctx, r.FailUpdate)
	return
}

func (api *sysUserApi) Delete(ctx *gin.Context) {
	reply.Success(ctx)
}