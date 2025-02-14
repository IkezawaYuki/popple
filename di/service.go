package di

import (
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/service"
)

func NewAdminService() service.AdminService {
	return service.NewAdminService(
		NewCustomerRepository(),
		NewAdminRepository(),
	)
}

func NewPostService() service.PostService {
	return service.NewPostService(
		NewPostRepository(),
	)
}

func NewCustomerService() service.CustomerService {
	return service.NewCustomerService(
		NewCustomerRepository(),
	)
}

func NewAuthService() service.AuthService {
	return service.NewAuthService(
		NewCustomerRepository(),
		NewRedisRepository(),
	)
}

func NewGraphAPI() service.GraphAPI {
	return service.NewGraph(
		infrastructure.NewHttpClient(),
	)
}

func NewFileTransfer() service.FileService {
	return service.NewFileService(
		infrastructure.NewHttpClient(),
	)
}

func NewSlackService() service.SlackService {
	return service.NewSlackService(
		infrastructure.NewHttpClient(),
	)
}
