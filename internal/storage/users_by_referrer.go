package storage

import (
	"context"
	"fmt"
	"referral-rest-api/internal/models"
	"time"
)

// UsersByReferrerID получает массив структур пользователей-рефералов
// по id их реферера.
func (db *DB) UsersByReferrerID(ctx context.Context, refID int64) ([]models.User, error) {
	const operation = "storage.UsersByReferrerID"

	rows, err := db.pool.Query(
		ctx,
		`SELECT users.id, users.email, users.passwd, users.created FROM users
		 INNER JOIN referrals ON users.id = referrals.id WHERE referrals.referrer_id = $1
		 ORDER BY users.created DESC;`,
		refID,
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	var users []models.User
	for rows.Next() {
		var user models.User
		var created int64
		err := rows.Scan(
			&user.ID,
			&user.Email,
			&user.PassHash,
			&created,
		)
		if err != nil {
			return nil, fmt.Errorf("%s: %w", operation, err)
		}
		user.Created = time.Unix(created, 0)
		users = append(users, user)
	}
	if rows.Err() != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%s: %w", operation, ErrUserIdNotFound)
	}

	return users, nil
}
