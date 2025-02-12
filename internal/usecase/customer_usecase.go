package usecase

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/res"
)

type customerUsecase struct {
	baseRepository   repository.BaseRepository
	customerService  service.CustomerService
	authService      service.AuthService
	postService      service.PostService
	wordpressRestApi service.WordpressRestAPI
	graphApi         service.GraphAPI
	fileTransfer     service.FileService
}

type CustomerUsecase interface {
	FetchAndPost(ctx context.Context, customerID int) (*res.Message, error)
	Login(ctx context.Context, user *entity.User) (string, error)
	GetCustomer(ctx context.Context, id int) (*model.Customer, error)
	GetPostsByCustomerID(ctx context.Context, customerID int) ([]*model.Post, error)
}

func NewCustomerUsecase(
	baseRepo repository.BaseRepository,
	customerSrv service.CustomerService,
	authSrv service.AuthService,
	postService service.PostService,
	graphApi service.GraphAPI,
	fileTransfer service.FileService,
) CustomerUsecase {
	return &customerUsecase{
		baseRepository:  baseRepo,
		customerService: customerSrv,
		authService:     authSrv,
		postService:     postService,
		graphApi:        graphApi,
		fileTransfer:    fileTransfer,
	}
}

func (c *customerUsecase) FindAll(ctx context.Context) ([]*model.Customer, error) {
	return c.customerService.FindAll(ctx)
}

func (c *customerUsecase) GetCustomer(ctx context.Context, id int) (*model.Customer, error) {
	return c.customerService.FindByID(ctx, id)
}

func (c *customerUsecase) Login(ctx context.Context, user *entity.User) (string, error) {
	customer, err := c.customerService.FindByEmail(ctx, user.Email)
	if err != nil {
		return "", err
	}
	if err := c.authService.CheckPassword(user, customer.Password); err != nil {
		return "", fmt.Errorf("invalid password: %w", err)
	}
	return c.authService.GenerateJWTCustomer(customer)
}

func (c *customerUsecase) FetchAndPost(ctx context.Context, customerID int) (*res.Message, error) {
	customer, err := c.customerService.FindByID(ctx, customerID)
	if err != nil {
		return nil, objects.ErrNotFound
	}
	if customer.FacebookToken == nil {
		return nil, fmt.Errorf("customer.FacebookToken is nil")
	}

	// インスタグラムの投稿を最新から50件取得する
	instagramPosts, err := c.graphApi.GetInstagramPosts(ctx, *customer.FacebookToken, *customer.InstagramID)
	if err != nil {
		return nil, err
	}
	for _, instagramMedia := range instagramPosts.Media.Data {
		isLinked, err := c.postService.IsLinked(ctx, instagramMedia.ID)
		if err != nil {
			return nil, err
		}
		// 連携済みのものは処理をスキップ
		if isLinked {
			continue
		}

		// 一時フォルダを作り、メディアをダウンロード
		err = c.fileTransfer.MakeTempDirectory(customerID)
		if err != nil {
			return nil, err
		}
		fileList, err := c.fileTransfer.DownloadMediaFiles(ctx, customerID, instagramMedia)
		if err != nil {
			return nil, err
		}

		// ワードプレスにメディアをアップロード
		wordpressMedia, err := c.wordpressRestApi.UploadFiles(ctx, customer.WordpressURL, fileList)
		if err != nil {
			return nil, err
		}

		wordpressResp, err := c.wordpressRestApi.CreatePost(ctx, customer.WordpressURL, instagramMedia, wordpressMedia)
		if err != nil {
			return nil, err
		}

		err = c.postService.Create(ctx, &model.Post{
			CustomerID:       customerID,
			InstagramMediaID: instagramMedia.ID,
			InstagramLink:    instagramMedia.MediaURL,
			WordpressMediaID: wordpressResp.PostId,
			WordpressLink:    wordpressResp.PostUrl,
		})
		if err != nil {
			return nil, err
		}

		err = c.fileTransfer.RemoveTempDirectory(customerID)
		if err != nil {
			return nil, err
		}
	}
	return &res.Message{Message: "ok"}, nil
}

func (c *customerUsecase) GetPostsByCustomerID(ctx context.Context, customerID int) ([]*model.Post, error) {
	customers, err := c.postService.FindByCustomerID(ctx, customerID)
	if err != nil {
		return nil, err
	}

}
