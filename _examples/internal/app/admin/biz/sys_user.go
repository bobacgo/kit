package biz

import (
	"github.com/bobacgo/kit/_examples/internal/app/admin/dao"
	"github.com/bobacgo/kit/_examples/internal/app/admin/model"
	"github.com/bobacgo/kit/g"
	"github.com/bobacgo/kit/web/r"
)

type sysUserServiceImpl struct{}

func (svc *sysUserServiceImpl) PageList(query model.PageQuery) (*r.PageResp[model.SysUser], *g.Error) {
	list, err := dao.SysUser.PageList(query)
	return list, err
}

func (svc *sysUserServiceImpl) Create(user model.SysUser) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Update(user model.SysUser) *g.Error {
	//TODO implement me
	panic("implement me")
}

func (svc *sysUserServiceImpl) Delete(id int) *g.Error {
	//TODO implement me
	panic("implement me")
}
