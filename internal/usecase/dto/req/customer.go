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
	Name           string `json:"name" example:"yuki"`
	Email          string `json:"email" example:"yuki@gmail.com"`
	Password       string `json:"password" example:"123456"`
	WordpressURL   string `json:"wordpress_url" example:"example.com"`
	DeleteHashFlag int    `json:"delete_hash_flag" example:"0"`
}

type UpdateCustomerBody struct {
	Name           *string    `json:"name"`
	FacebookToken  *string    `json:"facebook_token"`
	StartDate      *time.Time `json:"start_date"`
	InstagramID    *string    `json:"instagram_id"`
	InstagramName  *string    `json:"instagram_name"`
	DeleteHashFlag *int       `json:"delete_hash_flag"`
}
