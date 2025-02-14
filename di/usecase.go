package di

import (
	"github.com/IkezawaYuki/popple/internal/usecase"
)

func NewAdminUsecase() usecase.AdminUsecase {
	return usecase.NewAdminUsecase(
		NewBaseRepository(),
		NewAdminRepository(),
		NewCustomerRepository(),
		NewPostRepository(),
		NewAdminService(),
		NewAuthService(),
		NewCustomerService(),
	)
}

func NewCustomerUsecase() usecase.CustomerUsecase {
	return usecase.NewCustomerUsecase(
		NewBaseRepository(),
		NewPostRepository(),
		NewCustomerRepository(),
		NewCustomerService(),
		NewAuthService(),
		NewPostService(),
		NewGraphAPI(),
		NewFileTransfer(),
		NewRodutRepository(),
	)
}

func NewBatchUsecase() usecase.BatchUsecase {
	return usecase.NewBatchUsecase(
		NewCustomerService(),
		NewCustomerUsecase(),
		NewSlackService(),
	)
}
