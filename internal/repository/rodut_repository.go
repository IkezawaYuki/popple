package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/exreq"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/exres"
)

type RodutRepository interface {
	CreatePost(ctx context.Context, wordpressUrl string, email string, post *entity.WordpressPost) (*exres.CreatePostResponse, error)
	UploadMedia(ctx context.Context, wordpressUrl string, filePath string) (*exres.UploadMediaResponse, error)
}

func NewRodutRepository(httpClient infrastructure.HttpClient) RodutRepository {
	return &rodutRepository{
		httpClient: httpClient,
		ApiKey:     config.Env.RodutKey,
	}
}

type rodutRepository struct {
	httpClient infrastructure.HttpClient
	ApiKey     string
}

const createPostEndpoint = "create-post"

func (r *rodutRepository) CreatePost(ctx context.Context, wordpressUrl string, email string, post *entity.WordpressPost) (*exres.CreatePostResponse, error) {
	url := fmt.Sprintf("https://%s/rodut/v1/%s", wordpressUrl, createPostEndpoint)
	responseBody, err := r.httpClient.PostRequest(ctx, url, &exreq.CreatePostRequest{
		ApiKey:        r.ApiKey,
		Email:         email,
		Title:         post.Title,
		Content:       post.Content,
		FeaturedMedia: post.FeaturedMedia,
	}, "")
	if err != nil {
		return nil, err
	}
	createdPost := exres.CreatePostResponse{}
	err = json.Unmarshal(responseBody, &createdPost)
	if err != nil {
		return nil, err
	}
	return &createdPost, nil
}

const uploadMediaEndpoint = "upload-media"

func (r *rodutRepository) UploadMedia(ctx context.Context, wordpressUrl string, filePath string) (*exres.UploadMediaResponse, error) {
	url := fmt.Sprintf("https://%s/rodut/v1/%s", wordpressUrl, uploadMediaEndpoint)
	response, err := r.httpClient.UploadFile(ctx, url, filePath, "")
	if err != nil {
		return nil, err
	}
	uploadMedia := exres.UploadMediaResponse{}
	err = json.Unmarshal(response, &uploadMedia)
	if err != nil {
		return nil, err
	}
	return &uploadMedia, nil
}
