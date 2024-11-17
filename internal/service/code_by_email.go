package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"
)

// CodeByEmail возвращает реферальный код по его идентификатору.
func (a *App) CodeByEmail(ctx context.Context, email string) (models.RefCode, error) {
	const operation = "service.CodeByEmail"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("receiving referral code by email")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение реферального кода по email из БД.
	code, err := a.st.CodeByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrCodeNotFound) {
			log.Error("referral code from user with provided email not found", slog.String("email", email))
			return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		log.Error("internal DB error", logger.Err(err))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("referral code received successfully")
	return code, nil
}
