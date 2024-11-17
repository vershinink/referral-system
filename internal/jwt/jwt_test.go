package jwt

import (
	"referral-rest-api/internal/models"
	"testing"
	"time"
)

var secret string = "secret"

func TestNewToken(t *testing.T) {
	type args struct {
		user     models.User
		secret   string
		duration time.Duration
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "OK",
			args: args{
				user:     models.User{ID: 1, Email: "bob@gmail.com"},
				secret:   secret,
				duration: 15 * time.Minute,
			},
		},
		{
			name: "Empty_Secret",
			args: args{
				user:     models.User{ID: 1, Email: "bob@gmail.com"},
				secret:   "",
				duration: 15 * time.Minute,
			},
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewToken(tt.args.user, tt.args.secret, tt.args.duration)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == "" {
				t.Errorf("NewToken() = empty string")
			}
		})
	}
}
