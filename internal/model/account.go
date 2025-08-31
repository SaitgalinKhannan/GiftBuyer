package model

import "time"

type Account struct {
	ID        int64     `db:"id"`
	ApiID     int       `db:"api_id"`
	ApiHash   string    `db:"api_hash"`
	Phone     string    `db:"phone"`
	Username  string    `db:"username"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	IsActive  bool      `db:"is_active"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
