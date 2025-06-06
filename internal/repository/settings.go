package repository

import (
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/model"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type settingsRepository struct {
	db *sqlx.DB
}

func NewSettingsRepository(db *database.DB) SettingsRepository {
	return &settingsRepository{db: db.DB}
}

func (r *settingsRepository) GetByUserID(ctx context.Context, userID int) (*model.UserSettings, error) {
	var settings model.UserSettings
	query := `
        SELECT id, user_id, auto_buy_enabled, price_limit_from, 
               price_limit_to, supply_limit, auto_buy_cycles, channels,
               created_at, updated_at
        FROM user_settings
        WHERE user_id = $1`

	err := r.db.GetContext(ctx, &settings, query, userID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Настройки не найдены
		}
		return nil, fmt.Errorf("failed to get settings: %w", err)
	}

	return &settings, nil
}

func (r *settingsRepository) Update(ctx context.Context, settings *model.UserSettings) error {
	query := `
        UPDATE user_settings
        SET auto_buy_enabled = $1,
            price_limit_from = $2,
            price_limit_to = $3,
            supply_limit = $4,
            auto_buy_cycles = $5,
            channels = $6,
            updated_at = CURRENT_TIMESTAMP
        WHERE user_id = $7`

	_, err := r.db.ExecContext(ctx, query,
		settings.AutoBuyEnabled,
		settings.PriceLimitFrom,
		settings.PriceLimitTo,
		settings.SupplyLimit,
		settings.AutoBuyCycles,
		settings.Channels,
		settings.UserID,
	)
	if err != nil {
		return fmt.Errorf("failed to update settings: %w", err)
	}

	return nil
}

func (r *settingsRepository) Create(ctx context.Context, userID int) error {
	query := `
        INSERT INTO user_settings (user_id)
        VALUES ($1)`

	_, err := r.db.ExecContext(ctx, query, userID)
	if err != nil {
		return fmt.Errorf("failed to create default settings: %w", err)
	}

	return nil
}
