package service

import (
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"context"
	"fmt"
	"log"
)

type settingsService struct {
	settingsRepo repository.SettingsRepository
}

func NewSettingsService(settingsRepo repository.SettingsRepository) SettingsService {
	return &settingsService{settingsRepo: settingsRepo}
}

func (s *settingsService) GetByUserID(ctx context.Context, userID int) (*model.UserSettings, error) {
	settings, err := s.settingsRepo.GetByUserID(ctx, userID)
	if err != nil || settings == nil {
		// Если настройки не найдены, создаем новые
		err = s.Create(ctx, userID)
		if err != nil {
			log.Printf("Ошибка создания настроек: %v", err)
			return nil, fmt.Errorf("ошибка создания настроек")
		}

		// Получаем созданные настройки
		settings, err = s.settingsRepo.GetByUserID(ctx, userID)
		if err != nil || settings == nil {
			log.Printf("Не удалось получить настройки после создания: %v", err)
			return nil, fmt.Errorf("настройки не созданы")
		}
	}

	return settings, nil
}

func (s *settingsService) Update(ctx context.Context, settings *model.UserSettings) error {
	return s.settingsRepo.Update(ctx, settings)
}

func (s *settingsService) Create(ctx context.Context, userID int) error {
	return s.settingsRepo.Create(ctx, userID)
}
