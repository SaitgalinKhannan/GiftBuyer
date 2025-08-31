package repository

import (
	. "GiftBuyer/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	GetByID(ctx context.Context, telegramID int64) (*User, error)
	GetAll(ctx context.Context) ([]*User, error)
	UpdateBalance(ctx context.Context, telegramID int64, amount int) error
	GetBalance(ctx context.Context, telegramID int64) (float64, error)
	SetBalance(ctx context.Context, telegramID int64, balance float64) error
	DecrementBalance(ctx context.Context, telegramID int64, amount float64) error
	Update(ctx context.Context, user *User) error
	GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*User, error)
}

type GiftRepository interface {
	Create(ctx context.Context, gift *Gift) error
	GetById(ctx context.Context, id string) (*Gift, error)
	GetAll(ctx context.Context) ([]*Gift, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *Payment) error
}

type SettingsRepository interface {
	GetByUserID(ctx context.Context, userID int) (*UserSettings, error)
	GetAll(ctx context.Context) ([]*UserSettings, error)
	Update(ctx context.Context, settings *UserSettings) error
	Create(ctx context.Context, userID int) error
}

// AccountRepository определяет контракт для работы с аккаунтами
type AccountRepository interface {
	// Save сохраняет или обновляет аккаунт
	Save(account *Account) error

	// GetAll возвращает все активные аккаунты
	GetAll() ([]Account, error)

	// GetByID возвращает аккаунт по ID
	GetByID(id int64) (*Account, error)

	// Delete удаляет аккаунт по ID
	Delete(id int64) error

	// SetActive устанавливает статус активности аккаунта
	SetActive(id int64, active bool) error

	// UpdateUserInfo обновляет информацию о пользователе
	UpdateUserInfo(id int64, username, firstName, lastName, phone string) error

	// Close закрывает соединение с БД
	Close() error
}

type Repositories struct {
	User     UserRepository
	Gift     GiftRepository
	Payment  PaymentRepository
	Settings SettingsRepository
	Account  AccountRepository
}
