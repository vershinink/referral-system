package service

import (
	"context"
	"errors"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestApp_UsersByReferrer(t *testing.T) {

	tests := []struct {
		name      string
		userID    int64
		wantError bool
		mockError error
	}{
		{
			name:      "OK",
			userID:    1,
			wantError: false,
			mockError: nil,
		},
		{
			name:      "User_Not_Found",
			userID:    2,
			wantError: true,
			mockError: storage.ErrUserIdNotFound,
		},
		{
			name:      "Internal_Error",
			userID:    3,
			wantError: true,
			mockError: errors.New("DB error"),
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			dbMock := mocks.NewDatabase(t)

			// Исходя из тест-кейса устанавливаем поведение для мока
			// только если планируем дойти до него в тестируемой функции.
			if !tt.wantError || tt.mockError != nil {
				dbMock.On("UsersByReferrerID", mock.Anything, mock.AnythingOfType("int64")).
					Return(
						[]models.User{{ID: 1, Email: "bob@gmail.com"}},
						tt.mockError,
					).
					Once()
			}
			a := &App{
				log:        logger.SetupDiscard(),
				st:         dbMock,
				jwtSecret:  "qwerty",
				tokenTTL:   15 * time.Minute,
				refCodeTTL: 1 * time.Hour,
			}
			users, err := a.UsersByReferrer(context.Background(), tt.userID)
			if (err != nil) != tt.wantError {
				t.Errorf("App.CodeByID() error = %v, wantErr %v", err, tt.wantError)
			}
			if !tt.wantError {
				if len(users) == 0 {
					t.Errorf("App.CodeByID() error = empty users array")
				}
			}
		})
	}
}
