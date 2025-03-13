package biz

import (
	"github.com/bobacgo/kit/examples/internal/app/admin/dao"
	"github.com/bobacgo/kit/examples/internal/app/admin/model"
	"github.com/bobacgo/kit/web/r/page"
)

type sysUserServiceImpl struct{}

func (svc *sysUserServiceImpl) PageList(query model.PageQuery) (*page.Data[model.SysUser], error) {
	list, err := dao.SysUser.PageList(query)
	return list, err
}

func (svc *sysUserServiceImpl) Create(user model.SysUser) error {
	// TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Update(user model.SysUser) error {
	// TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Delete(id int) error {
	// TODO implement me
	panic("implement me")
}
