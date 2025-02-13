package req

type Pagination struct {
	Limit  *int `json:"limit" form:"limit"`
	Offset *int `json:"offset" form:"offset"`
}
