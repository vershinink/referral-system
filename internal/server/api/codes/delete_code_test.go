package codes

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/service"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/mock"
)

func TestDelete(t *testing.T) {

	tests := []struct {
		name       string
		codeID     string
		wantStatus int
		mockError  error
	}{
		{
			name:       "OK",
			codeID:     "1",
			wantStatus: http.StatusNoContent,
			mockError:  nil,
		},
		{
			name:       "Code_Not_Found",
			codeID:     "2",
			wantStatus: http.StatusNotFound,
			mockError:  service.ErrCodeNotFound,
		},
		{
			name:       "Code_Was_Used",
			codeID:     "3",
			wantStatus: http.StatusForbidden,
			mockError:  service.ErrCodeWasUsed,
		},
		{
			name:       "Internal_Error",
			codeID:     "10",
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
			if tt.wantStatus == http.StatusNoContent || tt.mockError != nil {
				appMock.On("DeleteCode", mock.Anything).
					Return(tt.mockError).
					Once()
			}
			appMock.On("LoggerSetup", mock.Anything, mock.AnythingOfType("string")).
				Return(logger.SetupDiscard()).
				Once()

			mux := chi.NewRouter()
			mux.Delete("/api/refcodes", Delete(appMock))
			srv := httptest.NewServer(mux)
			defer srv.Close()

			req := httptest.NewRequest(http.MethodDelete, "/api/refcodes", nil)
			rr := httptest.NewRecorder()

			mux.ServeHTTP(rr, req)

			if rr.Code != tt.wantStatus {
				t.Errorf("Delete() status = %v, want %v", rr.Code, tt.wantStatus)
			}
		})
	}
}
