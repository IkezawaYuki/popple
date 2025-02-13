package usecase

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/req"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/res"
)

type CustomerUsecase interface {
	FetchAndPost(ctx context.Context, customerID int) (*res.Message, error)
	Login(ctx context.Context, user *entity.User) (string, error)
	GetCustomer(ctx context.Context, customerID int) (*res.Customer, error)
	GetPosts(ctx context.Context, customerID int, req req.PostQuery) (*res.Posts, error)
}

type customerUsecase struct {
	baseRepo        repository.BaseRepository
	postRepo        repository.PostRepository
	customerRepo    repository.CustomerRepository
	rodutRepo       repository.RodutRepository
	customerService service.CustomerService
	authService     service.AuthService
	postService     service.PostService
	graphApi        service.GraphAPI
	fileTransfer    service.FileService
}

func NewCustomerUsecase(
	baseRepo repository.BaseRepository,
	postRepo repository.PostRepository,
	customerRepo repository.CustomerRepository,
	customerSrv service.CustomerService,
	authSrv service.AuthService,
	postService service.PostService,
	graphApi service.GraphAPI,
	fileTransfer service.FileService,
	rodutRepo repository.RodutRepository,
) CustomerUsecase {
	return &customerUsecase{
		baseRepo:        baseRepo,
		postRepo:        postRepo,
		customerRepo:    customerRepo,
		customerService: customerSrv,
		authService:     authSrv,
		postService:     postService,
		graphApi:        graphApi,
		fileTransfer:    fileTransfer,
		rodutRepo:       rodutRepo,
	}
}

func (c *customerUsecase) FindAll(ctx context.Context, query req.CustomerQuery) (*res.Customers, error) {
	f := &repository.CustomerFilter{
		PartialName:     query.PartialName,
		IsFacebookToken: query.IsFacebookToken,
		Limit:           query.Limit,
		Offset:          query.Offset,
	}
	customers, err := c.customerRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}
	count, err := c.customerRepo.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	return res.GetCustomers(customers, count), nil
}

func (c *customerUsecase) GetCustomer(ctx context.Context, id int) (*res.Customer, error) {
	customer, err := c.customerRepo.First(ctx, &repository.CustomerFilter{ID: &id})
	if err != nil {
		return nil, err
	}
	return res.GetCustomer(customer), nil
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
		wordpressMedia, err := c.rodutRepo.UploadMedias(ctx, customer.WordpressURL, fileList)
		if err != nil {
			return nil, err
		}

		createPost := entity.NewWordpressPost(instagramMedia, wordpressMedia)
		wordpressResp, err := c.rodutRepo.CreatePost(ctx, customer.WordpressURL, createPost)
		if err != nil {
			return nil, err
		}

		err = c.postRepo.Save(ctx, &model.Post{
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

func (c *customerUsecase) GetPosts(ctx context.Context, customerID int, req req.PostQuery) (*res.Posts, error) {
	f := &repository.PostFilter{
		CustomerID: &customerID,
		Limit:      req.Limit,
		Offset:     req.Offset,
	}
	posts, err := c.postRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}
	counts, err := c.postRepo.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	return res.GetPosts(posts, counts), nil
}
