package common

type Locality struct {
	Province string `json:"province"` // 省
	City     string `json:"city"`     // 市
	District string `json:"district"` // 区
}