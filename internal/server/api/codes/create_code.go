// Пакет codes содержит обработчики http запросов, работающие с объектами
// реферальных кодов.
package codes

import (
	"encoding/json"
	"errors"
	"net/http"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
)

// Create обрабатывает запрос на создание нового реферального кода. Тело
// запроса должно содержать ID пользователя, который создает код.
func Create(app ServiceCodes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.codes.Create"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to create new referral code")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Создание нового реферального кода.
		code, err := app.CreateCode(ctx)
		if err != nil {
			if errors.Is(err, service.ErrCodeExits) {
				http.Error(w, "valid referral code for user already exists", http.StatusBadRequest)
				return
			}
			if errors.Is(err, service.ErrUserNotFound) {
				http.Error(w, "user id not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		w.WriteHeader(http.StatusCreated)
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
		log.Info("new referral code created successfully")
	}
}
