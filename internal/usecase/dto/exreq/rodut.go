package exreq

type CreatePostRequest struct {
	ApiKey        string `json:"api_key"`
	Email         string `json:"email"`
	Title         string `json:"title"`
	Content       string `json:"content"`
	FeaturedMedia int    `json:"featured_media"`
}
