package codes

import (
	"errors"
	"net/http"
	"referral-rest-api/internal/service"
)

// Delete обрабатывает запрос на удаление неиспользованного реферального
// кода. URL-параметр id должен содержать идентификатор удаляемого кода.
func Delete(app ServiceCodes) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const operation = "server.codes.Delete"

		// Настройка логирования.
		ctx := r.Context()
		log := app.LoggerSetup(ctx, operation)
		log.Info("request to delete referral code")

		// Удаление реферального кода.
		err := app.DeleteCode(ctx)
		if err != nil {
			if errors.Is(err, service.ErrCodeNotFound) {
				http.Error(w, "referral code not found", http.StatusNotFound)
				return
			}
			if errors.Is(err, service.ErrCodeWasUsed) {
				http.Error(w, "referral code was used", http.StatusForbidden)
				return
			}
			http.Error(w, "internal error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusNoContent)
		log.Info("referral code deleted successfully")
	}
}
