package orm

import (
	"github.com/gogoclouds/gogo/pkg/util"
	"gorm.io/gorm"
)

// 关联关系表命名规范
// otm 一对多
// mto 多对一
// mtm 多对多

// Model
// ID 值使用UUID, 避免分布式环境下key冲突
//
// LocalTime
//
//	1.可以通过配置指定格式 app.timeFormat
//	1.1.string -> time (指定根式序列化)
//	1.2.time -> string (指定根式反序列化)
type Model struct {
	ID        string         `json:"id" gorm:"primarykey"`
	CreatedAt LocalTime      `json:"createdAt"`
	UpdatedAt LocalTime      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

func (u *Model) BeforeCreate(tx *gorm.DB) error {
	u.ID = util.UUID()
	return nil
}