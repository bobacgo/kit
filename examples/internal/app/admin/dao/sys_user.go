package dao

import (
	"github.com/bobacgo/kit/examples/internal/app/admin/model"
	"github.com/bobacgo/kit/web/r/page"
)

type sysUserDaoImpl struct{}

func (dao *sysUserDaoImpl) PageList(query model.PageQuery) (*page.Data[model.SysUser], error) {
	_ = query
	users := []model.SysUser{
		{"weilanjin", "abc123", "lanjin.wei"},
		{"gogo", "abc123", "gogo"},
	}

	data := page.New(int64(len(users)), users...)
	return data, nil
}

func (dao *sysUserDaoImpl) Create(user model.SysUser) error {
	// TODO implement me
	panic("implement me")
}

func (dao *sysUserDaoImpl) Update(user model.SysUser) error {
	// TODO implement me
	panic("implement me")
}

func (dao *sysUserDaoImpl) Delete(id int) error {
	// TODO implement me
	panic("implement me")
}
