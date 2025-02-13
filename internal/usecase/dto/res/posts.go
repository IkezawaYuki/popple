package res

import (
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"time"
)

type Posts struct {
	Posts []*Post `json:"posts"`
	Paginate
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
}

func GetPosts(posts []*model.Post, count int) *Posts {
	resp := make([]*Post, len(posts))
	for i, post := range posts {
		resp[i] = GetPost(post)
	}
	return &Posts{
		Posts: resp,
		Paginate: Paginate{
			Count: count,
			Total: len(posts),
		},
	}
}

func GetPost(post *model.Post) *Post {
	return &Post{
		ID:               post.ID,
		CustomerID:       post.CustomerID,
		InstagramMediaID: post.InstagramMediaID,
		InstagramLink:    post.InstagramLink,
		WordpressMediaID: post.WordpressMediaID,
		WordpressLink:    post.WordpressLink,
		CreatedAt:        post.CreatedAt,
		UpdatedAt:        post.UpdatedAt,
	}
}
