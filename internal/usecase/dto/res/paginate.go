package res

type Paginate struct {
	Total int64 `json:"total" binding:"required"`
	Count int64 `json:"count" binding:"required"`
}
