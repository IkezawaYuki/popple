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
	"golang.org/x/crypto/bcrypt"
)

type adminUsecase struct {
	baseRepository  repository.BaseRepository
	adminRepo       repository.AdminRepository
	customerRepo    repository.CustomerRepository
	adminService    service.AdminService
	authService     service.AuthService
	customerService service.CustomerService
	postService     service.PostService
}

type AdminUsecase interface {
}

func NewAdminUsecase(
	baseRepo repository.BaseRepository,
	adminRepo repository.AdminRepository,
	customerRepo repository.CustomerRepository,
	adminSrv service.AdminService,
	authSrv service.AuthService,
	customerService service.CustomerService,
) AdminUsecase {
	return &adminUsecase{
		baseRepository:  baseRepo,
		adminRepo:       adminRepo,
		customerRepo:    customerRepo,
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
