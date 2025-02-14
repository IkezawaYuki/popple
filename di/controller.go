package di

import (
	"github.com/IkezawaYuki/popple/internal/controller"
)

func NewCustomerController() controller.CustomerController {
	return controller.NewCustomerController(NewCustomerUsecase())
}

func NewAdminController() controller.AdminController {
	return controller.NewAdminController(NewAdminUsecase())
}

func NewBatchController() controller.BatchController {
	return controller.NewBatchController(NewBatchUsecase())
}
