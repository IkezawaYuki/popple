package presenter

import (
	"errors"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"log/slog"
	"net/http"
)

func Generate(err error, body any) (int, any) {
	slog.Info("Generate is invoked")
	if err == nil {
		return http.StatusOK, body
	}
	switch {
	case errors.Is(err, objects.ErrNotFound):
		return http.StatusNotFound, err.Error()
	case errors.Is(err, objects.ErrAuthentication):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, objects.ErrAuthorization):
		return http.StatusUnauthorized, err.Error()
	case errors.Is(err, objects.ErrDuplicateEmail):
		return http.StatusBadRequest, err.Error()
	case errors.Is(err, objects.ErrEmailUsed):
		return http.StatusBadRequest, err.Error()
	default:
		return http.StatusInternalServerError, err.Error()
	}
}
