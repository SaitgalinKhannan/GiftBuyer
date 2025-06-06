package database

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type DB struct {
	*sqlx.DB
}

func Connect(databaseURL string) (*DB, error) {
	db, err := sqlx.Connect("postgres", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Настройки пула соединений
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// Проверка соединения
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// Применение миграций
	if err := applyMigrations(databaseURL); err != nil {
		return nil, fmt.Errorf("failed to apply migrations: %w", err)
	}

	return &DB{db}, nil
}

func applyMigrations(databaseURL string) error {
	m, err := migrate.New(
		"file://migrations",
		databaseURL,
	)
	if err != nil {
		return fmt.Errorf("failed to initialize migrations: %w", err)
	}

	// Применить все миграции
	if err := m.Up(); err != nil {
		if err.Error() == "no change" {
			// Нет новых миграций — это нормально
			return nil
		}
		return fmt.Errorf("failed to run migrations: %w", err)
	}

	return nil
}

func (db *DB) Close() error {
	return db.DB.Close()
}
