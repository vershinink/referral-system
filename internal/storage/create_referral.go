package storage

import (
	"context"
	"errors"
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateReferral создает нового пользователя-реферала в БД, используя
// переданный в структуре реферальный код.
func (db *DB) CreateReferral(ctx context.Context, user models.User) error {
	const operation = "storage.CreateReferral"

	// Начало транзакции.
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	defer tx.Rollback(ctx)

	// Получение структуры реферального кода по значению кода.
	var code models.RefCode
	var created, expired int64
	err = tx.QueryRow(
		ctx,
		`SELECT id, code, owner, created, expired, is_used FROM codes WHERE code = $1;`,
		user.RefCode,
	).Scan(
		&code.ID,
		&code.Code,
		&code.OwnerID,
		&created,
		&expired,
		&code.IsUsed,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Проверка срока годности и использованности кода.
	if code.IsUsed {
		return fmt.Errorf("%s: %w", operation, ErrCodeWasUsed)
	}
	if time.Now().Unix() > expired {
		return fmt.Errorf("%s: %w", operation, ErrCodeExpired)
	}

	// Создание нового пользователя в таблице users.
	err = tx.QueryRow(
		ctx,
		`INSERT INTO users (email, passwd, created) VALUES ($1, $2, $3) RETURNING id;`,
		user.Email,
		user.PassHash,
		time.Now().Unix(),
	).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, checkErr(err))
	}

	// Создание записи реферала в таблице referrals.
	_, err = tx.Exec(
		ctx,
		`INSERT INTO referrals (id, referrer_id, code_id) VALUES ($1, $2, $3);`,
		user.ID,
		code.OwnerID,
		code.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Установка флага использованности реферального кода в таблице codes.
	_, err = tx.Exec(
		ctx,
		`UPDATE codes SET is_used = true WHERE id = $1`,
		code.ID,
	)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}

	// Завершение транзакции.
	err = tx.Commit(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	return nil
}
