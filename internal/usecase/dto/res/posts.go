package res

import "time"

type Posts struct {
	Posts []*Post `json:"posts"`
}

type Post struct {
	ID               int       `json:"id"`
	CustomerID       int       `json:"customer_id"`
	InstagramMediaID string    `json:"instagram_media_id"`
	InstagramLink    string    `json:"instagram_link"`
	WordpressMediaID string    `json:"wordpress_media_id"`
	WordpressLink    string    `json:"wordpress_link"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
	Paginate
}
