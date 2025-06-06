package handler

import (
	. "GiftBuyer/app"
	. "GiftBuyer/internal/keyboard"
	"GiftBuyer/internal/utils"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strconv"
	"strings"
)

func StateHandler(a *App) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		userID := update.Message.From.ID
		chatID := update.Message.Chat.ChatID()

		// Получаем текущее состояние
		state := a.StateStorage.GetState(userID)

		switch state {
		case StateWaitingTopUpAmount:
			if update.Message == nil || update.Message.Text == "" {
				return nil
			}

			amount, err := strconv.ParseUint(update.Message.Text, 10, 64)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Введите число!"))
				return nil
			}

			if amount <= 0 {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Введите положительное число!"))
				return nil
			}

			// Отправка инвойса
			link, err := ctx.Bot().CreateInvoiceLink(ctx, &CreateInvoiceLinkParams{
				Title:         "Пополнение баланса",
				Description:   fmt.Sprintf("Пополнение на %d звёзд", amount),
				Payload:       fmt.Sprintf("topup_%d", userID),
				ProviderToken: "",
				Currency:      "XTR",
				Prices: []LabeledPrice{
					{Label: "STARS", Amount: int(amount)},
				},
			})

			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Ошибка при создании инвойса."))
				return err
			}

			// Creating message
			msg := tu.Message(
				update.Message.Chat.ChatID(),
				"<b>Пополните счёт в боте по ссылке ниже.</b>",
			).WithReplyMarkup(TopUpBalanceKeyboard(int(amount), *link)).WithParseMode("HTML")

			_, err = ctx.Bot().SendMessage(ctx, msg)

			if err != nil {
				return err
			}

			// Сбрасываем состояние
			a.StateStorage.ClearState(userID)

		case StateWaitingChannelUsername:
			if update.Message == nil || update.Message.Text == "" {
				return nil
			}

			user, err := a.Services.User.GetByTelegramID(ctx, userID)
			if err != nil || user == nil {
				a.StateStorage.ClearState(userID)
				return fmt.Errorf("пользователь не найден")
			}

			settings, err := a.Services.Settings.GetByUserID(ctx, user.ID)
			if err != nil || settings == nil {
				a.StateStorage.ClearState(userID)
				return fmt.Errorf("настройки не найдены")
			}

			username := update.Message.Text
			if !strings.HasPrefix(username, "@") {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Введите в таком формате: <b>@username</b>!").WithParseMode("HTML"))
				return nil
			}

			newChannel := username
			channels := utils.StringToChannels(settings.Channels)

			if utils.Contains(channels, username) {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Такой канал уже добавлен!").WithReplyMarkup(GoToChannelsKeyboard()))
				a.StateStorage.ClearState(userID)
				return nil
			}

			channels = append(channels, newChannel)
			chStr := utils.ChannelsToString(channels)
			settings.Channels = &chStr

			err = a.Services.Settings.Update(ctx, settings)
			if err != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(chatID, "Не удалось добавить канал, обратитесь в поддержку!"))
				a.StateStorage.ClearState(userID)
			} else {
				msg := tu.Message(
					update.Message.Chat.ChatID(),
					fmt.Sprintf("Канал <b>%s</b> добавлен!", username),
				).WithReplyMarkup(GoToChannelsKeyboard()).WithParseMode("HTML")

				_, err = ctx.Bot().SendMessage(ctx, msg)
			}

			if err != nil {
				return err
			}

			a.StateStorage.ClearState(userID)

		default:
			return ctx.Next(update)
		}

		return nil
	}, th.AnyMessage()
}
