package repository

import (
	"GiftBuyer/internal/models"
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user *models.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*models.User, error)
	GetByID(ctx context.Context, id int) (*models.User, error)
	UpdateBalance(ctx context.Context, userID int, amount float64) error
	GetBalance(ctx context.Context, userID int) (float64, error)
	SetBalance(ctx context.Context, userID int, balance float64) error
	DecrementBalance(ctx context.Context, userID int, amount float64) error
	Update(ctx context.Context, user *models.User) error
	GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*models.User, error)
}

type GiftRepository interface {
	Create(ctx context.Context, gift *models.Gift) error
	GetById(ctx context.Context, id string) (*models.Gift, error)
	GetAll(ctx context.Context) ([]*models.Gift, error)
}

type PaymentRepository interface {
	Create(ctx context.Context, payment *models.Payment) error
}

type Repositories struct {
	User    UserRepository
	Gift    GiftRepository
	Payment PaymentRepository
}
