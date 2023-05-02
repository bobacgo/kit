package r

type PageInfo struct {
	Page     int `json:"page" form:"page" binding:"gte=1"`         // 页码
	PageSize int `json:"pageSize" form:"pageSize" binding:"gte=1"` // 每页大小
}

type IdReq struct {
	ID string `json:"id" form:"id" binding:"required"`
}

type IdsReq struct {
	Ids []string `json:"ids" form:"ids" binding:"gte=1"`
}