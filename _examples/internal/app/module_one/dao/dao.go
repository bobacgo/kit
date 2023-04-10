package dao

import (
	"github.com/gogoclouds/gogo/_examples/internal/app/module_one/model"
	"github.com/gogoclouds/gogo/g"
)

var SysUser g.IBase[model.SysUser, model.ReqPageQuery] = new(sysUserDaoImpl)
