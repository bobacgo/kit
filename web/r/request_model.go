package r

type IdReq struct {
	ID string `json:"id" form:"id" binding:"required"`
}

type IdsReq struct {
	Ids []string `json:"ids" form:"ids" binding:"gte=1"`
}