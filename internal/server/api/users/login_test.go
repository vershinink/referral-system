package users

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

func TestLogin(t *testing.T) {

	tests := []struct {
		name        string
		req         LoginRequest
		contentType string
		wantStatus  int
		mockError   error
	}{
		{
			name:        "OK",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: "12345678"},
			contentType: "Application/json",
			wantStatus:  http.StatusOK,
			mockError:   nil,
		},
		{
			name:        "Incorrect_Content_Type",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: "12345678"},
			contentType: "",
			wantStatus:  http.StatusUnsupportedMediaType,
			mockError:   nil,
		},
		{
			name:        "Incorrect_Email",
			req:         LoginRequest{Email: "bob", Passwd: "12345678"},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   nil,
		},
		{
			name:        "Empty_Password",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: ""},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   nil,
		},
		{
			name:        "User_Not_Found",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: "12345678"},
			contentType: "Application/json",
			wantStatus:  http.StatusNotFound,
			mockError:   service.ErrUserNotFound,
		},
		{
			name:        "Incorrect_Password",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: "00000000"},
			contentType: "Application/json",
			wantStatus:  http.StatusBadRequest,
			mockError:   service.ErrIncorrectPasswd,
		},
		{
			name:        "Internal_Error",
			req:         LoginRequest{Email: "bob@gmail.com", Passwd: "12345678"},
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
			if tt.wantStatus == http.StatusOK || tt.mockError != nil {
				appMock.On(
					"LoginUser",
					mock.Anything,
					mock.AnythingOfType("string"),
					mock.AnythingOfType("string"),
				).Return(
					models.JWT{AccessToken: "qwerty"},
					models.User{ID: 1, Email: tt.req.Email, PassHash: []byte(tt.req.Passwd)},
					tt.mockError,
				).Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Post("/api/users/login", Login(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			b, err := json.Marshal(tt.req)
			if err != nil {
				t.Fatalf("cannot encode new request, error = %s", err.Error())
			}

			req := httptest.NewRequest(http.MethodPost, "/api/users/login", bytes.NewReader(b))
			req.Header.Set("Content-Type", tt.contentType)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
