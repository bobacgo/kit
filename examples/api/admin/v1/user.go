package v1

type UserPageListReq struct {
	Keyword  string `json:"keyword"`
	Page     int    `json:"page" validate:"required"`
	PageSize int    `json:"page_size" validate:"required"`
}

type UserPageListResp struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}
