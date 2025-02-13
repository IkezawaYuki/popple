package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/req"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/res"
	"golang.org/x/crypto/bcrypt"
)

type adminUsecase struct {
	baseRepository  repository.BaseRepository
	adminRepo       repository.AdminRepository
	customerRepo    repository.CustomerRepository
	postRepo        repository.PostRepository
	adminService    service.AdminService
	authService     service.AuthService
	customerService service.CustomerService
	postService     service.PostService
}

type AdminUsecase interface {
	RegisterCustomer(ctx context.Context, body req.CreateCustomerBody) (resp *res.Customer, err error)
	RegisterAdmin(ctx context.Context, body req.CreateAdminBody) (*res.Admin, error)
	Login(ctx context.Context, user *entity.User) (string, error)
	GetCustomers(ctx context.Context, query req.CustomerQuery) (*res.Customers, error)
	GetCustomer(ctx context.Context, id int) (*res.Customer, error)
	GetAdmin(ctx context.Context, id int) (*res.Admin, error)
	GetAdmins(ctx context.Context, query req.AdminQuery) (*res.Admins, error)
	GetPosts(ctx context.Context, customerID int, query req.PostQuery) (*res.Posts, error)
}

func NewAdminUsecase(
	baseRepo repository.BaseRepository,
	adminRepo repository.AdminRepository,
	customerRepo repository.CustomerRepository,
	postRepo repository.PostRepository,
	adminSrv service.AdminService,
	authSrv service.AuthService,
	customerService service.CustomerService,
) AdminUsecase {
	return &adminUsecase{
		baseRepository:  baseRepo,
		adminRepo:       adminRepo,
		customerRepo:    customerRepo,
		postRepo:        postRepo,
		adminService:    adminSrv,
		authService:     authSrv,
		customerService: customerService,
	}
}

func (a *adminUsecase) RegisterCustomer(ctx context.Context, body req.CreateCustomerBody) (resp *res.Customer, err error) {
	tx := a.baseRepository.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	used, err := a.customerService.IsUsedEmailAddress(ctx, body.Email, tx)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, objects.ErrEmailUsed
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	customer := &model.Customer{
		Name:           body.Name,
		Email:          body.Email,
		Password:       string(passwordHash),
		WordpressURL:   body.WordpressURL,
		DeleteHashFlag: 0,
	}
	err = a.customerRepo.SaveTx(ctx, customer, tx)
	if err != nil {
		return nil, err
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}

	return res.GetCustomer(customer), nil
}

func (a *adminUsecase) RegisterAdmin(ctx context.Context, body req.CreateAdminBody) (*res.Admin, error) {
	used, err := a.adminService.IsUsedEmailAddress(ctx, body.Email)
	if err != nil {
		return nil, err
	}
	if used {
		return nil, objects.ErrEmailUsed
	}
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	admin := &model.Admin{
		Name:     body.Name,
		Email:    body.Email,
		Password: string(passwordHash),
	}
	err = a.adminRepo.Save(ctx, admin)
	if err != nil {
		return nil, err
	}
	return res.GetAdmin(admin), nil
}

func (a *adminUsecase) Login(ctx context.Context, user *entity.User) (string, error) {
	customer, err := a.adminService.FindByEmail(ctx, user.Email)
	if err != nil {
		return "", objects.ErrNotFound
	}
	if err := a.authService.CheckPassword(user, customer.Password); err != nil {
		return "", err
	}
	return a.authService.GenerateJWTAdmin(customer)
}

func (a *adminUsecase) GetCustomers(ctx context.Context, query req.CustomerQuery) (*res.Customers, error) {
	f := &repository.CustomerFilter{
		Email:           query.Email,
		PartialName:     query.PartialName,
		FacebookToken:   query.FacebookToken,
		IsFacebookToken: query.IsFacebookToken,
		Limit:           query.Limit,
		Offset:          query.Offset,
	}
	customers, err := a.customerRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}
	count, err := a.customerRepo.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	return res.GetCustomers(customers, count), nil
}

func (a *adminUsecase) GetCustomer(ctx context.Context, id int) (*res.Customer, error) {
	customer, err := a.customerService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res.GetCustomer(customer), nil
}

func (a *adminUsecase) GetAdmins(ctx context.Context, query req.AdminQuery) (*res.Admins, error) {
	f := &repository.AdminFilter{
		Email:       query.Email,
		PartialName: query.PartialName,
		Limit:       query.Limit,
		Offset:      query.Offset,
	}
	admins, err := a.adminRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}
	count, err := a.adminRepo.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	return res.GetAdmins(admins, count), nil
}

func (a *adminUsecase) GetAdmin(ctx context.Context, id int) (*res.Admin, error) {
	admin, err := a.adminService.FindByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return res.GetAdmin(admin), nil
}

func (a *adminUsecase) GetPosts(ctx context.Context, customerID int, query req.PostQuery) (*res.Posts, error) {
	f := &repository.PostFilter{
		CustomerID: &customerID,
		Limit:      query.Limit,
		Offset:     query.Offset,
	}
	posts, err := a.postRepo.Get(ctx, f)
	if err != nil {
		return nil, err
	}
	count, err := a.postRepo.Count(ctx, f)
	if err != nil {
		return nil, err
	}
	return res.GetPosts(posts, count), nil

}
