package storage

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
)

// DeleteCode удаляет неиспользованный реферальный код из БД по
// id его пользователя.
func (db *DB) DeleteCode(ctx context.Context, userID int64) error {
	const operation = "storage.DeleteCode"

	// Начало транзакции.
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("%s: %w", operation, err)
	}
	defer tx.Rollback(ctx)

	// Проверка на существование кода в БД, и был ли он использован
	// для регистрации реферала. Использованные коды удалять нельзя,
	// так как это нарушит консистентность БД.
	var codeID int64
	var used bool
	err = tx.QueryRow(
		ctx,
		`SELECT id, is_used FROM codes WHERE owner = $1 
		 ORDER BY created DESC LIMIT 1;`,
		userID,
	).Scan(&codeID, &used)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		return fmt.Errorf("%s: %w", operation, err)
	}

	if used {
		return fmt.Errorf("%s: %w", operation, ErrCodeWasUsed)
	}

	// Удаление неиспользованного кода.
	_, err = tx.Exec(
		ctx,
		`DELETE FROM codes WHERE id = $1`,
		codeID,
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
