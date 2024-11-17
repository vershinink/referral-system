package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"referral-rest-api/internal/jwt"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"

	"golang.org/x/crypto/bcrypt"
)

// LoginUser аутентифицирует пользователя. Метод проверяет совпадение
// электронного адреса и пароля в БД, в случае успеха генерирует JWT
// токен.
func (a *App) LoginUser(ctx context.Context, email, passwd string) (models.JWT, models.User, error) {
	const operation = "service.LoginUser"

	// Настройка логирования.
	log := a.LoggerSetup(ctx, operation)
	log.Debug("creating new user")

	// Проверка контекста запроса.
	if ctx.Err() != nil {
		log.Error("context error", logger.Err(ctx.Err()))
		return models.JWT{}, models.User{}, fmt.Errorf("%s: %w", operation, context.Canceled)
	}

	// Получение данных пользователя по его email из БД.
	user, err := a.st.UserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, storage.ErrEmailNotFound) {
			log.Error("user with provided email not found", slog.String("email", email))
			return models.JWT{}, models.User{}, fmt.Errorf("%s: %w", operation, ErrUserNotFound)
		}
		log.Error("internal DB error", logger.Err(err))
		return models.JWT{}, models.User{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	// Проверка совпадения хэшей паролей.
	err = bcrypt.CompareHashAndPassword(user.PassHash, []byte(passwd))
	if err != nil {
		log.Warn("provided password does not match", logger.Err(err))
		return models.JWT{}, models.User{}, fmt.Errorf("%s: %w", operation, ErrIncorrectPasswd)
	}

	// Генерация JWT токена.
	token, err := jwt.NewToken(user, a.jwtSecret, a.tokenTTL)
	if err != nil {
		return models.JWT{}, models.User{}, fmt.Errorf("%s: %w", operation, ErrInternal)
	}

	log.Debug("user data and token received successfully")
	return models.JWT{AccessToken: token}, user, nil
}
