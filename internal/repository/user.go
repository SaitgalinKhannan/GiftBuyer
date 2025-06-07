package repository

import (
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type userRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *database.DB) UserRepository {
	return &userRepository{db: db.DB}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO users (telegram_id, username, first_name, last_name, balance, is_active)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query,
		user.TelegramID,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Balance,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, telegram_id, username, first_name, last_name, balance, created_at, is_active
		FROM users
		WHERE telegram_id = $1`

	err := r.db.GetContext(ctx, &user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("failed to get user by telegram_id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetByID(ctx context.Context, telegramID int64) (*model.User, error) {
	var user model.User
	query := `
		SELECT id, telegram_id, username, first_name, last_name, balance, created_at, is_active
		FROM users
		WHERE telegram_id = $1`

	err := r.db.GetContext(ctx, &user, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) GetAll(ctx context.Context) ([]*model.User, error) {
	var users []*model.User
	query := `
		SELECT id, telegram_id, username, first_name, last_name, balance, created_at, is_active
		FROM users
		WHERE is_active = true
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &users, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get users: %w", err)
	}

	return users, nil
}

func (r *userRepository) UpdateBalance(ctx context.Context, telegramID int64, amount int) error {
	query := `
		UPDATE users
		SET balance = balance + $1
		WHERE telegram_id = $2`

	result, err := r.db.ExecContext(ctx, query, amount, telegramID)
	if err != nil {
		return fmt.Errorf("failed to update balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) GetBalance(ctx context.Context, telegramID int64) (float64, error) {
	var balance float64
	query := `SELECT balance FROM users WHERE telegram_id = $1`

	err := r.db.GetContext(ctx, &balance, query, telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

func (r *userRepository) Update(ctx context.Context, user *model.User) error {
	query := `
		UPDATE users
		SET username = $1, first_name = $2, last_name = $3, balance = $4, is_active = $5
		WHERE telegram_id = $6`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.FirstName,
		user.LastName,
		user.Balance,
		user.IsActive,
		user.TelegramID,
	)
	if err != nil {
		return fmt.Errorf("failed to update user: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) SetBalance(ctx context.Context, telegramID int64, balance float64) error {
	query := `
		UPDATE users
		SET balance = $1
		WHERE telegram_id = $2`

	result, err := r.db.ExecContext(ctx, query, balance, telegramID)
	if err != nil {
		return fmt.Errorf("failed to set balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	return nil
}

func (r *userRepository) DecrementBalance(ctx context.Context, telegramID int64, amount float64) error {
	// Используем транзакцию для безопасного списания
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if rbErr := tx.Rollback(); rbErr != nil {
			// Сохраняем ошибку Rollback, если основная ошибка отсутствует
			if err == nil {
				err = fmt.Errorf("rollback error: %w", rbErr)
			} else {
				// Или добавляем к основной ошибке
				err = fmt.Errorf("%v; rollback error: %w", err, rbErr)
			}
		}
	}()

	// Проверяем текущий баланс
	var currentBalance float64
	err = tx.GetContext(ctx, &currentBalance, "SELECT balance FROM users WHERE telegram_id = $1", telegramID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("user not found")
		}
		return fmt.Errorf("failed to get current balance: %w", err)
	}

	// Проверяем, достаточно ли средств
	if currentBalance < amount {
		return fmt.Errorf("insufficient balance: current %.2f, required %.2f", currentBalance, amount)
	}

	// Списываем средства
	result, err := tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE telegram_id = $2", amount, telegramID)
	if err != nil {
		return fmt.Errorf("failed to decrement balance: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("user not found")
	}

	// Коммитим транзакцию
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetUsersWithMinBalance Метод для получения пользователей с балансом больше указанной суммы
func (r *userRepository) GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*model.User, error) {
	var users []*model.User
	query := `
		SELECT id, telegram_id, username, first_name, last_name, balance, created_at, is_active
		FROM users
		WHERE balance >= $1 AND is_active = true
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &users, query, minBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with min balance: %w", err)
	}

	return users, nil
}
