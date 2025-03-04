package dao

import (
	"github.com/bobacgo/kit/_examples/internal/app/admin/model"
	"github.com/bobacgo/kit/g"
)

var SysUser g.IBase[model.SysUser, model.PageQuery] = new(sysUserDaoImpl)
