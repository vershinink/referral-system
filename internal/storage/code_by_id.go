package storage

import (
	"context"
	"errors"
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

// CodeByID получает структуру реферального кода из БД по id
// его пользователя.
func (db *DB) CodeByID(ctx context.Context, userID int64) (models.RefCode, error) {
	const operation = "storage.CodeByID"

	var code models.RefCode
	var created, expired int64
	err := db.pool.QueryRow(
		ctx,
		`SELECT id, code, owner, created, expired, is_used FROM codes 
		 WHERE owner = $1 ORDER BY created DESC LIMIT 1;`,
		userID,
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
			return models.RefCode{}, fmt.Errorf("%s: %w", operation, ErrCodeNotFound)
		}
		return models.RefCode{}, fmt.Errorf("%s: %w", operation, err)
	}

	code.Created = time.Unix(created, 0)
	code.Expired = time.Unix(expired, 0)
	return code, nil
}
