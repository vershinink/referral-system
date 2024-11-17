package service

import (
	"context"
	"errors"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
)

func TestApp_CreateUser(t *testing.T) {

	type args struct {
		email   string
		passwd  string
		refCode string
	}
	tests := []struct {
		name      string
		args      args
		wantError bool
		mockError error
	}{
		{
			name: "OK",
			args: args{
				email:   "bob@gmail.com",
				passwd:  "12345678",
				refCode: "",
			},
			wantError: false,
			mockError: nil,
		},
		{
			name: "User_Exists",
			args: args{
				email:   "bob2@gmail.com",
				passwd:  "12345678",
				refCode: "",
			},
			wantError: true,
			mockError: storage.ErrEmailExists,
		},
		{
			name: "Internal_Error",
			args: args{
				email:   "bob3@gmail.com",
				passwd:  "12345678",
				refCode: "",
			},
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
				dbMock.On("CreateUser", mock.Anything, mock.AnythingOfType("models.User")).
					Return(tt.mockError).
					Once()
			}
			a := &App{
				log:        logger.SetupDiscard(),
				st:         dbMock,
				jwtSecret:  "qwerty",
				tokenTTL:   15 * time.Minute,
				refCodeTTL: 1 * time.Hour,
			}
			err := a.CreateUser(context.Background(), tt.args.email, tt.args.passwd, tt.args.refCode)
			if (err != nil) != tt.wantError {
				t.Errorf("App.CreateUser() error = %v, wantErr %v", err, tt.wantError)
			}
		})
	}
}
