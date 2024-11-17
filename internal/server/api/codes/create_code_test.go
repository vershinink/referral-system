// Пакет codes содержит обработчики http запросов, работающие с объектами
// реферальных кодов.
package codes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/service"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

var testCode = models.RefCode{
	ID:      1,
	Code:    "qwerty",
	OwnerID: 1,
	Created: time.Now(),
	Expired: time.Now().Add(1 * time.Hour),
	IsUsed:  false,
}

func TestCreate(t *testing.T) {

	tests := []struct {
		name        string
		userID      int64
		contentType string
		wantStatus  int
		mockError   error
	}{
		{
			name:        "OK",
			userID:      1,
			contentType: "Application/json",
			wantStatus:  http.StatusCreated,
			mockError:   nil,
		},
		{
			name:        "Code_Exists",
			userID:      2,
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   service.ErrCodeExits,
		},
		{
			name:        "User_Not_Found",
			userID:      3,
			contentType: "Application/json",
			wantStatus:  http.StatusNotFound,
			mockError:   service.ErrUserNotFound,
		},
		{
			name:        "Internal_Error",
			userID:      10,
			contentType: "Application/json",
			wantStatus:  http.StatusInternalServerError,
			mockError:   errors.New("DB error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewServiceCodes(t)

			// Исходя из тест-кейса устанавливаем поведение для мока
			// только если планируем дойти до него в тестируемой функции.
			if tt.wantStatus == http.StatusCreated || tt.mockError != nil {
				appMock.On("CreateCode", mock.Anything).
					Return(models.RefCode{}, tt.mockError).
					Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Post("/api/refcodes", Create(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			req := httptest.NewRequest(http.MethodPost, "/api/refcodes", nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
