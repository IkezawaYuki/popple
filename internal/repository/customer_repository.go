package repository

import (
	"context"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
)

type CustomerRepository interface {
	Get(ctx context.Context, f infrastructure.Filter) ([]*model.Customer, error)
	GetTx(ctx context.Context, f infrastructure.Filter, tx infrastructure.Transaction) ([]*model.Customer, error)
	First(ctx context.Context, f infrastructure.Filter) (*model.Customer, error)
	Save(ctx context.Context, customer *model.Customer) error
	SaveTx(ctx context.Context, customer *model.Customer, tx infrastructure.Transaction) error
	Delete(ctx context.Context, f infrastructure.Filter) error
	DeleteTx(ctx context.Context, f infrastructure.Filter, tx infrastructure.Transaction) error
	Count(ctx context.Context, f infrastructure.Filter) (int, error)
}

func NewCustomerRepository(dbDriver infrastructure.DBDriver) CustomerRepository {
	return &customerRepository{
		dbDriver: dbDriver,
	}
}

type customerRepository struct {
	dbDriver infrastructure.DBDriver
}

func (c *customerRepository) Get(ctx context.Context, f infrastructure.Filter) ([]*model.Customer, error) {
	var customers []*model.Customer
	err := c.dbDriver.Find(ctx, &customers, f)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *customerRepository) GetTx(ctx context.Context, f infrastructure.Filter, tx infrastructure.Transaction) ([]*model.Customer, error) {
	var customers []*model.Customer
	err := c.dbDriver.FindTx(ctx, &customers, f, tx)
	if err != nil {
		return nil, err
	}
	return customers, nil
}

func (c *customerRepository) First(ctx context.Context, f infrastructure.Filter) (*model.Customer, error) {
	var customer model.Customer
	err := c.dbDriver.First(ctx, &customer, f)
	if err != nil {
		return nil, err
	}
	return &customer, nil
}

func (c *customerRepository) Save(ctx context.Context, customer *model.Customer) error {
	return c.dbDriver.Save(ctx, customer)
}

func (c *customerRepository) SaveTx(ctx context.Context, customer *model.Customer, tx infrastructure.Transaction) error {
	return c.dbDriver.SaveTx(ctx, customer, tx)
}

func (c *customerRepository) Delete(ctx context.Context, f infrastructure.Filter) error {
	return c.dbDriver.Delete(ctx, &model.Customer{}, f)
}

func (c *customerRepository) DeleteTx(ctx context.Context, f infrastructure.Filter, tx infrastructure.Transaction) error {
	return c.dbDriver.DeleteTx(ctx, &model.Customer{}, f, tx)
}

func (c *customerRepository) Count(ctx context.Context, f infrastructure.Filter) (int, error) {
	count, err := c.dbDriver.Count(ctx, &model.Customer{}, f)
	if err != nil {
		return 0, err
	}
	return int(count), nil
}
