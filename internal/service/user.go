package service

import (
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"context"
	"github.com/mymmrac/telego"
)

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{userRepo: userRepo}
}

func (u *userService) Create(ctx context.Context, user *telego.User) error {
	return u.userRepo.Create(ctx, &model.User{
		TelegramID: user.ID,
		Username:   user.Username,
		FirstName:  user.FirstName,
		LastName:   user.LastName,
		Balance:    0,
		IsActive:   true,
	})
}

func (u *userService) GetByTelegramID(ctx context.Context, telegramID int64) (*model.User, error) {
	return u.userRepo.GetByTelegramID(ctx, telegramID)
}

func (u *userService) GetByID(ctx context.Context, telegramID int64) (*model.User, error) {
	return u.userRepo.GetByID(ctx, telegramID)
}

func (u *userService) UpdateBalance(ctx context.Context, telegramID int64, amount int) error {
	return u.userRepo.UpdateBalance(ctx, telegramID, amount)
}

func (u *userService) GetBalance(ctx context.Context, telegramID int64) (float64, error) {
	return u.userRepo.GetBalance(ctx, telegramID)
}

func (u *userService) SetBalance(ctx context.Context, telegramID int64, balance float64) error {
	return u.userRepo.SetBalance(ctx, telegramID, balance)
}

func (u *userService) DecrementBalance(ctx context.Context, telegramID int64, amount float64) error {
	return u.userRepo.DecrementBalance(ctx, telegramID, amount)
}

func (u *userService) Update(ctx context.Context, user *model.User) error {
	return u.userRepo.Update(ctx, user)
}

func (u *userService) GetUsersWithMinBalance(ctx context.Context, minBalance float64) ([]*model.User, error) {
	return u.userRepo.GetUsersWithMinBalance(ctx, minBalance)
}

func (u *userService) CompareAndUpdate(ctx context.Context, user *model.User, telegramUser *telego.User) error {
	isUsernameSame := user.Username == telegramUser.Username
	isFirstNameSame := user.FirstName == telegramUser.FirstName
	isLastNameSame := user.LastName == telegramUser.LastName

	if !(isUsernameSame && isFirstNameSame && isLastNameSame) {
		user.Username = telegramUser.Username
		user.FirstName = telegramUser.FirstName
		user.LastName = telegramUser.LastName
		return u.Update(ctx, user)
	}

	return nil
}
