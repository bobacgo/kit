package g

type ReqPageInfo struct {
	CurrPage int    `json:"currPage" form:"currPage"` // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
}

type GetById struct {
	ID int `json:"id" form:"id"` // 主键ID
}

func (r *GetById) UintID() uint {
	return uint(r.ID)
}

type IdsReq struct {
	Ids []int `json:"ids" form:"ids"`
}
