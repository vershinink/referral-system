package storage

import (
	"context"
	"errors"
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

// UserByEmail получает структуру пользователя по его email из БД.
func (db *DB) UserByEmail(ctx context.Context, email string) (models.User, error) {
	const operation = "storage.UserByEmail"

	var user models.User
	var created int64
	err := db.pool.QueryRow(
		ctx,
		`SELECT id, email, passwd, created FROM users WHERE email = $1;`,
		email,
	).Scan(
		&user.ID,
		&user.Email,
		&user.PassHash,
		&created,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", operation, ErrEmailNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", operation, err)
	}

	user.Created = time.Unix(created, 0)
	return user, nil
}
