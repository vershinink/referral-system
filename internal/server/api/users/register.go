// Пакет users содержит обработчики http запросов, работающие с объектами
// пользователей.
package users

import (
	"errors"
	"log/slog"
	"net/http"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Register обрабатывает запрос на создание нового пользователя. Если тело запроса
// содержит реферальный код, то будет создан пользователь-реферал.
func Register(app ServiceUsers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.users.Register"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to register new user")

		// Проверка заголовка Content-Type на значение application/json.
		ct := strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(ct, "application/json") {
			log.Error("content-type is not application/json", slog.String("content-type", ct))
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		// Получение параметров запроса.
		r.Body = http.MaxBytesReader(w, r.Body, maxRequest)
		var req RegisterRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			log.Error("cannot decode request json", logger.Err(err))
			http.Error(w, "incorrect input data", http.StatusBadRequest)
			return
		}

		// Валидация полей запроса.
		valid := validator.New()
		err = valid.StructCtx(ctx, req)
		if err != nil {
			validateErr := err.(validator.ValidationErrors)
			log.Error("validation failed", logger.Err(validateErr))
			http.Error(w, "incorrect input data", http.StatusBadRequest)
			return
		}

		// Создание нового пользователя.
		err = app.CreateUser(ctx, req.Email, req.Passwd, req.RefCode)
		if err != nil {
			if errors.Is(err, service.ErrUserExists) {
				http.Error(w, "user with provided email already exists", http.StatusBadRequest)
				return
			}
			if errors.Is(err, service.ErrCodeNotFound) {
				http.Error(w, "referral code not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusCreated)
		log.Info("new user registered successfully")
	}
}
