package req

type PostQuery struct {
	ID         *int `form:"id"`
	CustomerID *int `form:"customer_id"`
	Pagination
}
