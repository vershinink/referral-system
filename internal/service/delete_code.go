package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/storage"
)

// DeleteCode удаляет неиспользованный реферальный код у пользователя.
// Функция получает идентификатор удаляющего код пользователя из payload
// JWT токена и проверяет его права. В случае успеха проверки прав удаляет
// неиспользованный код, если такой имеется.
func (a *App) DeleteCode(ctx context.Context) error {
	const operation = "service.DeleteCode"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("deleting referral code")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение id пользователя из JWT токена.
	userID, err := idFromToken(ctx)
	if err != nil {
		log.Error("failed to receive id from jwt token", logger.Err(err))
		return fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	// Удаление реферального кода в БД.
	err = a.st.DeleteCode(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrCodeNotFound) {
			log.Error("referral code for user not found", slog.Int64("userID", userID))
			return fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		if errors.Is(err, storage.ErrCodeWasUsed) {
			log.Error("referral code for user was used already", slog.Int64("userID", userID))
			return fmt.Errorf("%s: %w", operation, ErrCodeWasUsed)
		}
		log.Error("internal DB error", logger.Err(err))
		return fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("referral code deleted successfully")
	return nil
}
