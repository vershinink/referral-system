package storage

import (
	"context"
	"errors"
	"fmt"
	"referral-rest-api/internal/models"
	"time"

	"github.com/jackc/pgx/v5"
)

// CodeByEmail получает структуру текущего неиспользованного реферального
// кода по email создавшего его пользователя из БД.
func (db *DB) CodeByEmail(ctx context.Context, email string) (models.RefCode, error) {
	const operation = "storage.CodeByEmail"

	var code models.RefCode
	var created, expired int64
	err := db.pool.QueryRow(
		ctx,
		`SELECT codes.id, codes.code, codes.owner, codes.created, codes.expired, codes.is_used 
		FROM codes INNER JOIN users ON codes.owner = users.id 
		WHERE users.email = $1 AND codes.is_used = false
		ORDER BY codes.created DESC LIMIT 1;`,
		email,
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
