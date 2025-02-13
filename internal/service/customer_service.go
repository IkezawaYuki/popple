package service

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/util"
)

type CustomerService interface {
	FindAuthCustomers(ctx context.Context) ([]*model.Customer, error)
	FindByID(ctx context.Context, id int) (*model.Customer, error)
	FindByEmail(ctx context.Context, email string) (*model.Customer, error)
	IsUsedEmailAddress(ctx context.Context, email string, tx infrastructure.Transaction) (bool, error)
	FindAll(ctx context.Context) ([]*model.Customer, error)
	Create(ctx context.Context, customer *model.Customer) error
}

type customerService struct {
	customerRepository repository.CustomerRepository
}

func NewCustomerService(
	customerRepo repository.CustomerRepository,
) CustomerService {
	return &customerService{
		customerRepository: customerRepo,
	}
}

func (s *customerService) FindAuthCustomers(ctx context.Context) ([]*model.Customer, error) {
	return s.customerRepository.Get(ctx, &repository.CustomerFilter{
		IsFacebookToken: util.Pointer(true),
	})
}

func (s *customerService) FindAll(ctx context.Context) ([]*model.Customer, error) {
	return s.customerRepository.Get(ctx, nil)
}

func (s *customerService) FindByID(ctx context.Context, id int) (*model.Customer, error) {
	return s.customerRepository.First(ctx, &repository.CustomerFilter{
		ID: &id,
	})
}

func (s *customerService) FindByEmail(ctx context.Context, email string) (*model.Customer, error) {
	return s.customerRepository.First(ctx, &repository.CustomerFilter{
		Email: &email,
	})
}

func (s *customerService) IsUsedEmailAddress(ctx context.Context, email string, tx infrastructure.Transaction) (bool, error) {
	customers, err := s.customerRepository.GetTx(ctx, &repository.CustomerFilter{
		Email: &email,
	}, tx)
	if err != nil {
		return false, err
	}
	return len(customers) > 0, nil
}

func (s *customerService) Create(ctx context.Context, customer *model.Customer) error {
	return s.customerRepository.Save(ctx, customer)
}
