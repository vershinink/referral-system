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

// UsersByReferrer возвращает массив пользователе-рефералов по id их реферера.
func (a *App) UsersByReferrer(ctx context.Context, userID int64) ([]models.User, error) {
	const operation = "service.UsersByReferrer"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("receiving referrals by referrer id")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return nil, fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение массива пользователей из БД.
	users, err := a.st.UsersByReferrerID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrUserIdNotFound) {
			log.Error("referrals array is empty or user not found", slog.Int64("userID", userID))
			return nil, fmt.Errorf("%s: %w", operation, ErrUserNotFound)
		}
		log.Error("internal DB error", logger.Err(err))
		return nil, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("referrals array received successfully")
	return users, nil
}
