package codes

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
)

// CodeByEmail обрабатывает запрос на получение реферального кода по email
// пользователя. URL-параметр email должен содержать валидный электронный
// адрес.
func CodeByEmail(app ServiceCodes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.codes.CodeByEmail"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to receive referral code by email")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса и валидация.
		email := chi.URLParam(r, "email")
		email, err := url.QueryUnescape(email)
		if err != nil {
			log.Error("failed to decode email url parameter")
			http.Error(w, "incorrect email", http.StatusBadRequest)
			return
		}
		email = strings.ToLower(email)
		valid := validator.New()
		err = valid.VarCtx(ctx, email, "required,email")
		if err != nil {
			log.Error("incorrect email")
			http.Error(w, "incorrect email", http.StatusBadRequest)
			return
		}

		// Получение реферального кода.
		code, err := app.CodeByEmail(ctx, email)
		if err != nil {
			if errors.Is(err, service.ErrCodeNotFound) {
				http.Error(w, "referral code not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		w.WriteHeader(http.StatusOK)
		var resp CodeResponse
		resp.ID = code.ID
		resp.Code = code.Code
		resp.OwnerID = code.OwnerID
		resp.Created = code.Created
		resp.Expired = code.Expired
		resp.IsUsed = code.IsUsed
		err = json.NewEncoder(w).Encode(resp)
		if err != nil {
			log.Error("failed to encode response", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("referral code received successfully")
	}
}
