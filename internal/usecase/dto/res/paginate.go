package res

type Paginate struct {
	Total int `json:"total" binding:"required"`
	Count int `json:"count" binding:"required"`
}
