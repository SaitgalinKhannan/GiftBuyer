package service

import (
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"context"
)

type settingsService struct {
	settingsRepo repository.SettingsRepository
}

func NewSettingsService(settingsRepo repository.SettingsRepository) SettingsService {
	return &settingsService{settingsRepo: settingsRepo}
}

func (s *settingsService) GetByUserID(ctx context.Context, userID int) (*model.UserSettings, error) {
	return s.settingsRepo.GetByUserID(ctx, userID)
}

func (s *settingsService) Update(ctx context.Context, settings *model.UserSettings) error {
	return s.settingsRepo.Update(ctx, settings)
}

func (s *settingsService) Create(ctx context.Context, userID int) error {
	return s.settingsRepo.Create(ctx, userID)
}
