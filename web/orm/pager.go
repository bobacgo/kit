package orm

import (
	"github.com/gogoclouds/gogo/web/r"
	"gorm.io/gorm"
)

// Paginate 分页器
func Paginate(page r.PageInfo) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if page.Page <= 0 {
			page.Page = 1
		}
		if page.PageSize <= 0 {
			page.PageSize = 10
		}
		offset := (page.Page - 1) * page.PageSize
		return db.Offset(int(offset)).Limit(int(page.PageSize))
	}
}

// PageFind 分页查找
func PageFind[T any](db *gorm.DB, page r.PageInfo) (data *r.PageResp[T], err error) {
	var total int64
	if err = db.Count(&total).Error; err != nil && total == 0 {
		return
	}
	var list []T
	err = db.Scopes(Paginate(page)).Find(&list).Error
	return r.NewPage(list, page.Page, int(total), page.PageSize), err
}