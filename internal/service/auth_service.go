package service

import (
	"fmt"
	"github.com/IkezawaYuki/popple/config"
	"github.com/IkezawaYuki/popple/internal/domain/entity"
	"github.com/IkezawaYuki/popple/internal/domain/model"
	"github.com/IkezawaYuki/popple/internal/domain/objects"
	"github.com/IkezawaYuki/popple/internal/repository"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
	"strings"
	"time"
)

type authService struct {
	customerRepository *repository.CustomerRepository
	redisClient        *repository.RedisClient
}

type AuthService interface {
	CheckPassword(user *entity.User, password string) error
	GenerateJWTAdmin(admin *model.Admin) (string, error)
	GenerateJWTCustomer(c *model.Customer) (string, error)
}

func NewAuthService(customerRepo *repository.CustomerRepository, redisClient *repository.RedisClient) AuthService {
	return &authService{
		customerRepository: customerRepo,
		redisClient:        redisClient,
	}
}

func (a *authService) IsCustomerIsLogin(tokenString string) (int, error) {
	slog.Info("IsCustomerIsLogin is invoked")
	slog.Info(tokenString)
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		slog.Info(err.Error())
		return 0, objects.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, objects.ErrAuthorization
	}
	if !claims.VerifyAudience("customer", true) {
		return 0, objects.ErrAuthentication
	}

	return int(claims["sub"].(float64)), nil
}

func (a *authService) IsAdminLogin(tokenString string) (int, error) {
	slog.Info("IsAdminLogin is invoked")
	tokenString = strings.Replace(tokenString, "Bearer ", "", 1)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(config.Env.AccessSecretKey), nil
	})
	if err != nil {
		slog.Info(err.Error())
		return 0, objects.ErrAuthorization
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, objects.ErrAuthorization
	}
	if !claims.VerifyAudience("admin", true) {
		return 0, objects.ErrAuthentication
	}
	return int(claims["sub"].(float64)), nil
}

func (a *authService) GenerateJWTCustomer(c *model.Customer) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "customer",
		"sub":   c.ID,
		"email": c.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *authService) GenerateJWTAdmin(admin *model.Admin) (string, error) {
	claims := jwt.MapClaims{
		"iss":   "popple",
		"aud":   "admin",
		"sub":   admin.ID,
		"email": admin.Email,
		"exp":   time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(config.Env.AccessSecretKey))
}

func (a *authService) CheckPassword(user *entity.User, password string) error {
	if err := bcrypt.CompareHashAndPassword([]byte(password), []byte(user.Password)); err != nil {
		return fmt.Errorf("password is incorrect: %s, %v", err.Error(), objects.ErrAuthorization)
	}
	return nil
}
