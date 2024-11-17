package users

import (
	"context"
	"log/slog"
	"referral-rest-api/internal/models"
	"time"
)

const (
	// maxRequest - ограничение размера тела запроса.
	maxRequest int64 = 1 << 20
)

// Service - интерфейс сервисного слоя.
//
//go:generate go run github.com/vektra/mockery/v2@v2.47.0 --name=ServiceUsers
type ServiceUsers interface {
	CreateUser(ctx context.Context, email, passwd, refCode string) error
	LoginUser(ctx context.Context, email, passwd string) (models.JWT, models.User, error)
	UsersByReferrer(ctx context.Context, userID int64) ([]models.User, error)
	LoggerSetup(ctx context.Context, operation string) *slog.Logger
}

// RegisterRequest - структура тела запроса на регистрацию пользователя.
type RegisterRequest struct {
	Email   string `json:"email" validate:"required,email,max=100"`
	Passwd  string `json:"password" validate:"required,min=6,max=72"`
	RefCode string `json:"refcode" validate:"omitempty,max=100"`
}

// LoginRequest - структура тела запроса аутентификации.
type LoginRequest struct {
	Email  string `json:"email" validate:"required,email,max=100"`
	Passwd string `json:"password" validate:"required,min=6,max=72"`
}

// UserResponse - структура ответа с данными о пользователе.
type UserResponse struct {
	UserID  int64     `json:"user_id"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}

// LoginResponse - структура тела ответа аутентификации.
type LoginResponse struct {
	AccessToken string `json:"access_token"`
	UserResponse
}
