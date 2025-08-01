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

type SettingsService interface {
	GetByUserID(ctx context.Context, userID int) (*model.UserSettings, error)
	Update(ctx context.Context, settings *model.UserSettings) error
	Create(ctx context.Context, userID int) error
}

type GiftService interface {
	Create(ctx context.Context, gift *telego.Gift) error
	GetById(ctx context.Context, id string) (*model.Gift, error)
	GetAll(ctx context.Context) ([]*model.Gift, error)
	SaveNewGifts(ctx context.Context, newGifts []telego.Gift) error
	CompareGiftLists(gifts []*model.Gift, telegramGifts []telego.Gift) []telego.Gift
	GetAvailableGifts(ctx context.Context, bot *telego.Bot) ([]telego.Gift, error)
	NotifyUsers(ctx context.Context, newGifts []telego.Gift, bot *telego.Bot) error
	BuyGiftForChannel(ctx context.Context, gift telego.Gift, channel string, bot *telego.Bot) error
	BuyGiftForUser(ctx context.Context, gift telego.Gift, user *model.User, bot *telego.Bot) error
	CheckAndProcessNewGifts(ctx context.Context, bot *telego.Bot) error
}

type Services struct {
	Payment  PaymentService
	User     UserService
	Settings SettingsService
	Gift     GiftService
}
