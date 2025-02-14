package middleware

import (
	"fmt"
	"github.com/IkezawaYuki/popple/internal/presenter"
	"github.com/IkezawaYuki/popple/internal/service"
	"github.com/labstack/echo/v4"
	"log/slog"
)

func NewBatchAuthMiddleware(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			return next(c)
		}
	}
}

func NewAdminAuthMiddleware(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("AdminAuthMiddleware is invoked")
			token := c.Request().Header.Get("Authorization")
			adminUuid, err := authService.IsAdminLogin(token)
			if err != nil {
				return c.JSON(presenter.Generate(err, nil))
			}
			slog.Info("AdminAuthMiddleware is OK")
			c.Set("admin_id", adminUuid)
			return next(c)
		}
	}
}

func NewCustomerAuthMiddleware(authService service.AuthService) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			slog.Info("CustomerAuthMiddleware is invoked")
			fmt.Println(c.Request().Header)
			token := c.Request().Header.Get("Authorization")
			customerID, err := authService.IsCustomerIsLogin(token)
			if err != nil {
				return c.JSON(presenter.Generate(err, nil))
			}
			slog.Info("CustomerAuthMiddleware is OK")
			c.Set("customer_id", customerID)
			return next(c)
		}
	}
}
