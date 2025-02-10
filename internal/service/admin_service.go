package service

import (
	"context"
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type adminService struct {
	customerRepository repository.CustomerRepository
	adminRepository    repository.AdminRepository
}

type AdminService interface {
	FindByID(ctx context.Context, id int) (*model.Admin, error)
	FindByEmail(ctx context.Context, email string) (*model.Admin, error)
	FindAll(ctx context.Context) ([]*model.Admin, error)
}

func NewAdminService(customerRepo repository.CustomerRepository, adminRepo repository.AdminRepository) AdminService {
	return &adminService{
		customerRepository: customerRepo,
		adminRepository:    adminRepo,
	}
}

func (a *adminService) FindAll(ctx context.Context) ([]*model.Admin, error) {
	return a.adminRepository.Get(ctx, nil)
}

func (a *adminService) FindByEmail(ctx context.Context, email string) (*model.Admin, error) {
	return a.adminRepository.First(ctx, &repository.AdminFilter{Email: &email})
}

func (a *adminService) FindByID(ctx context.Context, id int) (*model.Admin, error) {
	return a.adminRepository.First(ctx, &repository.AdminFilter{ID: &id})
}

func (a *adminService) CreateAdmin(ctx context.Context, admin *model.Admin) error {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	adminModel := model.Admin{
		Name:     admin.Name,
		Email:    admin.Email,
		Password: string(passwordHash),
	}
	if err := a.adminRepository.Save(ctx, &adminModel); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return objects.ErrDuplicateEmail
		}
		return err
	}
	return nil
}
