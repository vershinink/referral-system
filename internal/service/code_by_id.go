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

// CodeByID возвращает реферальный код пользователя. Функция получает
// идентификатор пользователя из payload JWT токена и проверяет его права.
// В случае успеха проверки прав возвращает последний реферльный код, если
// такой имеется.
func (a *App) CodeByID(ctx context.Context) (models.RefCode, error) {
	const operation = "service.CodeByID"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("receiving referral code by id")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение id пользователя из JWT токена.
	userID, err := idFromToken(ctx)
	if err != nil {
		log.Error("failed to receive id from jwt token", logger.Err(err))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	// Получение реферального кода по id из БД.
	code, err := a.st.CodeByID(ctx, userID)
	if err != nil {
		if errors.Is(err, storage.ErrCodeNotFound) {
			log.Error("referral code for user not found", slog.Int64("userID", userID))
			return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		log.Error("internal DB error", logger.Err(err))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("referral code received successfully")
	return code, nil
}
