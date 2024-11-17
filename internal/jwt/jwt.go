package jwt

import (
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// NewToken создает новый JWT токен и заполняет его payload.
func NewToken(user models.User, secret string, dur time.Duration) (string, error) {
	const operation = "jwt.NewToken"

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["uid"] = user.ID
	claims["email"] = user.Email
	claims["exp"] = time.Now().Add(dur).Unix()

	tokenStr, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("%s: %w", operation, err)
	}

	return tokenStr, nil
}
