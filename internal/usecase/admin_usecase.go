package usecase

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
)

type adminUsecase struct {
	baseRepository  repository.BaseRepository
	adminService    service.AdminService
	authService     service.AuthService
	customerService service.CustomerService
	postService     service.PostService
}

type AdminUsecase interface{}

func NewAdminUsecase(
	baseRepo repository.BaseRepository,
	adminSrv service.AdminService,
	authSrv service.AuthService,
	customerService service.CustomerService,
) AdminUsecase {
	return &adminUsecase{
		baseRepository:  baseRepo,
		adminService:    adminSrv,
		authService:     authSrv,
		customerService: customerService,
	}
}

func (a *adminUsecase) RegisterCustomer(ctx context.Context, customer *model.Customer) (*model.Customer, error) {
	panic("implement me")
}

func (a *adminUsecase) RegisterAdmin(ctx context.Context, admin *model.Admin) error {
	panic("implement me")
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

func (a *adminUsecase) GetCustomers(ctx context.Context) ([]*model.Customer, error) {
	return a.customerService.FindAll(ctx)
}

func (a *adminUsecase) GetCustomer(ctx context.Context, id int) (*model.Customer, error) {
	return a.customerService.FindByID(ctx, id)
}

func (a *adminUsecase) GetAdmins(ctx context.Context) ([]*model.Admin, error) {
	return a.adminService.FindAll(ctx)
}

func (a *adminUsecase) GetAdmin(ctx context.Context, id int) (*model.Admin, error) {
	return a.adminService.FindByID(ctx, id)
}

func (a *adminUsecase) GetPostByCustomer(ctx context.Context, customerId int) ([]*model.Post, error) {
	return a.postService.FindByCustomerID(ctx, customerId)
}
