package model

type Payment struct {
	ID                      int    `json:"id" db:"id"`
	UserID                  int64  `json:"user_id" db:"user_id"`
	Currency                string `json:"currency" db:"currency"`
	Amount                  int    `json:"amount" db:"amount"`
	Payload                 string `json:"payload" db:"payload"`
	TelegramPaymentChargeID string `json:"telegram_payment_charge_id" json:"telegram_payment_charge_id" db:"telegram_payment_charge_id"`
	CreatedAt               string `json:"created_at" db:"created_at"`
}
