package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type graphAPI struct {
	httpClient *infrastructure.HttpClient
	baseURL    string
}

type GraphAPI interface {
	GetInstagramBusinessAccountID(ctx context.Context, facebookToken string) (string, error)
	GetInstagramPosts(ctx context.Context, facebookToken string, instagramID string) (*entity.InstagramPosts, error)
}

func NewGraph(httpClient *infrastructure.HttpClient) GraphAPI {
	return &graphAPI{
		httpClient: httpClient,
		baseURL:    config.Env.GraphApiURL,
	}
}

const getInstagramBusinessAccountURL = "/me?fields=id,name,accounts{instagram_business_account}"

type GraphApiMeResponse struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Accounts struct {
		Data []struct {
			InstagramBusinessAccount struct {
				ID string `json:"id"`
			} `json:"instagram_business_account"`
			ID string `json:"id"`
		} `json:"data"`
	} `json:"accounts"`
}

func (i *graphAPI) GetInstagramBusinessAccountID(ctx context.Context, facebookToken string) (string, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+getInstagramBusinessAccountURL,
		fmt.Sprintf("Bearer %s", facebookToken))
	if err != nil {
		return "", err
	}
	var instagram GraphApiMeResponse
	err = json.Unmarshal(resp, &instagram)
	if err != nil {
		return "", err
	}
	if len(instagram.Accounts.Data) == 0 {
		return "", objects.ErrNotFound
	}
	return instagram.Accounts.Data[0].InstagramBusinessAccount.ID, nil
}

const getInstagramPosts = "%s?fields=media{id,permalink,caption,timestamp,media_type,media_url,children{media_type,media_url}}"

func (i *graphAPI) GetInstagramPosts(ctx context.Context, facebookToken string, instagramID string) (*entity.InstagramPosts, error) {
	resp, err := i.httpClient.GetRequest(ctx,
		i.baseURL+fmt.Sprintf(getInstagramPosts, instagramID),
		fmt.Sprintf("Bearer %s", facebookToken),
	)
	if err != nil {
		return nil, err
	}
	var posts entity.InstagramPosts
	if err := json.Unmarshal(resp, &posts); err != nil {
		return nil, err
	}
	return &posts, nil
}
