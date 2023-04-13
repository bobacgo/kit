package r

type Page[T any] struct {
	Total int `json:"total"` // 总页数
	List  []T `json:"list"`  // 列表数据
}

// PageResp 分页数据响应体
// T 列表每一项的数据类型
type PageResp[T any] struct {
	PageInfo
	Page[T]
}

// PageMetaResp 分页数据响应体 （携带额外数据）
// T 列表每一项的数据类型
// M 非列表数据的数据类型
type PageMetaResp[T any, M any] struct {
	PageInfo
	Page[T]
	Meta M `json:"meta"`
}

// NewPage 分页数据组装
//
// page 当前数据是第几页
// total 总的条数
// pageSize 每一页多少条数据
func NewPage[T any](list []T, page, total, pageSize int) *PageResp[T] {
	return &PageResp[T]{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		Page:     Page[T]{total, list},
	}
}

// NewPageMeta 分页数据组装-携带非列表数据
//
// page 当前数据是第几页
// total 总的条数
// pageSize 每一页多少条数据
// meta 非列表数据
func NewPageMeta[T any, M any](list []T, page, total, pageSize int, meta M) *PageMetaResp[T, M] {
	return &PageMetaResp[T, M]{
		PageInfo: PageInfo{Page: page, PageSize: pageSize},
		Page:     Page[T]{total, list},
		Meta:     meta,
	}
}