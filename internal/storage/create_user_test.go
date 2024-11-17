package storage

import (
	"context"
	"referral-rest-api/internal/models"
	"testing"
)

func TestDB_CreateUser(t *testing.T) {

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

	type args struct {
		email string
		hash  []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "OK",
			args:    args{email: "bob2@gmail.com", hash: []byte("12345678")},
			wantErr: false,
		},
		{
			name:    "User_Exists",
			args:    args{email: "bob2@gmail.com", hash: []byte("12345678")},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := models.User{Email: tt.args.email, PassHash: tt.args.hash}
			err := pool.CreateUser(context.Background(), user)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
