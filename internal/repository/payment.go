package repository

import (
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/model"
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type paymentRepository struct {
	db *sqlx.DB
}

func NewPaymentRepository(db *database.DB) PaymentRepository {
	return &paymentRepository{db: db.DB}
}

func (p *paymentRepository) Create(ctx context.Context, payment *model.Payment) error {
	fmt.Println(*payment)
	query := `
		INSERT INTO payments (user_id, currency, amount, payload, telegram_payment_charge_id)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at`

	err := p.db.QueryRowxContext(ctx, query,
		payment.UserID,
		payment.Currency,
		payment.Amount,
		payment.Payload,
		payment.TelegramPaymentChargeID,
	).Scan(&payment.ID, &payment.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create payment: %w", err)
	}

	return nil
}
