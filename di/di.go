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

func NewAuthService(dbDriver infrastructure.DBDriver, redisCli *redis.Client) service.AuthService {
	customerRepo := repository.NewCustomerRepository(dbDriver)
	redisClient := repository.NewRedisRepository(redisCli)
	return service.NewAuthService(customerRepo, redisClient)
}

func NewCustomerService(dbDriver infrastructure.DBDriver) service.CustomerService {
	customerRepo := repository.NewCustomerRepository(dbDriver)
	return service.NewCustomerService(customerRepo)
}

func NewCustomerController(db *gorm.DB, redisCli *redis.Client) controller.CustomerController {
	dbDriver := infrastructure.NewDBDriver(db)
	pre := presenter.NewPresenter()
	customerUsecase := NewCustomerUsecase(dbDriver, redisCli)
	return controller.NewCustomerController(customerUsecase, pre)
}

func NewAdminController(db *gorm.DB, redisCli *redis.Client) controller.AdminController {
	dbDriver := infrastructure.NewDBDriver(db)
	baseRepo := repository.NewBaseRepository(dbDriver)
	customerRepo := repository.NewCustomerRepository(dbDriver)
	adminRepo := repository.NewAdminRepository(dbDriver)
	redisClient := repository.NewRedisRepository(redisCli)
	pre := presenter.NewPresenter()
	customerService := service.NewCustomerService(customerRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	adminService := service.NewAdminService(customerRepo, adminRepo)
	adminUsecase := usecase.NewAdminUsecase(baseRepo, adminService, authService, customerService)
	return controller.NewAdminController(adminUsecase, pre)
}

func NewBatchController(dbDriver infrastructure.DBDriver, redisCli *redis.Client) controller.BatchController {
	pre := presenter.NewPresenter()
	httpClient := infrastructure.NewHttpClient()
	slack := service.NewSlackService(httpClient)
	customerUsecase := NewCustomerUsecase(dbDriver, redisCli)
	batchUsecase := usecase.NewBatchUsecase(customerUsecase, slack)
	return controller.NewBatchController(batchUsecase, pre)
}

func NewCustomerUsecase(dbDriver infrastructure.DBDriver, redisCli *redis.Client) usecase.CustomerUsecase {
	httpClient := infrastructure.NewHttpClient()
	baseRepo := repository.NewBaseRepository(dbDriver)
	customerRepo := repository.NewCustomerRepository(dbDriver)
	postRepo := repository.NewPostRepository(dbDriver)
	redisClient := repository.NewRedisRepository(redisCli)
	customerService := service.NewCustomerService(customerRepo)
	authService := service.NewAuthService(customerRepo, redisClient)
	postService := service.NewPostService(postRepo)
	graphApi := service.NewGraph(httpClient)
	fileTransfer := service.NewFileService(httpClient)
	return usecase.NewCustomerUsecase(
		baseRepo,
		customerService,
		authService,
		postService,
		graphApi,
		fileTransfer)
}
