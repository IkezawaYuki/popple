package service

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/exres"
)

type WordpressRestAPI struct {
	httpClient *infrastructure.HttpClient
}

func NewWordpressRestAPI(httpClient *infrastructure.HttpClient) *WordpressRestAPI {
	return &WordpressRestAPI{
		httpClient: httpClient,
	}
}

const wpPostUrl = "https://%s/wp-json/wp/v2/posts"

func (w *WordpressRestAPI) CreatePost(ctx context.Context, wordpressURL string, instaDetail entity.InstagramPost, wpMedia []*entity.WordpressMedia) (*exres.CreatePostResponse, error) {
	posts := entity.NewWordpressPost(instaDetail, wpMedia)
	resp, err := w.httpClient.PostRequest(ctx, fmt.Sprintf(wpPostUrl, wordpressURL), posts, BasicAuthHeader())
	if err != nil {
		return nil, err
	}
	var response exres.CreatePostResponse
	if err := json.Unmarshal(resp, &response); err != nil {
		return nil, err
	}
	return &response, nil
}

func (w *WordpressRestAPI) UploadFiles(ctx context.Context, wordpressURL string, pathList []string) ([]*entity.WordpressMedia, error) {
	var wordpressMediaList []*entity.WordpressMedia
	for _, path := range pathList {
		wordpressMedia, err := w.UploadFile(ctx, wordpressURL, path)
		if err != nil {
			return nil, err
		}
		wordpressMediaList = append(wordpressMediaList, wordpressMedia)
	}
	return wordpressMediaList, nil
}

const wpMediaUrl = "https://%s/wp-json/wp/v2/media"

func (w *WordpressRestAPI) UploadFile(ctx context.Context, wordpressURL string, path string) (*entity.WordpressMedia, error) {
	resp, err := w.httpClient.UploadFile(ctx, fmt.Sprintf(wpMediaUrl, wordpressURL), path, BasicAuthHeader())
	if err != nil {
		return nil, err
	}
	var wordpressMedia entity.WordpressMedia
	if err := json.Unmarshal(resp, &wordpressMedia); err != nil {
		return nil, err
	}
	return &wordpressMedia, nil
}

func BasicAuthHeader() string {
	auth := config.Env.WordpressAdminEmail + ":" + config.Env.WordpressAdminPassword
	encoded := base64.StdEncoding.EncodeToString([]byte(auth))
	return "Basic " + encoded
}
