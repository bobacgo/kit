package biz

import (
	"github.com/bobacgo/kit/examples/internal/app/admin/model"
	"github.com/bobacgo/kit/g"
)

var SysUser g.IBase[model.SysUser, model.PageQuery] = new(sysUserServiceImpl)