package repository

import (
	"GiftBuyer/internal/model"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *model.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	GetByID(ctx context.Context, telegramID int64) (*model.User, error)
	GetAll(ctx context.Context) ([]*model.User, error)
	UpdateBalance(ctx context.Context, telegramID int64, amount int) error
	GetBalance(ctx context.Context, telegramID int64) (float64, error)
	SetBalance(ctx context.Context, telegramID int64, balance float64) error
	DecrementBalance(ctx context.Context, telegramID int64, amount float64) error
	Update(ctx context.Context, user *model.User) error
	GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*model.User, error)
}

type GiftRepository interface {
	Create(ctx context.Context, gift *model.Gift) error
	GetById(ctx context.Context, id string) (*model.Gift, error)
	GetAll(ctx context.Context) ([]*model.Gift, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *model.Payment) error
}

type SettingsRepository interface {
	GetByUserID(ctx context.Context, userID int) (*model.UserSettings, error)
	GetAll(ctx context.Context) ([]*model.UserSettings, error)
	Update(ctx context.Context, settings *model.UserSettings) error
	Create(ctx context.Context, userID int) error
}

type Repositories struct {
	User     UserRepository
	Gift     GiftRepository
	Payment  PaymentRepository
	Settings SettingsRepository
}
