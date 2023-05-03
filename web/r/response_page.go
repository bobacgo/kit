package r

type PageInfo struct {
	Page     uint `json:"page" form:"page" binding:"gte=1"`         // 页码
	PageSize uint `json:"pageSize" form:"pageSize" binding:"gte=1"` // 每页大小
}

type PageData[T any] struct {
	Total int `json:"total"` // 总页数
	List  []T `json:"list"`  // 列表数据
}

// PageResp 分页数据响应体
// T 列表每一项的数据类型
type PageResp[T any] struct {
	PageInfo
	PageData[T]
}

// PageAnyResp 分页数据响应体
// 使用场景：匿名结构体返回的时候
type PageAnyResp struct {
	PageInfo
	Total int `json:"total"` // 总页数
	List  any `json:"list"`  // 列表数据
}

// PageMetaResp 分页数据响应体 （携带额外数据）
// T 列表每一项的数据类型
// M 非列表数据的数据类型
type PageMetaResp[T any, M any] struct {
	PageInfo
	PageData[T]
	Meta M `json:"meta"`
}

// PageAnyMetaResp 分页数据响应体 （携带额外数据）
// M 非列表数据的数据类型
type PageAnyMetaResp[M any] struct {
	PageInfo
	Total int `json:"total"` // 总页数
	List  any `json:"list"`  // 列表数据
	Meta  M   `json:"meta"`
}

// NewPage 分页数据组装
//
// page 当前数据是第几页
// total 总的条数
// pageSize 每一页多少条数据
func NewPage[T any](list []T, total int, page, pageSize uint) *PageResp[T] {
	return &PageResp[T]{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		PageData: PageData[T]{total, list},
	}
}

// NewPageMeta 分页数据组装-携带非列表数据
//
// page 当前数据是第几页
// total 总的条数
// pageSize 每一页多少条数据
// meta 非列表数据
func NewPageMeta[T any, M any](list []T, total int, page, pageSize uint, meta M) *PageMetaResp[T, M] {
	return &PageMetaResp[T, M]{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		PageData: PageData[T]{total, list},
		Meta:     meta,
	}
}

func NewPageAny(list any, total int, page, pageSize uint) *PageAnyResp {
	return &PageAnyResp{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		Total:    total, // 总页数
		List:     list,  // 列表数据
	}
}

func NewPageAnyMeta[M any](list any, total int, page, pageSize uint, meta M) *PageAnyMetaResp[M] {
	return &PageAnyMetaResp[M]{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		Total:    total, // 总页数
		List:     list,  // 列表数据
		Meta:     meta,
	}
}