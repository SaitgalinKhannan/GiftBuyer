package service

import (
	"GiftBuyer/config"
	"GiftBuyer/internal/model"
	"GiftBuyer/internal/repository"
	"GiftBuyer/internal/utils"
	"GiftBuyer/logging"
	"context"
	"fmt"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"sync"
)

type giftService struct {
	giftRepo     repository.GiftRepository
	userRepo     repository.UserRepository
	settingsRepo repository.SettingsRepository
	conf         *config.Config
}

func NewGiftService(giftRepo repository.GiftRepository, userService repository.UserRepository, settingsRepo repository.SettingsRepository, conf *config.Config) GiftService {
	return &giftService{giftRepo: giftRepo, userRepo: userService, settingsRepo: settingsRepo, conf: conf}
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

func (g *giftService) BuyGiftForChannel(ctx context.Context, gift telego.Gift, channel string, bot *telego.Bot) error {
	err := bot.SendGift(ctx, &telego.SendGiftParams{
		ChatID: telego.ChatID{Username: channel},
		GiftID: gift.ID,
	})
	if err != nil {
		return fmt.Errorf("ошибка покупки подарка для канала %s: %w", channel, err)
	}
	return nil
}

func (g *giftService) BuyGiftForUser(ctx context.Context, gift telego.Gift, user *model.User, bot *telego.Bot) error {
	if user == nil {
		return fmt.Errorf("user is nil")
	}

	err := bot.SendGift(ctx, &telego.SendGiftParams{
		UserID: user.TelegramID,
		GiftID: gift.ID,
	})
	if err != nil {
		return fmt.Errorf("ошибка покупки подарка для пользователя %d: %w", user.TelegramID, err)
	}
	return nil
}

func (g *giftService) GetAutoBuyUsers(ctx context.Context, bot *telego.Bot) ([]*model.User, error) {
	// 1. Получаем всех пользователей
	users, err := g.userRepo.GetAll(ctx)
	if err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, fmt.Errorf("ошибка при получении пользователей: %w", err))
		return nil, fmt.Errorf("ошибка при получении пользователей: %w", err)
	}

	// 2. Получаем все настройки
	settingsList, sErr := g.settingsRepo.GetAll(ctx)
	if sErr != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, fmt.Errorf("ошибка при получении настроек: %w", sErr))
		return nil, fmt.Errorf("ошибка при получении настроек: %w", sErr)
	}

	// 3. Создаем мапу настроек по UserID для быстрого поиска
	settingsMap := make(map[int]*model.UserSettings)
	for _, s := range settingsList {
		settingsMap[s.UserID] = s
	}

	// 4. Фильтруем пользователей по AutoBuyEnabled
	var autoBuyUsers []*model.User
	for _, user := range users {
		if settings, exists := settingsMap[user.ID]; exists && settings.AutoBuyEnabled {
			autoBuyUsers = append(autoBuyUsers, user)
		}
	}

	return autoBuyUsers, nil
}

func (g *giftService) ProcessAutoBuy(ctx context.Context, newGifts []telego.Gift, bot *telego.Bot) {
	users, err := g.GetAutoBuyUsers(ctx, bot)
	if err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		fmt.Println(err)
		return
	}

	var wg sync.WaitGroup
	for _, user := range users {
		wg.Add(1)
		go func(u *model.User) {
			defer wg.Done()
			g.processUserAutoBuy(ctx, u, newGifts, bot)
		}(user)
	}
	wg.Wait()
}

func (g *giftService) processUserAutoBuy(ctx context.Context, user *model.User, newGifts []telego.Gift, bot *telego.Bot) {
	settings, err := g.settingsRepo.GetByUserID(ctx, user.ID)
	if err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		return
	}

	if !settings.AutoBuyEnabled {
		return
	}

	channels := utils.StringToChannels(settings.Channels)

	maxCycles := 1000
	if settings.AutoBuyCycles > 0 {
		maxCycles = settings.AutoBuyCycles
	}

	for i := 0; i < maxCycles; i++ {
		minPrice := 1_000_000
		for _, gift := range newGifts {
			if gift.StarCount < minPrice {
				minPrice = gift.StarCount
			}
		}

		if user.Balance < minPrice {
			break
		}

		for _, gift := range newGifts {
			if !isGiftMatchSettings(gift, settings) {
				continue
			}

			if user.Balance < gift.StarCount {
				continue
			}

			if len(channels) == 0 {
				// Покупаем подарок
				if err := g.BuyGiftForUser(ctx, gift, user, bot); err != nil {
					logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
					continue
				}
			} else {
				targetChannel := channels[i%len(channels)]
				if err := g.BuyGiftForChannel(ctx, gift, targetChannel, bot); err != nil {
					logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
					continue
				}
			}

			// Обновляем баланс и остатки
			user.Balance -= gift.StarCount
			fmt.Printf("Баланс юзера ID %d = %d\n", user.ID, user.Balance)
			_ = g.UpdateUserBalance(ctx, user)
		}
	}
}

func isGiftMatchSettings(gift telego.Gift, settings *model.UserSettings) bool {
	if settings.PriceLimitFrom != nil && gift.StarCount < *settings.PriceLimitFrom {
		return false
	}
	if settings.PriceLimitTo != nil && gift.StarCount > *settings.PriceLimitTo {
		return false
	}
	if settings.SupplyLimit != nil && gift.TotalCount > *settings.SupplyLimit {
		return false
	}
	return true
}

func (g *giftService) UpdateUserBalance(ctx context.Context, user *model.User) error {
	return g.userRepo.UpdateBalance(ctx, user.TelegramID, user.Balance)
}

/*func (g *giftService) UpdateGiftSupply(ctx context.Context, gift *telego.Gift) error {
	return g.giftRepo.UpdateRemainingCount(ctx, gift.ID, gift.RemainingCount)
}*/

func (g *giftService) CheckAndProcessNewGifts(ctx context.Context, bot *telego.Bot) error {
	// 1. Получить текущие подарки из Telegram
	telegramGifts, err := g.GetAvailableGifts(ctx, bot)
	if err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		log.Printf("Получить текущие подарки из Telegram: %v", err)
		return err
	}

	// 2. Получить сохраненные подарки из БД
	savedGifts, err := g.GetAll(ctx)
	if err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		log.Printf("Получить сохраненные подарки из БД: %v", err)
		return err
	}

	// 3. Найти новые подарки
	newGifts := g.CompareGiftLists(savedGifts, telegramGifts)
	if len(newGifts) == 0 {
		return nil
	}

	// 4. Уведомить пользователей
	if err := g.NotifyUsers(ctx, newGifts, bot); err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		log.Printf("Ошибка уведомления пользователей: %v", err)
	}

	// 5. Сохранить новые подарки в БД
	if err := g.SaveNewGifts(ctx, newGifts); err != nil {
		logging.SendLogErrorToTelegram(ctx, bot, g.conf.LogChatId, err)
		log.Printf("Ошибка сохранения новых подарков: %v", err)
	}

	g.ProcessAutoBuy(ctx, newGifts, bot)
	// 6. Купить подарки для пользователей
	/*for _, gift := range newGifts {
		if err := g.BuyGiftForChannel(ctx, gift, someUserID); err != nil {
			log.Printf("Ошибка покупки подарка %s: %v", gift.ID, err)
		}
	}*/

	// 1. берем пользователей из бд
	// 2. запускаем 5 или сколько-то горутин, параллельно надо обрабатывать пользователей, которые поставили автоматическую покупку
	// 3. берем настройки покупки подарков пользователя
	// 4. если его баланс достаточен и новые подарки по цене и саплаю подходят,
	// то покупаем каждый подарок по 1 разу указанное кол-во циклов (если не указан, просто покупаем пока балансе хватает) в настройках, пока баланса хватает.
	// Если у него указаны каналы, то отправляем подарок туда, если 3 канала и 3 вида подарков,
	// то один канал на один вид подарков, если каналов меньше, то придется с начала списка каналов начать и другие подарки тоже туда скупать
	// 5. дальше идем по пользователям и повторяем предыдущие шаги
	// additional надо проверять, что новые подарки все еще есть, т.е. не раскуплены

	return nil
}
