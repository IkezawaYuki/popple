package req

type Post struct {
	ID         *int `form:"id"`
	CustomerID *int `form:"customer_id"`
	Pagination
}
