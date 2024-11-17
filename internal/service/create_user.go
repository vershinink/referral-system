package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

// CreateUser создает нового пользователя. Если передан refCode, то создает
// пользователя-реферала с использованием этого значения.
func (a *App) CreateUser(ctx context.Context, email, passwd, refCode string) error {
	const operation = "service.CreateUser"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("creating new user")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Генерация хэша из переданного пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(passwd), bcrypt.MinCost)
	if err != nil {
		log.Error("failed to generate password hash", logger.Err(err))
		return fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	user := models.User{
		Email:    email,
		PassHash: passHash,
		RefCode:  refCode,
	}

	// Проверка на использование реферального кода при создании пользователя.
	switch user.RefCode {
	case "":
		err = a.st.CreateUser(ctx, user)
	default:
		err = a.st.CreateReferral(ctx, user)
	}
	if err != nil {
		if errors.Is(err, storage.ErrEmailExists) {
			log.Warn("user with provided email already exists", slog.String("email", user.Email))
			return fmt.Errorf("%s: %w", operation, ErrUserExists)
		}
		if errors.Is(err, storage.ErrCodeNotFound) {
			log.Error("referral code not found", slog.String("refCode", user.RefCode))
			return fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		if errors.Is(err, storage.ErrCodeWasUsed) {
			log.Error("referral code was used already", slog.String("refCode", user.RefCode))
			return fmt.Errorf("%s: %w", operation, ErrCodeWasUsed)
		}
		if errors.Is(err, storage.ErrCodeExpired) {
			log.Error("referral code expired", slog.String("refCode", user.RefCode))
			return fmt.Errorf("%s: %w", operation, ErrCodeExpired)
		}
		log.Error("internal DB error", logger.Err(err))
		return fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("user created successfully")
	return nil
}
