package codes

import (
	"encoding/json"
	"errors"
	"net/http"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
)

// CodeByID обрабатывает запрос на получение реферального кода по его
// идентификатору. URL-параметр id должен содержать идентификатор
// запрашиваемого кода.
func CodeByID(app ServiceCodes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.codes.CodeByID"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to receive referral code by id")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение реферального кода.
		code, err := app.CodeByID(ctx)
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
