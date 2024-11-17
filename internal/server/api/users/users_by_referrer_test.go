package users

import (
	"encoding/json"
	"fmt"
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

func TestUsersByReferrer(t *testing.T) {

	tests := []struct {
		name       string
		userID     string
		wantArrLen int
		wantStatus int
		mockError  error
	}{
		{
			name:       "OK",
			userID:     "1",
			wantArrLen: 1,
			wantStatus: http.StatusOK,
			mockError:  nil,
		},
		{
			name:       "Incorrect_UserID",
			userID:     "asdf",
			wantArrLen: 0,
			wantStatus: http.StatusBadRequest,
			mockError:  nil,
		},
		{
			name:       "User_Not_Found",
			userID:     "2",
			wantArrLen: 0,
			wantStatus: http.StatusNotFound,
			mockError:  service.ErrUserNotFound,
		},
		{
			name:       "Internal_Error",
			userID:     "3",
			wantArrLen: 0,
			wantStatus: http.StatusInternalServerError,
			mockError:  service.ErrInternal,
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
					"UsersByReferrer",
					mock.Anything,
					mock.AnythingOfType("int64"),
				).Return(
					[]models.User{{ID: 5, Email: "bob@gmail.com", Created: time.Now()}},
					tt.mockError,
				).Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Get("/api/users/{id}", UsersByReferrer(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			uri := fmt.Sprintf("/api/users/%s", tt.userID)

			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Create() status = %v, want %v", rr.Code, tt.wantStatus)
			}

			if rr.Code == http.StatusOK {
				var resp []UserResponse
				err := json.NewDecoder(rr.Body).Decode(&resp)
				if err != nil {
					t.Fatalf("cannot decode response, error = %s", err.Error())
				}
				if len(resp) != tt.wantArrLen {
					t.Errorf("CodeByID() response len = %v, want %v", len(resp), tt.wantArrLen)
				}
			}
		})
	}
}
