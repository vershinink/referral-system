// Пакет storage содержит методы для работы с базой данных PostgreSQL.
package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"referral-rest-api/internal/config"
	"strings"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	ErrEmailExists    = errors.New("user with provided email already exists")
	ErrUserIdNotFound = errors.New("user with provided id not found")
	ErrEmailNotFound  = errors.New("user with provided email not found")
	ErrCodeExits      = errors.New("valid referral code already exists")
	ErrCodeNotFound   = errors.New("referral code with provided id not found")
	ErrCodeWasUsed    = errors.New("referral code was used")
	ErrCodeExpired    = errors.New("referral code expired")
	ErrUndefined      = errors.New("undefined database error")
)

const (
	// tmConn - таймаут на создание пула подключений к БД.
	tmConn time.Duration = time.Second * 10
)

// DB - пул подключений к БД.
type DB struct {
	pool *pgxpool.Pool
}

// New - обертка для конструктора пула подключений new.
func New(cfg *config.Config) *DB {
	// Формирование адреса для подключения к БД из полей конфига.
	elem := []string{
		"postgres://",
		cfg.StorageUser,
		":",
		cfg.StoragePasswd,
		"@",
		cfg.StoragePath,
		"/",
		cfg.StorageDB,
	}
	addr := strings.Join(elem, "")

	storage, err := new(addr)
	if err != nil {
		log.Fatalf("failed to init storage: %s", err.Error())
	}
	return storage
}

// new - конструктор пула подключений к БД.
func new(path string) (*DB, error) {
	const operation = "storage.new"

	ctx, cancel := context.WithTimeout(context.Background(), tmConn)
	defer cancel()

	// Создание пула подключений и проверка коннекта.
	pool, err := pgxpool.New(ctx, path)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}
	err = pool.Ping(ctx)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", operation, err)
	}

	return &DB{pool: pool}, nil
}

// Close - обертка для закрытия пула подключений.
func (db *DB) Close() {
	db.pool.Close()
}

// checkErr проверяет код ошибки из БД и возвращает соответствующую
// ему ошибку пакета storage для передачи в сервисный слой.
func checkErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case pgerrcode.UniqueViolation:
			return ErrEmailExists
		case pgerrcode.ForeignKeyViolation:
			return ErrUserIdNotFound
		default:
			return err
		}
	}
	return err
}
