package model

import "time"

type UserSettings struct {
	ID              int       `json:"id" db:"id"`
	UserID          int       `json:"user_id" db:"user_id"`
	AutoBuyEnabled  bool      `json:"auto_buy_enabled" db:"auto_buy_enabled"`
	OnlyPremiumGift bool      `json:"only_premium_gift" db:"only_premium_gift"`
	PriceLimitFrom  *int      `json:"price_limit_from,omitempty" db:"price_limit_from"`
	PriceLimitTo    *int      `json:"price_limit_to,omitempty" db:"price_limit_to"`
	SupplyLimit     *int      `json:"supply_limit,omitempty" db:"supply_limit"`
	AutoBuyCycles   int       `json:"auto_buy_cycles,omitempty" db:"auto_buy_cycles"`
	Channels        *string   `json:"channels,omitempty" db:"channels"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
	UpdatedAt       time.Time `json:"updated_at" db:"updated_at"`
}
