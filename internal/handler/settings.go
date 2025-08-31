package handler

import (
	"GiftBuyer/app"
	. "GiftBuyer/internal/keyboard"
	"GiftBuyer/internal/utils"
	"context"
	"fmt"
	"log"
	"slices"
	"strconv"
	"strings"

	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

// Предикат для обработки нажатий кнопок настройки
func isSettingsCallbackQuery(_ context.Context, update Update) bool {
	// Проверяем, что CallbackQuery и Message не nil
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		return false
	}

	// Список допустимых значений data
	allowedData := map[string]struct{}{
		"login":             {},
		"settings":          {},
		"toggle_auto_buy":   {},
		"only_premium_gift": {},
		"price_from":        {},
		"price_to":          {},
		"supply_limit":      {},
		"auto_buy_cycles":   {},
		"channels":          {},
		"add_channel":       {},
	}

	// Проверяем, есть ли update.CallbackQuery.Data в allowedData
	_, ok := allowedData[update.CallbackQuery.Data]
	return ok
}

func HandleSettingsCallback(a *app.App, updates <-chan Update) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ChatID()
		msgID := update.CallbackQuery.Message.GetMessageID()
		data := update.CallbackQuery.Data

		switch data {
		case "login":
			if err := a.AccountManager.AddNewAccount(userID, a.Config.ApiID, a.Config.ApiHash, a.Bot, updates); err != nil {
				log.Printf("Failed to add new account: %v\n", err)
			}
		case "settings":
			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем текущие настройки
			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        utils.FormatAutoBuySettings(settings),
				ReplyMarkup: SettingsKeyboard(settings),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "toggle_auto_buy":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем текущие настройки
			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Меняем состояние автопокупки
			settings.AutoBuyEnabled = !settings.AutoBuyEnabled
			err = a.Services.Settings.Update(ctx, settings)

			if err != nil {
				log.Printf("Ошибка обновления настроек: %v", err)
				return fmt.Errorf("ошибка обновления настроек")
			}

			// Обновляем сообщение с новой клавиатурой
			_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        utils.FormatAutoBuySettings(settings),
				ReplyMarkup: SettingsKeyboard(settings),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "only_premium_gift":
			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем текущие настройки
			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Меняем состояние автопокупки
			settings.OnlyPremiumGift = !settings.OnlyPremiumGift
			err = a.Services.Settings.Update(ctx, settings)

			if err != nil {
				log.Printf("Ошибка обновления настроек: %v", err)
				return fmt.Errorf("ошибка обновления настроек")
			}

			fmt.Printf("AutoBuyEnabled: %t\n", settings.AutoBuyEnabled)
			fmt.Printf("OnlyPremiumGift: %t\n", settings.OnlyPremiumGift)

			// Обновляем сообщение с новой клавиатурой
			_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        utils.FormatAutoBuySettings(settings),
				ReplyMarkup: SettingsKeyboard(settings),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "price_from":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Обновляем сообщение с новой клавиатурой
			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        "<b>Выбери новый минимум цены для автопокупки:\n(бот не отправит подарок дешевле установленного лимита)</b>",
				ReplyMarkup: SetPriceFromKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "price_to":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Обновляем сообщение с новой клавиатурой
			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        "<b>Выбери новый максимум цены для автопокупки:\n(бот не отправит подарок дороже установленного лимита)</b>",
				ReplyMarkup: SetPriceToKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "supply_limit":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Обновляем сообщение с новой клавиатурой
			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   update.CallbackQuery.Message.GetMessageID(),
				Text:        "<b>Выбери новый лимит саплая для автопокупки:\n(бот не отправит подарок, если их выпущено больше установленного лимита)</b>",
				ReplyMarkup: SetSupplyLimitKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "auto_buy_cycles":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Обновляем сообщение с новой клавиатурой
			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:    update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Text: "<b>Выбери новое количество циклов автопокупки:\n(то, сколько раз бот купит новый подарок, например: " +
					"выходит 3 подарка, циклов 2 - бот купит каждый по 2 раза = 6 подарков)</b>",
				ReplyMarkup: SetAutoBuyCyclesKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "channels":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})

			a.StateStorage.ClearState(userID)

			// Получаем пользователя
			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем настройки
			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			channels := utils.StringToChannels(settings.Channels)

			// Обновляем сообщение с новой клавиатурой
			_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID:   msgID,
				Text:        "<b>Ваши каналы, на которые бот будет отправлять подарки:</b>",
				ReplyMarkup: ChannelsKeyboard(channels),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

		case "add_channel":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{
				CallbackQueryID: update.CallbackQuery.ID,
			})

			// Обновляем сообщение
			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      chatID,
				MessageID:   msgID,
				Text:        "Введите <b>username</b> канала в формате <b>@username</b>",
				ReplyMarkup: GoToChannelsKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}

			a.StateStorage.SetState(userID, app.StateWaitingChannelUsername)

		default:
			return ctx.Next(update)
		}

		return nil
	}, isSettingsCallbackQuery
}

func HandlePriceLimitUpdateCallback(a *app.App) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		// Подтверждаем callback
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ChatID()
		msgID := update.CallbackQuery.Message.GetMessageID()

		// Получаем пользователя
		user, err := a.Services.User.GetByTelegramID(ctx, userID)
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				chatID, "Пользователь не найден. Обратитесь в поддержку.",
			))
			return err
		}

		// Получаем настройки
		settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
		if err != nil {
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
			return err
		}
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		// Извлекаем тип и значение из callback_data
		callbackData := update.CallbackQuery.Data

		var (
			priceType string
			rawValue  string
		)

		if strings.HasPrefix(callbackData, "set_price_from=") {
			priceType = "from"
			rawValue = strings.TrimPrefix(callbackData, "set_price_from=")
		} else if strings.HasPrefix(callbackData, "set_price_to=") {
			priceType = "to"
			rawValue = strings.TrimPrefix(callbackData, "set_price_to=")
		} else {
			return nil
		}

		// Обрабатываем значение
		var value *int
		switch rawValue {
		case "nil":
			value = nil
		default:
			num, parseIntErr := strconv.Atoi(rawValue)
			if parseIntErr != nil {
				// Некорректное значение — показываем ошибку
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Некорректное значение лимита. Попробуйте ещё раз.",
				))
				return nil
			}
			value = &num
		}

		// Обновляем настройки
		switch priceType {
		case "from":
			settings.PriceLimitFrom = value
		case "to":
			settings.PriceLimitTo = value
		}

		// Сохраняем изменения
		err = a.Services.Settings.Update(ctx, settings)
		if err != nil {
			log.Printf("Ошибка обновления лимитов цены: %v", err)
			return err
		}

		// Обновляем сообщение
		_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   msgID,
			Text:        utils.FormatAutoBuySettings(settings),
			ReplyMarkup: SettingsKeyboard(settings),
			ParseMode:   "HTML",
		})

		if err != nil {
			return err
		}

		return nil
	}, th.CallbackDataPrefix("set_price_")
}

func HandleSupplyLimitUpdateCallback(a *app.App) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		// Подтверждаем callback
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ChatID()
		msgID := update.CallbackQuery.Message.GetMessageID()

		// Получаем пользователя
		user, err := a.Services.User.GetByTelegramID(ctx, userID)
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				chatID, "Пользователь не найден. Обратитесь в поддержку.",
			))
			return err
		}

		// Получаем настройки
		settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
		if err != nil {
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
			return err
		}
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		// Извлекаем значение из callback_data
		callbackData := update.CallbackQuery.Data
		rawValue := strings.TrimPrefix(callbackData, "set_supply_limit=")

		// Обрабатываем значение
		var value *int
		switch rawValue {
		case "nil":
			value = nil
		default:
			num, err := strconv.Atoi(rawValue)
			if err != nil {
				// Некорректное значение — показываем ошибку
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Некорректное значение лимита. Попробуйте ещё раз.",
				))
				return nil
			}
			value = &num
		}

		// Обновляем настройки
		settings.SupplyLimit = value

		// Сохраняем изменения
		err = a.Services.Settings.Update(ctx, settings)
		if err != nil {
			log.Printf("Ошибка обновления лимита кол-ва подарков: %v", err)
			return err
		}

		// Обновляем сообщение
		_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   msgID,
			Text:        utils.FormatAutoBuySettings(settings),
			ReplyMarkup: SettingsKeyboard(settings),
			ParseMode:   "HTML",
		})
		if err != nil {
			return err
		}

		return nil
	}, th.CallbackDataPrefix("set_supply_limit=")
}

func HandleAutoBuyCyclesUpdateCallback(a *app.App) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		// Подтверждаем callback
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{
			CallbackQueryID: update.CallbackQuery.ID,
		})

		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ChatID()
		msgID := update.CallbackQuery.Message.GetMessageID()

		// Получаем пользователя
		user, err := a.Services.User.GetByTelegramID(ctx, userID)
		if err != nil {
			_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
				chatID, "Пользователь не найден. Обратитесь в поддержку.",
			))
			return err
		}

		// Получаем настройки
		settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
		if err != nil {
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
			return err
		}
		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		// Извлекаем значение из callback_data
		callbackData := update.CallbackQuery.Data
		rawValue := strings.TrimPrefix(callbackData, "set_auto_buy_cycles=")

		// Обрабатываем значение
		var value int

		switch rawValue {
		case "infinite":
			value = -1 // Используем -1 для обозначения бесконечных циклов
		default:
			num, parseIntErr := strconv.Atoi(rawValue)
			if parseIntErr != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Некорректное значение циклов. Попробуйте ещё раз.",
				))
				return nil
			}
			value = num
		}

		// Обновляем настройки
		settings.AutoBuyCycles = value

		// Сохраняем изменения
		err = a.Services.Settings.Update(ctx, settings)
		if err != nil {
			log.Printf("Ошибка обновления циклов покупки: %v", err)
			return err
		}

		// Обновляем сообщение
		_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
			ChatID:      chatID,
			MessageID:   msgID,
			Text:        utils.FormatAutoBuySettings(settings),
			ReplyMarkup: SettingsKeyboard(settings),
			ParseMode:   "HTML",
		})
		if err != nil {
			return err
		}

		return nil
	}, th.CallbackDataPrefix("set_auto_buy_cycles=")
}

// Предикат для обработки нажатий кнопок настройки канала
func isChannelSettingsCallbackQuery(_ context.Context, update Update) bool {
	// Проверяем, что CallbackQuery и Message не nil
	if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
		return false
	}

	// Список допустимых значений data
	allowedData := map[string]struct{}{
		"channel":        {},
		"delete_channel": {},
	}

	// Проверяем, есть ли update.CallbackQuery.Data в allowedData
	parts := strings.Split(update.CallbackQuery.Data, "=")

	if len(parts) > 0 {
		str := parts[0]
		_, ok := allowedData[str]
		return ok
	}

	return false
}

func HandleChannelSettingsCallback(a *app.App) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		userID := update.CallbackQuery.From.ID
		chatID := update.CallbackQuery.Message.GetChat().ChatID()
		msgID := update.CallbackQuery.Message.GetMessageID()
		parts := strings.Split(update.CallbackQuery.Data, "=")
		var callbackData string
		if len(parts) > 0 {
			callbackData = parts[0]
		}
		switch callbackData {
		case "channel":
			// Получаем пользователя
			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем настройки
			_, err = a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Извлекаем значение из callback_data
			data := update.CallbackQuery.Data
			username := strings.TrimPrefix(data, "channel=")

			// Обновляем сообщение
			_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:      chatID,
				MessageID:   msgID,
				Text:        fmt.Sprintf("Канал <b>%s</b>", username),
				ReplyMarkup: ChannelSettingsKeyboard(username),
				ParseMode:   "HTML",
			})
			if err != nil {
				return err
			}
		case "delete_channel":
			// Получаем пользователя
			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					chatID, "Пользователь не найден. Обратитесь в поддержку.",
				))
				return err
			}

			// Получаем настройки
			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil {
				_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID, Text: err.Error(), ShowAlert: true})
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			// Извлекаем значение из callback_data
			data := update.CallbackQuery.Data
			username := strings.TrimPrefix(data, "delete_channel=")
			channels := utils.StringToChannels(settings.Channels)

			// Находим индекс элемента
			index := slices.Index(channels, username)

			// Если элемент найден, удаляем его
			if index != -1 {
				channels = slices.Delete(channels, index, index+1)
			}

			chStr := utils.ChannelsToString(channels)
			settings.Channels = &chStr

			err = a.Services.Settings.Update(ctx, settings)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Не удалось удалить канал, обратитесь в поддержку!"))
				a.StateStorage.ClearState(userID)
			} else {
				// Обновляем сообщение
				_, err = ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
					ChatID:      chatID,
					MessageID:   msgID,
					Text:        fmt.Sprintf("Канал <b>%s</b> удалён!", username),
					ReplyMarkup: GoToChannelsKeyboard(),
					ParseMode:   "HTML",
				})
			}
			if err != nil {
				return err
			}
		default:
			return ctx.Next(update)
		}

		return nil
	}, isChannelSettingsCallbackQuery
}
