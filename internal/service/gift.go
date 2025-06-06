package service

import (
	"GiftBuyer/config"
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
)

type giftService struct {
	giftRepo repository.GiftRepository
	conf     *config.Config
}

func NewGiftService(giftRepo repository.GiftRepository, conf *config.Config) GiftService {
	return &giftService{giftRepo: giftRepo, conf: conf}
}

func (g *giftService) Create(ctx context.Context, gift *telego.Gift) error {
	newGift := &model.Gift{
		ID:               gift.ID,
		StarCount:        gift.StarCount,
		UpgradeStarCount: gift.UpgradeStarCount,
		TotalCount:       gift.TotalCount,
		RemainingCount:   gift.RemainingCount,
	}
	return g.giftRepo.Create(ctx, newGift)
}

func (g *giftService) GetById(ctx context.Context, id string) (*model.Gift, error) {
	return g.giftRepo.GetById(ctx, id)
}

func (g *giftService) GetAll(ctx context.Context) ([]*model.Gift, error) {
	return g.giftRepo.GetAll(ctx)
}

func (g *giftService) SaveNewGifts(ctx context.Context, newGifts []telego.Gift) error {
	for _, gift := range newGifts {
		modelGift := &model.Gift{
			ID:               gift.ID,
			StarCount:        gift.StarCount,
			UpgradeStarCount: gift.UpgradeStarCount,
			TotalCount:       gift.TotalCount,
			RemainingCount:   gift.RemainingCount,
		}
		if err := g.giftRepo.Create(ctx, modelGift); err != nil {
			log.Printf("Ошибка сохранения подарка %s: %v", gift.ID, err)
		}
	}
	return nil
}

func (g *giftService) CompareGiftLists(gifts []*model.Gift, telegramGifts []telego.Gift) []telego.Gift {
	// Создаем мапу существующих подарков по ID
	existing := make(map[string]struct{})
	for _, gift := range gifts {
		existing[gift.ID] = struct{}{}
	}

	// Собираем новые подарки
	var newGifts []telego.Gift
	for _, tg := range telegramGifts {
		if _, exists := existing[tg.ID]; !exists {
			newGifts = append(newGifts, tg)
		}
	}

	return newGifts
}

func (g *giftService) GetAvailableGifts(ctx context.Context, bot *telego.Bot) ([]telego.Gift, error) {
	// Предположим, что у вас есть экземпляр бота в App
	gifts, err := bot.GetAvailableGifts(ctx)
	if err != nil {
		return nil, fmt.Errorf("ошибка получения подарков: %w", err)
	}
	return gifts.Gifts, nil
}

func (g *giftService) NotifyUsers(ctx context.Context, newGifts []telego.Gift, bot *telego.Bot) error {
	if len(newGifts) == 0 {
		return nil
	}
	if g.conf == nil {
		return fmt.Errorf("config is nil")
	}

	text := g.buildNotificationMessage(newGifts)
	_, err := bot.SendMessage(ctx, tu.Message(
		telego.ChatID{ID: g.conf.NotificationChannelId},
		text,
	))
	if err != nil {
		fmt.Println(err)
	}

	for _, gift := range newGifts {
		_, err := bot.SendMessage(ctx, tu.Message(
			telego.ChatID{ID: g.conf.NotificationChannelId},
			gift.Sticker.Emoji,
		))

		if err != nil {
			fmt.Printf("can`t send sticker %s to notification chat: %v\n", gift.ID, err)
		}
	}

	return nil
}

func (g *giftService) buildNotificationMessage(gifts []telego.Gift) string {
	// Реализуйте формирование текста уведомления
	return fmt.Sprintf("Новые подарки доступны! Количество: %d", len(gifts))
}

func (g *giftService) BuyGift(ctx *th.Context, gift telego.Gift, userID int64) error {
	err := ctx.Bot().SendGift(ctx, &telego.SendGiftParams{
		UserID: userID,
		GiftID: gift.ID,
	})
	if err != nil {
		return fmt.Errorf("ошибка покупки подарка: %w", err)
	}
	return nil
}

func (g *giftService) CheckAndProcessNewGifts(ctx context.Context, bot *telego.Bot) error {
	// 1. Получить текущие подарки из Telegram
	telegramGifts, err := g.GetAvailableGifts(ctx, bot)
	if err != nil {
		return err
	}

	// 2. Получить сохраненные подарки из БД
	savedGifts, err := g.GetAll(ctx)
	if err != nil {
		return err
	}

	// 3. Найти новые подарки
	newGifts := g.CompareGiftLists(savedGifts, telegramGifts)
	if len(newGifts) == 0 {
		return nil
	}

	// 4. Уведомить пользователей
	if err := g.NotifyUsers(ctx, newGifts, bot); err != nil {
		log.Printf("Ошибка уведомления пользователей: %v", err)
	}

	// 5. Сохранить новые подарки в БД
	if err := g.SaveNewGifts(ctx, newGifts); err != nil {
		log.Printf("Ошибка сохранения новых подарков: %v", err)
	}

	// 6. Опционально: купить подарки для пользователей
	// for _, gift := range newGifts {
	//     if err := g.BuyGift(ctx, gift, someUserID); err != nil {
	//         log.Printf("Ошибка покупки подарка %s: %v", gift.ID, err)
	//     }
	// }

	return nil
}
