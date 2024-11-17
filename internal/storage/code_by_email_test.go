package storage

import (
	"context"
	"testing"
)

func TestDB_CodeByEmail(t *testing.T) {

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
		email    string
		wantCode string
		wantErr  bool
	}{

		{
			name:     "OK",
			email:    "jane@gmail.com",
			wantCode: "zxcvbn",
			wantErr:  false,
		},
		{
			name:     "OK_Expired",
			email:    "bill@gmail.com",
			wantCode: "asdfgh",
			wantErr:  false,
		},
		{
			name:     "Was_Used",
			email:    "bob@gmail.com",
			wantCode: "",
			wantErr:  true,
		},
		{
			name:     "Code_Not_Found",
			email:    "jill@gmail.com",
			wantCode: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pool.CodeByEmail(context.Background(), tt.email)
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
