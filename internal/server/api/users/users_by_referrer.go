package users

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/service"
	"strconv"

	"github.com/go-chi/chi/v5"
)

// UsersByReferrer обрабатывает запрос на получение списка пользователей
// рефералов по id их реферера.
func UsersByReferrer(app ServiceUsers) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.users.UsersByReferrer"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to receive referrals by referrer id")

		// Установка типа контента для ответа.
		w.Header().Set("Content-Type", "application/json")

		// Получение параметров запроса и валидация.
		param := chi.URLParam(r, "id")
		userID, err := strconv.Atoi(param)
		if err != nil || userID < 1 {
			log.Error("incorrect user id", slog.String("id", param), logger.Err(err))
			http.Error(w, "incorrect user id", http.StatusBadRequest)
			return
		}

		// Получение массива пользователей-рефералов.
		users, err := app.UsersByReferrer(ctx, int64(userID))
		if err != nil {
			if errors.Is(err, service.ErrUserNotFound) {
				http.Error(w, "users not found", http.StatusNotFound)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		// Кодирование ответа в JSON.
		var response []UserResponse
		for _, user := range users {
			var resp UserResponse
			resp.UserID = user.ID
			resp.Email = user.Email
			resp.Created = user.Created
			response = append(response, resp)
		}

		w.WriteHeader(http.StatusOK)
		err = json.NewEncoder(w).Encode(response)
		if err != nil {
			log.Error("failed to encode response", logger.Err(err))
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		log.Info("refarrals by referrer id received successfully")
	}
}
