package handler

import (
	"GiftBuyer/app"
	. "GiftBuyer/internal/keyboard"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func HandleCallback(storage *app.StateStorage) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		data := update.CallbackQuery.Data
		switch data {
		case "back_to_start":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:    update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Text: "<b>Привет! Это удобный бот для покупки подарков в Telegram</b>\n\n" +
					"С ним ты можешь моментально и автоматически покупать новые подарки и обеспечить себе здоровый сон.",
				ReplyMarkup: StartKeyboard(),
				ParseMode:   "HTML",
			})

			// Сбрасываем состояние
			storage.ClearState(update.CallbackQuery.From.ID)

			if err != nil {
				return err
			}

		case "profile":
			_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

			_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
				ChatID:    update.CallbackQuery.Message.GetChat().ChatID(),
				MessageID: update.CallbackQuery.Message.GetMessageID(),
				Text: "<b>👤 Ваш профиль:</b>\n\n" +
					fmt.Sprintf("<b>⭐️ Баланс звёзд в боте:</b> %d\n\n", 0) +
					fmt.Sprintf("<b>⭐️ Подарков куплено:</b> %d на сумму %d ⭐️", 0, 0),
				ReplyMarkup: ProfileKeyboard(),
				ParseMode:   "HTML",
			})

			if err != nil {
				return err
			}
		}

		return nil
	}, th.AnyCallbackQueryWithMessage()
}
