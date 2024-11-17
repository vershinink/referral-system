package storage

import (
	"context"
	"testing"
)

func TestDB_UsersByReferrerID(t *testing.T) {

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
		name       string
		userID     int64
		wantArrLen int
		wantErr    bool
	}{

		{
			name:       "OK",
			userID:     1,
			wantArrLen: 1,
			wantErr:    false,
		},
		{
			name:       "No_Referrals",
			userID:     2,
			wantArrLen: 0,
			wantErr:    true,
		},
		{
			name:       "User_Not_Found",
			userID:     5,
			wantArrLen: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pool.UsersByReferrerID(context.Background(), tt.userID)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if len(got) != tt.wantArrLen {
				t.Errorf("DB.CreateCode() len = %v, want %v", len(got), tt.wantArrLen)
			}
		})
	}
}
