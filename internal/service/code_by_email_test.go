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

var testCode = models.RefCode{
	ID:      1,
	Code:    "qwerty",
	OwnerID: 1,
	Created: time.Now(),
	Expired: time.Now().Add(1 * time.Hour),
	IsUsed:  false,
}

func TestApp_CodeByEmail(t *testing.T) {

	tests := []struct {
		name      string
		email     string
		wantError bool
		mockError error
	}{
		{
			name:      "OK",
			email:     "bob@gmail.com",
			wantError: false,
			mockError: nil,
		},
		{
			name:      "Code_Not_Found",
			email:     "bob2@gmail.com",
			wantError: true,
			mockError: storage.ErrCodeNotFound,
		},
		{
			name:      "Internal_Error",
			email:     "bob3@gmail.com",
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
				dbMock.On("CodeByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(testCode, tt.mockError).
					Once()
			}
			a := &App{
				log:        logger.SetupDiscard(),
				st:         dbMock,
				jwtSecret:  "qwerty",
				tokenTTL:   15 * time.Minute,
				refCodeTTL: 1 * time.Hour,
			}
			code, err := a.CodeByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantError {
				t.Errorf("App.CodeByEmail() error = %v, wantErr %v", err, tt.wantError)
			}
			if !tt.wantError && code.ID != testCode.ID {
				t.Errorf("App.CodeByEmail() code.ID = %v, wantErr %v", code.ID, testCode.ID)
			}
		})
	}
}
