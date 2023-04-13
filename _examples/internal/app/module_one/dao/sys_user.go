package dao

import (
	"github.com/gogoclouds/gogo/_examples/internal/app/module_one/model"
	"github.com/gogoclouds/gogo/g"
	"github.com/gogoclouds/gogo/web/r"
)

type sysUserDaoImpl struct{}

func (dao *sysUserDaoImpl) PageList(query model.PageQuery) (*r.PageResp[model.SysUser], *g.Error) {
	_ = query
	users := []model.SysUser{
		{"weilanjin", "abc123", "lanjin.wei"},
		{"gogo", "abc123", "gogo"},
	}
	return r.NewPage(users, 0, 2, 10), nil
}

func (dao *sysUserDaoImpl) Create(user model.SysUser) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (dao *sysUserDaoImpl) Update(user model.SysUser) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (dao *sysUserDaoImpl) Delete(id int) *g.Error {
	//TODO implement me
	panic("implement me")
}