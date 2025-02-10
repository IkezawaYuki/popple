package exres

type CreatePostResponse struct {
	PostId  string `json:"post_id"`
	PostUrl string `json:"post_url"`
}

type UploadMediaResponse struct {
	ID        string `json:"id"`
	SourceUrl string `json:"source_url"`
}
