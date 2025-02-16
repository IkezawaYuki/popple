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
}

func NewCustomerController(customerUsecase usecase.CustomerUsecase) CustomerController {
	return CustomerController{
		customerUsecase: customerUsecase,
	}
}

// Login godoc
//
// @Summary	ログイン
// @Description	顧客としてログインします
// @Tags Customer
// @Accept application/json
// @Param  body body req.User  true  "ユーザー情報"
// @Success      200   {object}  res.Auth
// @Router /customer/login [post]
func (ctr *CustomerController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user entity.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, "invalid value")
	}
	token, err := ctr.customerUsecase.Login(c.Request().Context(), &user)
	return c.JSON(presenter.Generate(err, token))
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
	customer, err := ctr.customerUsecase.GetCustomer(c.Request().Context(), customerId)
	return c.JSON(presenter.Generate(err, customer))
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
	return c.JSON(presenter.Generate(err, posts))
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
	resp, err := ctr.customerUsecase.FetchAndPost(c.Request().Context(), customerID)
	return c.JSON(presenter.Generate(err, resp))
}
