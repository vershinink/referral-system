package service

import (
	"bytes"
	"context"
	"errors"
	"referral-rest-api/internal/logger"
	"referral-rest-api/internal/mocks"
	"referral-rest-api/internal/models"
	"referral-rest-api/internal/storage"
	"testing"
	"time"

	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestApp_LoginUser(t *testing.T) {

	hash, _ := bcrypt.GenerateFromPassword([]byte("12345678"), bcrypt.MinCost)
	errPass := errors.New("err")

	type args struct {
		email  string
		passwd string
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
				email:  "bob@gmail.com",
				passwd: "12345678",
			},
			wantError: false,
			mockError: nil,
		},
		{
			name: "User_Not_Found",
			args: args{
				email:  "bob2@gmail.com",
				passwd: "12345678",
			},
			wantError: true,
			mockError: storage.ErrEmailNotFound,
		},
		{
			name: "Incorrect_Password",
			args: args{
				email:  "bob@gmail.com",
				passwd: "00000000",
			},
			wantError: true,
			mockError: errPass,
		},
		{
			name: "Internal_Error",
			args: args{
				email:  "bob3@gmail.com",
				passwd: "12345678",
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
				dbMock.On("UserByEmail", mock.Anything, mock.AnythingOfType("string")).
					Return(func(ctx context.Context, email string) (models.User, error) {
						user := models.User{ID: 1, Email: "bob@gmail.com", PassHash: hash}
						if tt.mockError == errPass {
							return user, nil
						}
						return user, tt.mockError
					}).
					Once()
			}
			a := &App{
				log:        logger.SetupDiscard(),
				st:         dbMock,
				jwtSecret:  "qwerty",
				tokenTTL:   15 * time.Minute,
				refCodeTTL: 1 * time.Hour,
			}
			token, user, err := a.LoginUser(context.Background(), tt.args.email, tt.args.passwd)
			if (err != nil) != tt.wantError {
				t.Errorf("App.CreateUser() error = %v, wantErr %v", err, tt.wantError)
			}
			if !tt.wantError {
				if token.AccessToken == "" {
					t.Errorf("App.CreateUser() error = empty token string")
				}
				if !bytes.Equal(user.PassHash, hash) {
					t.Errorf("App.CreateUser() error = incorrect hash")
				}
			}
		})
	}
}
