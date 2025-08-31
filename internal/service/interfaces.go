package service

import (
	. "GiftBuyer/internal/model"
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
	GetByTelegramID(ctx context.Context, telegramID int64) (*User, error)
	GetByID(ctx context.Context, telegramID int64) (*User, error)
	UpdateBalance(ctx context.Context, telegramID int64, amount int) error
	GetBalance(ctx context.Context, telegramID int64) (float64, error)
	SetBalance(ctx context.Context, telegramID int64, balance float64) error
	DecrementBalance(ctx context.Context, telegramID int64, amount float64) error
	Update(ctx context.Context, user *User) error
	GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*User, error)
	CompareAndUpdate(ctx context.Context, user *User, telegramUser *telego.User) error
}

type SettingsService interface {
	GetByUserID(ctx context.Context, userID int) (*UserSettings, error)
	Update(ctx context.Context, settings *UserSettings) error
	Create(ctx context.Context, userID int) error
}

type GiftService interface {
	Create(ctx context.Context, gift *telego.Gift) error
	GetById(ctx context.Context, id string) (*Gift, error)
	GetAll(ctx context.Context) ([]*Gift, error)
	SaveNewGifts(ctx context.Context, newGifts []telego.Gift) error
	CompareGiftLists(gifts []*Gift, telegramGifts []telego.Gift) []telego.Gift
	GetAvailableGifts(ctx context.Context, bot *telego.Bot) ([]telego.Gift, error)
	NotifyUsers(ctx context.Context, newGifts []telego.Gift, bot *telego.Bot) error
	BuyGiftForChannel(ctx context.Context, gift telego.Gift, channel string, bot *telego.Bot) error
	BuyGiftForUser(ctx context.Context, gift telego.Gift, user *User, bot *telego.Bot) error
	CheckAndProcessNewGifts(ctx context.Context, bot *telego.Bot) error
}

// AccountService определяет контракт для бизнес-логики работы с аккаунтами
type AccountService interface {
	// Create создает или обновляет аккаунт
	Create(account *Account) error

	// GetAll возвращает все активные аккаунты
	GetAll() ([]*Account, error)

	// GetByID возвращает аккаунт по ID
	GetByID(id int64) (*Account, error)

	// Delete удаляет аккаунт по ID
	Delete(id int64) error

	// SetActive устанавливает статус активности аккаунта
	SetActive(id int64, active bool) error

	// UpdateUserInfo обновляет информацию о пользователе: username, имя, фамилию, телефон
	UpdateUserInfo(id int64, username, firstName, lastName, phone string) error

	Close() error
}

type Services struct {
	Payment  PaymentService
	User     UserService
	Settings SettingsService
	Gift     GiftService
	Account  AccountService
}
