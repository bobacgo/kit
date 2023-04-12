package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/gogoclouds/gogo/_examples/internal/app/module_one/biz"
	"github.com/gogoclouds/gogo/_examples/internal/app/module_one/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/logger"
	"github.com/gogoclouds/gogo/web/gin/reply"
	"github.com/gogoclouds/gogo/web/r"
)

type sysUserApi struct{}

func (api *sysUserApi) PageList(ctx *gin.Context) {
	var req model.ReqPageQuery
	_ = ctx.ShouldBindJSON(&req)

	pageResp, err := biz.SysUser.PageList(req)
	if err != nil {
		logger.Error(err.Error())
		reply.FailMsg(ctx, err.Text)
		return
	}
	reply.SuccessData(ctx, pageResp)
}

func (api *sysUserApi) Create(ctx *gin.Context) {
	req, check := g.BindAndValidate[model.CreateUserReq](ctx)
	if check {
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