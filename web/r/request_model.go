package r

type PageInfo struct {
	Page     uint   `json:"page" form:"page" binding:"gte=1"`         // 页码
	PageSize uint   `json:"pageSize" form:"pageSize" binding:"gte=1"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`                   //关键字
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
