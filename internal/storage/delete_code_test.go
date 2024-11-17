package storage

import (
	"context"
	"testing"
)

func TestDB_DeleteCode(t *testing.T) {

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
		userID  int64
		wantErr bool
	}{

		{
			name:    "OK",
			userID:  3,
			wantErr: false,
		},
		{
			name:    "OK_Expired",
			userID:  2,
			wantErr: false,
		},
		{
			name:    "Error_Was_Used",
			userID:  1,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := pool.DeleteCode(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
