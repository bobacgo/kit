package biz

import (
	"github.com/bobacgo/kit/examples/internal/app/admin/dao"
	"github.com/bobacgo/kit/examples/internal/app/admin/model"
	"github.com/bobacgo/kit/g"
	"github.com/bobacgo/kit/web/r/page"
)

type sysUserServiceImpl struct{}

func (svc *sysUserServiceImpl) PageList(query model.PageQuery) (*page.Data[model.SysUser], *g.Error) {
	list, err := dao.SysUser.PageList(query)
	return list, err
}

func (svc *sysUserServiceImpl) Create(user model.SysUser) *g.Error {
	// TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Update(user model.SysUser) *g.Error {
	// TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Delete(id int) *g.Error {
	// TODO implement me
	panic("implement me")
}