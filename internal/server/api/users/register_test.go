// Пакет users содержит обработчики http запросов, работающие с объектами
// пользователей.
package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

func TestRegister(t *testing.T) {

	tests := []struct {
		name        string
		req         RegisterRequest
		contentType string
		wantStatus  int
		mockError   error
	}{
		{
			name:        "OK",
			req:         RegisterRequest{Email: "bob@gmail.com", Passwd: "12345678", RefCode: "qwerty"},
			contentType: "Application/json",
			wantStatus:  http.StatusCreated,
			mockError:   nil,
		},
		{
			name:        "Incorrect_Content_Type",
			req:         RegisterRequest{Email: "bob@gmail.com", Passwd: "12345678", RefCode: "qwerty"},
			contentType: "",
			wantStatus:  http.StatusUnsupportedMediaType,
			mockError:   nil,
		},
		{
			name:        "Incorrect_Email",
			req:         RegisterRequest{Email: "bob", Passwd: "12345678", RefCode: "qwerty"},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   nil,
		},
		{
			name:        "Incorrect_Password",
			req:         RegisterRequest{Email: "bob@gmail.com", Passwd: "", RefCode: "qwerty"},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   nil,
		},
		{
			name:        "User_Exists",
			req:         RegisterRequest{Email: "bob@gmail.com", Passwd: "12345678", RefCode: "qwerty"},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   service.ErrUserExists,
		},
		{
			name:        "Internal_Error",
			req:         RegisterRequest{Email: "bob@gmail.com", Passwd: "12345678", RefCode: "qwerty"},
			contentType: "Application/json",
			wantStatus:  http.StatusInternalServerError,
			mockError:   service.ErrInternal,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewServiceUsers(t)

			// Исходя из тест-кейса устанавливаем поведение для мока
			// только если планируем дойти до него в тестируемой функции.
			if tt.wantStatus == http.StatusCreated || tt.mockError != nil {
				appMock.On(
					"CreateUser",
					mock.Anything,
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
				).Return(tt.mockError).Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Post("/api/users", Register(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			b, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("cannot encode new request, error = %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, "/api/users", bytes.NewReader(b))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
