package service

import (
	"context"
	"fmt"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/ory/dockertest/v3"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"testing"
	"time"
)

var (
	adminSrv     *adminService
	authSrv      *AuthService
	customerSrv  *CustomerService
	postSrv      *postService
	wordpressSrv *WordpressRestAPI
	graphSrv     *GraphAPI
	fileTransfer *FileService
	slackSrv     *SlackService
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	if err := pool.Client.Ping(); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}
	resourceMysql, err := pool.Run("mysql", "latest", []string{"MYSQL_ROOT_PASSWORD=secret"})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	resourceRedis, err := pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start Redis resource: %s", err)
	}

	dsn := fmt.Sprintf("root:secret@(localhost:%s)/mysql?parseTime=true", resourceMysql.GetPort("3306/tcp"))
	var dbInit *gorm.DB
	var client *redis.Client

	if err := pool.Retry(func() error {
		var err error
		dbInit, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
		if err != nil {
			return err
		}
		sqlDB, err := dbInit.DB()
		if err != nil {
			return err
		}
		err = sqlDB.Ping()
		if err != nil {
			return err
		}

		err = dbInit.Exec("CREATE DATABASE IF NOT EXISTS c_root;").Error
		if err != nil {
			return err
		}
		return nil
	}); err != nil {
		_ = pool.Purge(resourceMysql)
		_ = pool.Purge(resourceRedis)
	}

	// Redisへの接続をリトライ
	err = pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr: fmt.Sprintf("localhost:%s", resourceRedis.GetPort("6379/tcp")),
		})

		// Pingを使ってRedisへの接続を確認
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		_, err := client.Ping(ctx).Result()
		return err
	})
	if err != nil {
		_ = pool.Purge(resourceMysql)
		_ = pool.Purge(resourceRedis)
	}

	var db *gorm.DB
	db, err = gorm.Open(mysql.Open(fmt.Sprintf("root:secret@(localhost:%s)/c_root?parseTime=true", resourceMysql.GetPort("3306/tcp"))), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	if err := db.AutoMigrate(&model.Customer{}); err != nil {
		log.Fatalf("AutoMigrate err: %v\n", err)
	}
	if err := db.AutoMigrate(&model.Post{}); err != nil {
		log.Fatalf("AutoMigrate err: %v\n", err)
	}
	if err := db.AutoMigrate(&model.Admin{}); err != nil {
		log.Fatalf("AutoMigrate err: %v\n", err)
	}

	customerRepo := repository.NewCustomerRepository(db)
	adminRepo := repository.NewAdminRepository(db)
	postRepo := repository.NewPostRepository(db)
	redisRepo := repository.NewRedisClient(client)
	httpClient := infrastructure.NewHttpClient()

	adminSrv = NewAdminService(customerRepo, adminRepo)
	authSrv = NewAuthService(customerRepo, redisRepo)
	customerSrv = NewCustomerService(customerRepo, postRepo)
	postSrv = NewPostService(postRepo)
	wordpressSrv = NewWordpressRestAPI(httpClient)
	graphSrv = NewGraph(httpClient)
	fileTransfer = NewFileService(httpClient)
	slackSrv = NewSlackService(httpClient)

	admins, err := adminSrv.FindAll(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(admins)

	m.Run()
}
