package codes

import (
	"context"
	"log/slog"
	"referral-rest-api/internal/models"
	"time"
)

// ServiceCodes - интерфейс сервисного слоя.
//
//go:generate go run github.com/vektra/mockery/v2@v2.47.0 --name=ServiceCodes
type ServiceCodes interface {
	CreateCode(ctx context.Context) (models.RefCode, error)
	DeleteCode(ctx context.Context) error
	CodeByID(ctx context.Context) (models.RefCode, error)
	CodeByEmail(ctx context.Context, email string) (models.RefCode, error)
	LoggerSetup(ctx context.Context, operation string) *slog.Logger
}

// CodeResponse - структура тела ответа на создание кода.
type CodeResponse struct {
	ID      int64     `json:"id"`
	Code    string    `json:"code"`
	OwnerID int64     `json:"owner_id"`
	Created time.Time `json:"created"`
	Expired time.Time `json:"expired"`
	IsUsed  bool      `json:"is_used"`
}
