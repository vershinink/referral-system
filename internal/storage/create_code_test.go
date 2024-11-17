package storage

import (
	"context"
	"referral-rest-api/internal/models"
	"testing"
	"time"
)

func TestDB_CreateCode(t *testing.T) {

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

	tests := []struct {
		name       string
		code       models.RefCode
		wantCodeID int64
		wantErr    bool
	}{

		{
			name: "OK_With_Used",
			code: models.RefCode{
				Code:    "qwerty1",
				OwnerID: 1,
				Created: time.Now(),
				Expired: time.Now().Add(1 * time.Hour),
			},
			wantCodeID: 4,
			wantErr:    false,
		},
		{
			name: "OK_With_Expired",
			code: models.RefCode{
				Code:    "qwerty2",
				OwnerID: 2,
				Created: time.Now(),
				Expired: time.Now().Add(1 * time.Hour),
			},
			wantCodeID: 5,
			wantErr:    false,
		},
		{
			name: "OK_No_Codes",
			code: models.RefCode{
				Code:    "qwerty4",
				OwnerID: 4,
				Created: time.Now(),
				Expired: time.Now().Add(1 * time.Hour),
			},
			wantCodeID: 6,
			wantErr:    false,
		},
		{
			name: "Valid_Code_Exists",
			code: models.RefCode{
				Code:    "qwerty3",
				OwnerID: 3,
				Created: time.Now(),
				Expired: time.Now().Add(1 * time.Hour),
			},
			wantCodeID: 3,
			wantErr:    true,
		},
		{
			name: "User_Not_Found",
			code: models.RefCode{
				Code:    "qwerty5",
				OwnerID: 5,
				Created: time.Now(),
				Expired: time.Now().Add(1 * time.Hour),
			},
			wantCodeID: 0,
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := pool.CreateCode(context.Background(), tt.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("DB.CreateCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantCodeID {
				t.Errorf("DB.CreateCode() = %v, want %v", got, tt.wantCodeID)
			}
		})
	}
}
