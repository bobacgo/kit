package orm

import (
	"github.com/bobacgo/kit/web/r/page"
	"gorm.io/gorm"
)

// Paginate 分页器
func Paginate(q page.Query) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Offset(q.Offset()).Limit(q.Limit())
	}
}

// PageFind 分页查找
func PageFind[T any](db *gorm.DB, q page.Query) (data *page.Data[T], err error) {
	data = page.New[T](0)

	var total int64
	if err = db.Count(&total).Error; err != nil || total == 0 {
		return
	}
	err = db.Scopes(Paginate(q)).Find(&data.List).Error
	data.Total = total
	return
}
