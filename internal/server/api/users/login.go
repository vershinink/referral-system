package users

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
)

// Login обрабатывает запрос на аутентификацию пользователя. Тело запроса
// должно содержать электронный адрес и пароль.
func Login(app ServiceUsers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.users.Login"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to login user")

		// Проверка заголовка Content-Type на значение application/json.
		ct := strings.ToLower(r.Header.Get("Content-Type"))
		if !strings.Contains(ct, "application/json") {
			log.Error("content-type is not application/json", slog.String("content-type", ct))
			http.Error(w, "unsupported content type", http.StatusUnsupportedMediaType)
			return
		}

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса.
		r.Body = http.MaxBytesReader(w, r.Body, maxRequest)
		var req LoginRequest
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

		// Получение JWT токена и данных пользователя.
		token, user, err := app.LoginUser(ctx, req.Email, req.Passwd)
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				http.Error(w, "user not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, service.ErrIncorrectPasswd) {
				http.Error(w, "incorrect password", http.StatusBadRequest)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		w.WriteHeader(http.StatusOK)
		var resp LoginResponse
		resp.AccessToken = token.AccessToken
		resp.UserID = user.ID
		resp.Email = user.Email
		resp.Created = user.Created
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode response", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}
		log.Info("user signed in successfully")
	}
}
