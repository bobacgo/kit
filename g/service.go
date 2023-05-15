package g

import (
	"errors"
	"github.com/gogoclouds/gogo/web/r"
	"golang.org/x/exp/maps"
	"gorm.io/gorm"
)

type FindByIDService[T any] struct{}

func (*FindByIDService[T]) FindByID(tx *gorm.DB, id string) (T, *Error) {
	var m T
	err := tx.Where("id = ?", id).First(&m).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return m, WrapError(err, r.FailRecordNotFound)
	}
	return m, WrapError(err, r.FailRead)
}

// UniqueService 校验传入值是否已经在数据库中存在
// T 表示要查询表的Model
// k为 'id' 作为排除，比如更新时
//
//	type UserService struct{
//		g.Unique[model.User]
//	}
type UniqueService[T any] struct{}

func (*UniqueService[T]) Verify(tx *gorm.DB, q map[string]any) *Error {
	tx = tx.Model(new(T)).Select(maps.Keys(q))
	for k, v := range q {
		if k != "id" {
			tx = tx.Or(k, v) // 传入的查询条件用 or 请注意索引是否生效
		}
	}
	var list []map[string]any
	if err := tx.Find(&list).Error; err != nil {
		return WrapError(err, "校验唯一出错")
	}
	if len(list) > 0 {
		msg := make(map[string]any, 0)
		for _, m := range list {
			if q["id"] != nil && q["id"] == m["id"] { // 注意：id 值类型 string, 更新唯一校验，排除自身
				continue
			}
			for k, v := range q {
				if v == m[k] {
					msg[k] = v
				}
			}
		}
		if len(msg) > 0 {
			return &Error{
				Text: ErrRecordRepeat.Error(),
				Misc: msg,
			}
		}
	}
	return nil
}
