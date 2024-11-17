package storage

import (
	"context"
	"testing"
)

func TestDB_UserByEmail(t *testing.T) {

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

	tests := []struct {
		name    string
		email   string
		wantID  int64
		wantErr bool
	}{

		{
			name:    "OK",
			email:   "bob@gmail.com",
			wantID:  1,
			wantErr: false,
		},
		{
			name:    "User_Not_Found",
			email:   "tom@gmail.com",
			wantID:  0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pool.UserByEmail(context.Background(), tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.ID != tt.wantID {
				t.Errorf("DB.CreateCode() = %v, want %v", got.ID, tt.wantID)
			}
		})
	}
}
