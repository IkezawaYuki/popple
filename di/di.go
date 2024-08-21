package di

import (
	"github.com/IkezawaYuki/popple/internal/controller"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func NewAuthService(db *gorm.DB, redisCli *redis.Client) *service.AuthService {
	customerRepo := repository.NewCustomerRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	return service.NewAuthService(customerRepo, redisClient)
}

func NewCustomerService(db *gorm.DB) *service.CustomerService {
	customerRepo := repository.NewCustomerRepository(db)
	instaRepo := repository.NewInstagramRepository(db)
	instaWordpressRepo := repository.NewInstagramWordpressRepository(db)
	return service.NewCustomerService(customerRepo, instaRepo, instaWordpressRepo)
}

func NewCustomerController(db *gorm.DB, redisCli *redis.Client) controller.CustomerController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	instaRepo := repository.NewInstagramRepository(db)
	instagramWordpressRepo := repository.NewInstagramWordpressRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(customerRepo, instaRepo, instagramWordpressRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	httpClient := infrastructure.NewHttpClient()
	wordpressRestApi := service.NewWordpressRestAPI(httpClient)
	graphApi := service.NewGraph(httpClient)
	fileTransfer := service.NewFileService(httpClient)
	customerUsecase := usecase.NewCustomerUsecase(baseRepo, customerService, authService, wordpressRestApi, graphApi, fileTransfer)
	return controller.NewCustomerController(customerUsecase, pre)
}

func NewAdminController(db *gorm.DB, redisCli *redis.Client) controller.AdminController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	instaRepo := repository.NewInstagramRepository(db)
	instagramWordpressRepo := repository.NewInstagramWordpressRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(customerRepo, instaRepo, instagramWordpressRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	adminService := service.NewAdminService(customerRepo, adminRepo)
	adminUsecase := usecase.NewAdminUsecase(baseRepo, adminService, authService, customerService)
	return controller.NewAdminController(adminUsecase, pre)
}

func NewBatchController(db *gorm.DB, redisCli *redis.Client) controller.BatchController {
	baseRepo := repository.NewBaseRepository(db)
	customerRepo := repository.NewCustomerRepository(db)
	instaRepo := repository.NewInstagramRepository(db)
	instagramWordpressRepo := repository.NewInstagramWordpressRepository(db)
	redisClient := repository.NewRedisClient(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(customerRepo, instaRepo, instagramWordpressRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	httpClient := infrastructure.NewHttpClient()
	wordpressRestApi := service.NewWordpressRestAPI(httpClient)
	graphApi := service.NewGraph(httpClient)
	fileTransfer := service.NewFileService(httpClient)
	customerUsecase := usecase.NewCustomerUsecase(baseRepo, customerService, authService, wordpressRestApi, graphApi, fileTransfer)
	slack := service.NewSlackService(httpClient)
	batchUsecase := usecase.NewBatchUsecase(customerUsecase, slack)
	return controller.NewBatchController(batchUsecase, pre)
}
