// Пакет service содержит структуру App и ее методы, которые реализуют
// бизнес-логику веб приложения.
package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/config"
	"referral-rest-api/internal/models"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/jwtauth/v5"
	"github.com/lestrrat-go/jwx/v2/jwt"
	"github.com/pzentenoe/go-cache"
	"github.com/sqids/sqids-go"
)

var (
	ErrUserExists      = errors.New("user with provided email already exists")
	ErrUserNotFound    = errors.New("user not found")
	ErrCodeExits       = errors.New("valid referral code already exists")
	ErrIncorrectPasswd = errors.New("incorrect password")
	ErrCodeNotFound    = errors.New("referral code not found")
	ErrCodeWasUsed     = errors.New("referral code was used")
	ErrCodeExpired     = errors.New("referral code expired")
	ErrInternal        = errors.New("internal error")
)

// App - структура сервисного слоя приложения, осуществляет бизнес-логику.
type App struct {
	log        *slog.Logger
	st         Database
	codeCache  *cache.Cache
	jwtAuth    *jwtauth.JWTAuth
	jwtSecret  string
	tokenTTL   time.Duration
	refCodeTTL time.Duration
}

// Database - интерфейс пула подключений БД.
//
//go:generate go run github.com/vektra/mockery/v2@v2.47.0 --name=Database
type Database interface {
	CreateUser(ctx context.Context, user models.User) error
	CreateReferral(ctx context.Context, user models.User) error
	UserByEmail(ctx context.Context, email string) (models.User, error)
	UsersByReferrerID(ctx context.Context, refID int64) ([]models.User, error)

	CreateCode(ctx context.Context, code models.RefCode) (int64, error)
	DeleteCode(ctx context.Context, userID int64) error
	CodeByID(ctx context.Context, userID int64) (models.RefCode, error)
	CodeByEmail(ctx context.Context, email string) (models.RefCode, error)
}

// New - конструктор сервиса.
func New(cfg *config.Config, log *slog.Logger, st Database) *App {
	cc := cache.New(cache.NoExpiration, cache.NoExpiration)

	j := jwtauth.New("HS256", []byte(cfg.JwtSecret), nil, jwt.WithAcceptableSkew(30*time.Second))
	j.ValidateOptions()

	a := &App{
		log:        log,
		st:         st,
		codeCache:  cc,
		jwtAuth:    j,
		jwtSecret:  cfg.JwtSecret,
		tokenTTL:   cfg.TokenTTL,
		refCodeTTL: cfg.CodeTTL,
	}
	return a
}

// LoggerSetup дополняет логирование операцией и id запроса.
func (a *App) LoggerSetup(ctx context.Context, operation string) *slog.Logger {
	log := a.log.With(
		slog.String("op", operation),
		slog.String("request_id", middleware.GetReqID(ctx)),
	)
	return log
}

// JWTauth возвращает указатель на JWTAuth.
func (a *App) JWTauth() *jwtauth.JWTAuth {
	return a.jwtAuth
}

// generateCode генерирует уникальное стрковое значение из текущей
// временной метки.
func generateCode() string {
	tm := time.Now()
	sec := uint64(tm.Unix())
	nano := uint64(tm.Nanosecond())

	s, err := sqids.New(sqids.Options{MinLength: 16})
	if err != nil {
		return ""
	}
	code, err := s.Encode([]uint64{sec, nano})
	if err != nil {
		return ""
	}
	return code
}

// idFromToken получает идентификатор пользователя из токена.
func idFromToken(ctx context.Context) (int64, error) {
	const operation = "service.accessFromToken"
	var userID int64

	_, claims, err := jwtauth.FromContext(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	userID = int64(claims["uid"].(float64))

	return userID, nil
}
