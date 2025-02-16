package controller

import (
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/usecase"
	"github.com/IkezawaYuki/popple/internal/usecase/dto/req"
	"github.com/labstack/echo/v4"
	"log/slog"
	"net/http"
	"strconv"
)

type AdminController struct {
	adminUsecase usecase.AdminUsecase
}

func NewAdminController(adminUsecase usecase.AdminUsecase) AdminController {
	return AdminController{
		adminUsecase: adminUsecase,
	}
}

// RegisterCustomer godoc
//
// @Summary      顧客の作成
// @Description  顧客を作成します
// @Tags         Admin
// @Accept application/json
// @Produce      application/json
// @Security     Token
// @Param        body body req.CreateCustomerBody  true  "顧客作成用リクエスト"
// @Success      200   {object}  res.Customer
// @Router       /admin/register/customer [post]
func (a *AdminController) RegisterCustomer(c echo.Context) error {
	slog.Info("RegisterCustomer is invoked")
	var registerCustomer req.CreateCustomerBody
	if err := c.Bind(&registerCustomer); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"error": err.Error()})
	}

	resp, err := a.adminUsecase.RegisterCustomer(c.Request().Context(), registerCustomer)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": err.Error()})
	}
	return c.JSON(presenter.Generate(err, resp))
}

// Login godoc
//
// @Summary      ログイン
// @Description  管理者としてログインします
// @Tags         Admin
// @Accept       application/json
// @Produce      application/json
// @Param        body body req.User  true  "ユーザー情報"
// @Success      200   {object}  res.Auth "トークン"
// @Router		/admin/login [post]
func (a *AdminController) Login(c echo.Context) error {
	slog.Info("Login is invoked")
	var user entity.User
	if err := c.Bind(&user); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	token, err := a.adminUsecase.Login(c.Request().Context(), &user)
	return c.JSON(presenter.Generate(err, token))
}

// GetCustomers godoc
// @Summary	顧客一覧取得
// @Tags Admin
// @Description	全顧客を一覧で取得します
// @Produce	application/json
// @Accept  application/json
// @Security Token
// @Param query query req.CustomerQuery  true  "ユーザー情報"
// @Router /admin/customers [get]
func (a *AdminController) GetCustomers(c echo.Context) error {
	slog.Info("GetCustomers is invoked")
	var query req.CustomerQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	customers, err := a.adminUsecase.GetCustomers(c.Request().Context(), query)
	return c.JSON(presenter.Generate(err, customers))
}

// GetCustomer godoc
//
//	@Summary		顧客情報取得
//	@Tags			Admin
//	@Description	顧客情報を取得します
//	@Produce		application/json
//	@Security		Token
//	@Param customerId path int  true  "顧客ID"
//	@Router			/admin/customers/{customerId} [get]
func (a *AdminController) GetCustomer(c echo.Context) error {
	slog.Info("GetCustomer is invoked")
	customerIdParam := c.Param("customerId")
	customerId, err := strconv.Atoi(customerIdParam)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	customer, err := a.adminUsecase.GetCustomer(c.Request().Context(), customerId)
	return c.JSON(presenter.Generate(err, customer))
}

// GetPostsByCustomer godoc
//
//	@Summary		投稿取得
//	@Tags			Admin
//	@Description	顧客ごとの投稿データを一覧で取得します
//	@Produce		application/json
//	@Security		Token
//	@Param customerId path int  true  "顧客ID"
//	@Router			/admin/customers/{customerId}/posts [get]
func (a *AdminController) GetPostsByCustomer(c echo.Context) error {
	slog.Info("GetPostsByCustomer is invoked")
	customerIdParam := c.Param("customer_id")
	customerId, err := strconv.Atoi(customerIdParam)
	if err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	var query req.PostQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}

	posts, err := a.adminUsecase.GetPosts(c.Request().Context(), customerId, query)
	return c.JSON(presenter.Generate(err, posts))
}

// GetAdmins godoc
//
// @Summary		管理者ユーザー一覧取得
// @Tags		Admin
// @Description	管理者ユーザーを一覧で取得します
// @Param query query req.AdminQuery  true  "管理者情報"
// @Produce		application/json
// @Security	Token
// @Router		/admin/admins [get]
func (a *AdminController) GetAdmins(c echo.Context) error {
	slog.Info("GetAdmins is invoked")
	var query req.AdminQuery
	if err := c.Bind(&query); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	admins, err := a.adminUsecase.GetAdmins(c.Request().Context(), query)
	return c.JSON(presenter.Generate(err, admins))
}

// GetAdmin godoc
//
//	@Summary		ログイン中の管理者情報取得
//	@Description	ログイン中の管理者情報取得します
//	@Tags			Admin
//	@Produce		application/json
//	@Security		Token
//	@Param			adminId		query	int	true	"Admin ID"
//	@Router			/admin/admins/i [get]
func (a *AdminController) GetAdmin(c echo.Context) error {
	slog.Info("GetAdmin is invoked")
	adminId := c.Get("admin_id").(int)
	admin, err := a.adminUsecase.GetAdmin(c.Request().Context(), adminId)
	return c.JSON(presenter.Generate(err, admin))
}

// RegisterAdmin godoc
//
// @Summary	管理者ユーザーの作成
// @Description	管理者ユーザーを作成します
// @Tags Admin
// @Produce	application/json
// @Accept application/json
// @Security Token
// @Param  body body req.CreateAdminBody  true  "管理者情報"
// @Success      200   {object}  res.Admin
// @Router	/admin/register/admin [post]
func (a *AdminController) RegisterAdmin(c echo.Context) error {
	slog.Info("RegisterAdmin is invoked")
	var admin req.CreateAdminBody
	if err := c.Bind(&admin); err != nil {
		return c.String(http.StatusBadRequest, err.Error())
	}
	resp, err := a.adminUsecase.RegisterAdmin(c.Request().Context(), admin)
	return c.JSON(presenter.Generate(err, resp))
}
