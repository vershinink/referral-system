package storage

import (
	"context"
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

// CreateCode создает новый реферальный код в БД. При успехе возвращает его id.
func (db *DB) CreateCode(ctx context.Context, code models.RefCode) (int64, error) {
	const operation = "storage.CreateCode"

	// Начало транзакции.
	tx, err := db.pool.Begin(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	defer tx.Rollback(ctx)

	// Получение последнего неиспользованного кода пользователя. Если
	// код получен, то проверяется его срок годности. Если такого кода
	// нет, значит либо у пользователя еще нет созданных кодов, либо
	// они все использованы. Либо пользователя с таким id не существует.
	var lastID, exp int64
	err = tx.QueryRow(
		ctx,
		`SELECT id, expired FROM codes WHERE owner = $1 AND is_used = false ORDER BY created DESC LIMIT 1;`,
		code.OwnerID,
	).Scan(
		&lastID,
		&exp,
	)
	// Ошибка возникает, если код не получен.
	if err != nil {
		if err != pgx.ErrNoRows {
			return 0, fmt.Errorf("%s: %w", operation, err)
		}
		// Создание нового кода и получение его id.
		id, err := db.createCode(ctx, tx, code)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", operation, checkErr(err))
		}
		// При успешном создании кода завершение транзакции и возврат id
		// созданного кода.
		err = tx.Commit(ctx)
		if err != nil {
			return 0, fmt.Errorf("%s: %w", operation, err)
		}
		return id, nil
	}

	// При успешном получении последнего кода проверяется его срок годности.
	// Если срок не вышел, то возвращается ошибка.
	if time.Now().Unix() < exp {
		return lastID, fmt.Errorf("%s: %w", operation, ErrCodeExits)
	}

	// Если срок годности вышел, то старый код удаляется, и создается новый.
	// Таким образом у каждого пользователя одновременно может быть только
	// один валидный код.
	res, err := tx.Exec(
		ctx,
		`DELETE FROM codes WHERE id = $1;`,
		lastID,
	)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	if res.RowsAffected() == 0 {
		return 0, fmt.Errorf("%s: %w", operation, ErrUndefined)
	}

	id, err := db.createCode(ctx, tx, code)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, checkErr(err))
	}

	// Завершение транзакции.
	err = tx.Commit(ctx)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", operation, err)
	}
	return id, nil
}

// createCode - обертка для операции вставки нового кода в БД.
func (db *DB) createCode(ctx context.Context, tx pgx.Tx, code models.RefCode) (int64, error) {
	err := tx.QueryRow(
		ctx,
		`INSERT INTO codes (code, owner, created, expired) VALUES ($1, $2, $3, $4) RETURNING id;`,
		code.Code,
		code.OwnerID,
		code.Created.Unix(),
		code.Expired.Unix(),
	).Scan(&code.ID)
	if err != nil {
		return 0, err
	}
	return code.ID, nil
}
