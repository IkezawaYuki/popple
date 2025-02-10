package req

type RegisterCustomer struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	WordpressURL string `json:"wordpress_url"`
}
