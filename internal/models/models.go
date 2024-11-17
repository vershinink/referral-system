// Пакет models содержит структуры объектов веб приложения.
package models

import "time"

// User - структура пользователя.
type User struct {
	ID       int64
	Email    string
	PassHash []byte
	RefCode  string
	Created  time.Time
	// AccessCode string
}

// RefCode - структура реферального кода.
type RefCode struct {
	ID      int64
	Code    string
	OwnerID int64
	Created time.Time
	Expired time.Time
	IsUsed  bool
}

// JWT - структура со строковым значением JWT токена.
type JWT struct {
	AccessToken string
}
