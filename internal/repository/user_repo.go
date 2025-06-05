package repository

import (
	"GiftBuyer/internal/models"
	"GiftBuyer/pkg/database"
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

func (r *userRepository) Create(ctx context.Context, user *models.User) error {
	query := `
		INSERT INTO users (telegram_id, username, balance, is_active)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`

	err := r.db.QueryRowxContext(ctx, query,
		user.TelegramID,
		user.Username,
		user.Balance,
		user.IsActive,
	).Scan(&user.ID, &user.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

func (r *userRepository) GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, telegram_id, username, balance, created_at, is_active
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

func (r *userRepository) UpdateBalance(ctx context.Context, userID int, amount float64) error {
	query := `
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, amount, userID)
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

func (r *userRepository) GetBalance(ctx context.Context, userID int) (float64, error) {
	var balance float64
	query := `SELECT balance FROM users WHERE id = $1`

	err := r.db.GetContext(ctx, &balance, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, fmt.Errorf("user not found")
		}
		return 0, fmt.Errorf("failed to get balance: %w", err)
	}

	return balance, nil
}

// GetByID Дополнительные полезные методы
func (r *userRepository) GetByID(ctx context.Context, id int) (*models.User, error) {
	var user models.User
	query := `
		SELECT id, telegram_id, username, balance, created_at, is_active
		FROM users
		WHERE id = $1`

	err := r.db.GetContext(ctx, &user, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get user by id: %w", err)
	}

	return &user, nil
}

func (r *userRepository) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET username = $1, balance = $2, is_active = $3
		WHERE id = $4`

	result, err := r.db.ExecContext(ctx, query,
		user.Username,
		user.Balance,
		user.IsActive,
		user.ID,
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

func (r *userRepository) SetBalance(ctx context.Context, userID int, balance float64) error {
	query := `
		UPDATE users
		SET balance = $1
		WHERE id = $2`

	result, err := r.db.ExecContext(ctx, query, balance, userID)
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

func (r *userRepository) DecrementBalance(ctx context.Context, userID int, amount float64) error {
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
	err = tx.GetContext(ctx, &currentBalance, "SELECT balance FROM users WHERE id = $1", userID)
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
	result, err := tx.ExecContext(ctx, "UPDATE users SET balance = balance - $1 WHERE id = $2", amount, userID)
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
func (r *userRepository) GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*models.User, error) {
	var users []*models.User
	query := `
		SELECT id, telegram_id, username, balance, created_at, is_active
		FROM users
		WHERE balance >= $1 AND is_active = true
		ORDER BY created_at DESC`

	err := r.db.SelectContext(ctx, &users, query, minBalance)
	if err != nil {
		return nil, fmt.Errorf("failed to get users with min balance: %w", err)
	}

	return users, nil
}
