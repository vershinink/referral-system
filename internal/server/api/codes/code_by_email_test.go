package codes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

func TestCodeByEmail(t *testing.T) {

	tests := []struct {
		name       string
		email      string
		wantStatus int
		mockError  error
	}{
		{
			name:       "OK",
			email:      "bob@gmail.com",
			wantStatus: http.StatusOK,
			mockError:  nil,
		},
		{
			name:       "Incorrect_Email",
			email:      "asdf",
			wantStatus: http.StatusBadRequest,
			mockError:  nil,
		},
		{
			name:       "Code_Not_Found",
			email:      "bob@gmail.com",
			wantStatus: http.StatusNotFound,
			mockError:  service.ErrCodeNotFound,
		},
		{
			name:       "Internal_Error",
			email:      "bob@gmail.com",
			wantStatus: http.StatusInternalServerError,
			mockError:  errors.New("DB error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			appMock := mocks.NewServiceCodes(t)

			// Исходя из тест-кейса устанавливаем поведение для мока
			// только если планируем дойти до него в тестируемой функции.
			if tt.wantStatus == http.StatusOK || tt.mockError != nil {
				appMock.On("CodeByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(testCode, tt.mockError).
					Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Get("/api/refcodes/email/{email}", CodeByEmail(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			uri := fmt.Sprintf("/api/refcodes/email/%s", tt.email)

			req := httptest.NewRequest(http.MethodGet, uri, nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("CodeByEmail() status = %v, want %v", rr.Code, tt.wantStatus)
			}
			if rr.Code == http.StatusOK {
				var resp CodeResponse
				err := json.NewDecoder(rr.Body).Decode(&resp)
				if err != nil {
					t.Fatalf("cannot decode response, error = %s", err.Error())
				}
				if resp.Code != testCode.Code {
					t.Errorf("CodeByEmail() response code = %v, want %v", resp.Code, testCode.Code)
				}
			}
		})
	}
}
