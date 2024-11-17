// Пакет storage содержит методы для работы с базой данных PostgreSQL.
package storage

import (
	"context"
	"testing"
)

// testDB - адрес для подключения к тестовой БД.
const testDB = "postgres://tester:tester@localhost:5432/tests"

func Test_new(t *testing.T) {

	pool, err := new(testDB)
	if err != nil {
		t.Fatal(err)
	}
	defer pool.Close()

	err = pool.truncate()
	if err != nil {
		t.Fatal(err)
	}

	err = pool.populate()
	if err != nil {
		t.Fatal(err)
	}
}

// truncate очищает все таблицы в базе данных tests.
func (db *DB) truncate() error {
	_, err := db.pool.Exec(context.Background(), `TRUNCATE referrals, codes, users RESTART IDENTITY;`)
	if err != nil {
		return err
	}
	return nil
}

// populate заполняет таблицы тестовыми данными в базе данных tests.
func (db *DB) populate() error {
	_, err := db.pool.Exec(
		context.Background(),
		`INSERT INTO users (email, passwd, created) VALUES
		('bob@gmail.com', '12345678'::bytea, extract (epoch from now())),
		('bill@gmail.com', '12345678'::bytea, extract (epoch from now())),
		('jane@gmail.com', '12345678'::bytea, extract (epoch from now())),
		('jill@gmail.com', '12345678'::bytea, extract (epoch from now()));
		INSERT INTO codes (code, owner, created, expired, is_used) VALUES 
		('qwerty', 1, extract (epoch from now()), extract (epoch from (now() + interval '1 HOUR')), true),
		('asdfgh', 2, extract (epoch from now()), extract (epoch from now()), false),
		('zxcvbn', 3, extract (epoch from now()), extract (epoch from (now() + interval '1 HOUR')), false);
		INSERT INTO referrals (id, referrer_id, code_id) VALUES (4, 1, 1);`,
	)
	if err != nil {
		return err
	}
	return nil
}
