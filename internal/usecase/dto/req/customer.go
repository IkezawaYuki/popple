package req

import "time"

type CustomerQuery struct {
	IsFacebookToken *bool   `query:"isFacebookToken"`
	PartialName     *string `query:"partialName"`
	FacebookToken   *string `query:"facebookToken"`
	Email           *string `query:"email"`
	Pagination
}

type CreateCustomerBody struct {
	Name           string `json:"name"`
	Email          string `json:"email"`
	Password       string `json:"password"`
	WordpressURL   string `json:"wordpress_url"`
	DeleteHashFlag int    `json:"delete_hash_flag"`
}

type UpdateCustomerBody struct {
	Name           *string    `json:"name"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag *int       `json:"delete_hash_flag"`
}
