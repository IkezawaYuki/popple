package main

import (
	"github.com/IkezawaYuki/popple/di"
	"github.com/IkezawaYuki/popple/docs"
	"github.com/IkezawaYuki/popple/internal/infrastructure"
	"github.com/IkezawaYuki/popple/internal/middleware"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/labstack/echo/v4"
	middleware2 "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"
	"net/http"
)

// @securityDefinitions.apikey BearerAuth
// @in							header
// @name						Authorization
func main() {
	db := infrastructure.GetMysqlConnection()
	conn, err := db.DB()
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = conn.Close()
	}()

	redisCli := infrastructure.GetRedisConnection()
	defer func() {
		_ = redisCli.Close()
	}()

	customerController := di.NewCustomerController(db, redisCli)
	adminController := di.NewAdminController(db, redisCli)
	authService := di.NewAuthService(db, redisCli)
	batchController := di.NewBatchController(db, redisCli)
	pres := presenter.NewPresenter()

	customerAuthMiddleware := middleware.NewCustomerAuthMiddleware(authService, pres)
	adminAuthMiddleware := middleware.NewAdminAuthMiddleware(authService, pres)
	badgeAuthMiddleware := middleware.NewBatchAuthMiddleware(authService, pres)

	e := echo.New()

	e.Use(middleware2.CORSWithConfig(middleware2.CORSConfig{
		AllowOrigins: []string{
			"*",
		},
		AllowHeaders: []string{
			"Authorization",
			"Content-Type",
		},
		AllowMethods: []string{
			echo.GET,
			echo.POST,
			echo.PUT,
			echo.DELETE,
		},
	}))

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	})
	v1 := e.Group("/api/v1")
	{
		v1.POST("/customer/login", customerController.Login)
		v1.POST("/admin/login", adminController.Login)

		customerHandler := v1.Group("/customer")
		customerHandler.Use(customerAuthMiddleware)
		customerHandler.GET("/:id", func(c echo.Context) error {
			return customerController.GetCustomer(c)
		})

		customerHandler.GET("/:id/instagram", func(c echo.Context) error {
			return c.String(http.StatusOK, c.Param("id"))
		})
		customerHandler.POST("/:id/facebook_token", func(c echo.Context) error {
			return c.String(http.StatusOK, c.Param("id"))
		})
		customerHandler.POST("/i/fetch/post", customerController.FetchAndPost)

		adminHandler := v1.Group("/admin")
		adminHandler.Use(adminAuthMiddleware)
		adminHandler.POST("/register/customer", adminController.RegisterCustomer)
		adminHandler.POST("/register/admin", adminController.RegisterAdmin)
		adminHandler.GET("/admins/i", adminController.GetAdmin)
		adminHandler.GET("/admin/customers", adminController.GetCustomers)
		adminHandler.GET("/admin/customers/:customerId", adminController.GetCustomer)
		adminHandler.GET("/admin/customers/:customerId/posts", adminController.GetPostsByCustomer)

		batchHandler := v1.Group("/batch")
		batchHandler.Use(badgeAuthMiddleware)
		batchHandler.GET("/execute", batchController.Execute)
	}

	docs.SwaggerInfo.Title = "Popple API"
	docs.SwaggerInfo.Description = "Popple is very very exciting api!!!"
	docs.SwaggerInfo.Version = "0.1"
	docs.SwaggerInfo.Host = "127.0.0.1:1323"
	docs.SwaggerInfo.BasePath = "/api/v1"
	docs.SwaggerInfo.Schemes = []string{"http", "https"}

	e.GET("/swagger/*", echoSwagger.WrapHandler)

	e.Logger.Fatal(e.Start(":1323"))
}
