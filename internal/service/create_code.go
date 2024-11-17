package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"
	"time"
)

// CreateCode создает новый реферальный код. Функция получает идентификатор
// создающего код пользователя из payload JWT токена и проверяет его права.
// В случае успеха проверки прав и отсутствия другого валидного кода у этого
// пользователя - создает новый.
func (a *App) CreateCode(ctx context.Context) (models.RefCode, error) {
	const operation = "service.CreateCode"
	var code models.RefCode

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("creating new referral code")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return code, fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение id пользователя из JWT токена.
	userID, err := idFromToken(ctx)
	if err != nil {
		log.Error("failed to receive id from jwt token", logger.Err(err))
		return code, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	// Генерация уникального реферального кода.
	gen := generateCode()
	if gen == "" {
		log.Error("failed to generate referral code")
		return code, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	// Формирование структуры реферального кода.
	code.OwnerID = userID
	code.Code = gen
	code.Created = time.Now()
	code.Expired = code.Created.Add(a.refCodeTTL)

	// Создание нового реферального кода в БД.
	id, err := a.st.CreateCode(ctx, code)
	if err != nil {
		if errors.Is(err, storage.ErrCodeExits) {
			log.Warn(
				"valid referral code for provided userID already exists",
				slog.Int64("userID", userID),
				slog.Int64("codeID", id),
			)
			return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrCodeExits)
		}
		if errors.Is(err, storage.ErrUserIdNotFound) {
			log.Error("user with provided id not found", slog.Int64("userID", userID))
			return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrUserNotFound)
		}
		log.Error("internal DB error", logger.Err(err))
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	code.ID = id

	log.Debug("referral code created successfully")
	return code, nil
}
