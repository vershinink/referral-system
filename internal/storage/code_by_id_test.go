package storage

import (
	"context"
	"testing"
)

func TestDB_CodeByID(t *testing.T) {

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
		name     string
		userID   int64
		wantCode string
		wantErr  bool
	}{

		{
			name:     "OK",
			userID:   3,
			wantCode: "zxcvbn",
			wantErr:  false,
		},
		{
			name:     "OK_Expired",
			userID:   2,
			wantCode: "asdfgh",
			wantErr:  false,
		},
		{
			name:     "OK_Was_Used",
			userID:   1,
			wantCode: "qwerty",
			wantErr:  false,
		},
		{
			name:     "Code_Not_Found",
			userID:   4,
			wantCode: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pool.CodeByID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.Code != tt.wantCode {
				t.Errorf("DB.CreateCode() = %v, want %v", got.Code, tt.wantCode)
			}
		})
	}
}
