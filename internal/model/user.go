package model

import "time"

type User struct {
	ID         int       `json:"id" db:"id"`
	TelegramID int64     `json:"telegram_id" db:"telegram_id"`
	Username   string    `json:"username" db:"username"`
	FirstName  string    `json:"first_name" db:"first_name"`
	LastName   string    `json:"last_name" db:"last_name"`
	Balance    int       `json:"balance" db:"balance"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	IsActive   bool      `json:"is_active" db:"is_active"`
}
