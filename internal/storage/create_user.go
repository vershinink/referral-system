package storage

import (
	"context"
	"fmt"
	"referral-rest-api/internal/models"
	"time"
)

// CreateUser создает нового пользователя не реферала в БД.
func (db *DB) CreateUser(ctx context.Context, user models.User) error {
	const operation = "storage.CreateUser"

	res, err := db.pool.Exec(
		ctx,
		`INSERT INTO users (email, passwd, created) VALUES ($1, $2, $3);`,
		user.Email,
		user.PassHash,
		time.Now().Unix(),
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, checkErr(err))
	}
	// Если по какой то причине вставка строки в таблицу не произошла,
	// то возвращаем неопределенную ошибку.
	if res.RowsAffected() == 0 {
		return fmt.Errorf("%s: %w", operation, ErrUndefined)
	}

	return nil
}
