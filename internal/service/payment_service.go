package service

import (
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type paymentService struct {
	paymentRepo repository.PaymentRepository
	userRepo    repository.UserRepository
}

func NewPaymentService(paymentRepo repository.PaymentRepository, userRepo repository.UserRepository) PaymentService {
	return &paymentService{paymentRepo: paymentRepo, userRepo: userRepo}
}

func (s *paymentService) ValidatePreCheckout(ctx *th.Context, query *telego.PreCheckoutQuery) error {
	// Проверка пользователя
	user, err := s.userRepo.GetByID(ctx, query.From.ID)
	if user == nil || err != nil {
		return fmt.Errorf("user not found ID %d: %w", query.From.ID, err)
	}

	return nil
}

func (s *paymentService) ProcessSuccessfulPayment(ctx *th.Context, payment *telego.SuccessfulPayment, userID int64) error {
	// Проверка пользователя
	user, err := s.userRepo.GetByID(ctx, userID)
	if user == nil || err != nil {
		return fmt.Errorf("user not found ID %d: %w", userID, err)
	}

	err = s.userRepo.UpdateBalance(ctx, userID, payment.TotalAmount)
	if err != nil {
		return fmt.Errorf("unable to update user balance with ID %d : %w", userID, err)
	}

	err = s.paymentRepo.Create(ctx, &model.Payment{
		UserID:                  userID,
		Currency:                "XTR",
		Amount:                  payment.TotalAmount,
		Payload:                 payment.InvoicePayload,
		TelegramPaymentChargeID: payment.TelegramPaymentChargeID,
	})
	if err != nil {
		fmt.Printf("unable to save user payment with ID %d : %s", userID, err)
	}

	return nil
}
