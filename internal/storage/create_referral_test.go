package storage

import (
	"context"
	"referral-rest-api/internal/models"
	"testing"
	"time"
)

func TestDB_CreateReferral(t *testing.T) {

	pool, err := new(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	err = pool.truncate()
	if err != nil {
		t.Fatal(err)
	}

	err = pool.populate()
	if err != nil {
		t.Fatal(err)
	}

	// Остановка теста на 3 секунды, чтобы данные в БД немного устарели.
	time.Sleep(3 * time.Second)

	type args struct {
		email   string
		hash    []byte
		refCode string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "User_Exists",
			args:    args{email: "bob@gmail.com", hash: []byte("12345678"), refCode: "zxcvbn"},
			wantErr: true,
		},
		{
			name:    "OK",
			args:    args{email: "tom1@gmail.com", hash: []byte("12345678"), refCode: "zxcvbn"},
			wantErr: false,
		},
		{
			name:    "Code_Was_Used",
			args:    args{email: "tom2@gmail.com", hash: []byte("12345678"), refCode: "qwerty"},
			wantErr: true,
		},
		{
			name:    "Code_Expired",
			args:    args{email: "tom3@gmail.com", hash: []byte("12345678"), refCode: "asdfgh"},
			wantErr: true,
		},
		{
			name:    "Code_Not_Found",
			args:    args{email: "tom4@gmail.com", hash: []byte("12345678"), refCode: "mnbvcx"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{Email: tt.args.email, PassHash: tt.args.hash, RefCode: tt.args.refCode}
			err := pool.CreateReferral(context.Background(), user)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
