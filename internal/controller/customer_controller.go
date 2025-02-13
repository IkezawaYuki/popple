package controller

import (
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/req"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
)

type CustomerController struct {
	customerUsecase usecase.CustomerUsecase
	presenter       *presenter.Presenter
}

func NewCustomerController(customerUsecase usecase.CustomerUsecase, presenter2 *presenter.Presenter) CustomerController {
	return CustomerController{
		customerUsecase: customerUsecase,
		presenter:       presenter2,
	}
}

// Login godoc
//
//	@Summary		ログイン
//	@Description	顧客としてログインします
//	@Tags			Customer
//	@Accept			application/x-www-form-urlencoded
//	@Param			email			formData	string	false	"Email"
//	@Param			password		formData	string	false	"Password"
//	@Router			/customer/login [post]
func (ctr *CustomerController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user entity.User
	user.Email = c.FormValue("email")
	user.Password = c.FormValue("password")
	if user.Email == "" || user.Password == "" {
		return c.String(http.StatusBadRequest, "invalid value")
	}
	ctx := c.Request().Context()
	token, err := ctr.customerUsecase.Login(ctx, &user)
	return c.JSON(ctr.presenter.Generate(err, token))
}

// GetCustomer godoc
//
//	@Summary		顧客情報の取得
//	@Description	自分の顧客情報を取得します
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Token
//	@Router			/customer/i [get]
func (ctr *CustomerController) GetCustomer(c echo.Context) error {
	slog.Info("GetCustomer is invoked")
	customerId := c.Get("customer_id").(int)
	ctx := c.Request().Context()
	customer, err := ctr.customerUsecase.GetCustomer(ctx, customerId)
	return c.JSON(ctr.presenter.Generate(err, customer))
}

// GetPosts godoc
//
//	@Summary		顧客の投稿一覧の取得
//	@Description	顧客ごとの投稿を一覧で取得します
//	@Tags			Customer
//	@Accept			json
//	@Produce		json
//	@Security		Token
//	@Router			/customer/i/posts [get]
func (ctr *CustomerController) GetPosts(c echo.Context) error {
	slog.Info("GetPosts is invoked")
	customerId := c.Get("customer_id").(int)
	var query req.PostQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	posts, err := ctr.customerUsecase.GetPosts(c.Request().Context(), customerId, query)
	return c.JSON(ctr.presenter.Generate(err, posts))
}

// FetchAndPost godoc
//
//	@Summary		インスタグラムとWordpressの連携
//	@Description	インスタグラムとWordpressの連携
//	@Tags			Customer
//	@Produce		json
//	@Security		Token
//	@Param			customer_id	path	int	true	"Customer ID"
//	@Router			/customer/i/fetch/post [post]
func (ctr *CustomerController) FetchAndPost(c echo.Context) error {
	slog.Info("FetchAndPost is invoked")
	customerID := c.Get("customer_id").(int)
	ctx := c.Request().Context()
	resp, err := ctr.customerUsecase.FetchAndPost(ctx, customerID)
	return c.JSON(ctr.presenter.Generate(err, resp))
}
