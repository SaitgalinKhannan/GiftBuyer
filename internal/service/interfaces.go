package service

import (
	"GiftBuyer/internal/model"
	"context"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type PaymentService interface {
	ValidatePreCheckout(ctx *th.Context, query *telego.PreCheckoutQuery) error
	ProcessSuccessfulPayment(ctx *th.Context, payment *telego.SuccessfulPayment, userID int64) error
}

type UserService interface {
	Create(ctx context.Context, user *telego.User) error
	GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error)
	GetByID(ctx context.Context, telegramID int64) (*model.User, error)
	UpdateBalance(ctx context.Context, telegramID int64, amount int) error
	GetBalance(ctx context.Context, telegramID int64) (float64, error)
	SetBalance(ctx context.Context, telegramID int64, balance float64) error
	DecrementBalance(ctx context.Context, telegramID int64, amount float64) error
	Update(ctx context.Context, user *model.User) error
	GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*model.User, error)
	CompareAndUpdate(ctx context.Context, user *model.User, telegramUser *telego.User) error
}

type Services struct {
	Payment PaymentService
	User    UserService
}
