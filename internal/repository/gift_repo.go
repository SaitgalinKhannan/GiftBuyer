package repository

import (
	"GiftBuyer/internal/model"
	"GiftBuyer/pkg/database"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type giftRepository struct {
	db *sqlx.DB
}

func NewGiftRepository(db *database.DB) GiftRepository {
	return &giftRepository{db: db.DB}
}

func (g giftRepository) Create(ctx context.Context, gift *model.Gift) error {
	query := `
		INSERT INTO gifts (id, star_count, upgrade_star_count, total_count, remaining_count)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING created_at`

	err := g.db.QueryRowxContext(ctx, query,
		gift.ID,
		gift.StarCount,
		gift.UpgradeStarCount,
		gift.TotalCount,
		gift.RemainingCount,
	).Scan(&gift.CreatedAt)

	if err != nil {
		return fmt.Errorf("failed to create gift: %w", err)
	}

	return nil
}

func (g giftRepository) GetById(ctx context.Context, id string) (*model.Gift, error) {
	var gift model.Gift
	query := `
		SELECT id, star_count, upgrade_star_count, total_count, remaining_count, created_at
		FROM gifts
		WHERE id = $1`

	err := g.db.GetContext(ctx, &gift, query, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Пользователь не найден
		}
		return nil, fmt.Errorf("failed to get gift by id: %w", err)
	}

	return &gift, nil
}

func (g giftRepository) GetAll(ctx context.Context) ([]*model.Gift, error) {
	var gifts []*model.Gift
	query := `
		SELECT id, star_count, upgrade_star_count, total_count, remaining_count, created_at
		FROM gifts`

	err := g.db.SelectContext(ctx, &gifts, query)
	if err != nil {
		return nil, fmt.Errorf("failed to get all gifts: %w", err)
	}

	return gifts, nil
}
