package handler

import (
	. "GiftBuyer/internal/keyboard"
	"GiftBuyer/internal/service"
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStartCommand(userService service.UserService) (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update telego.Update) error {
		user, err := userService.GetByTelegramID(ctx, update.Message.Chat.ID)

		if err != nil {
			return err
		}

		if user == nil {
			createErr := userService.Create(ctx, update.Message.From)
			if createErr != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					tu.ID(update.Message.Chat.ID),
					"Не удалось зарегистрировать вас в системе, обратитесь в поддержку!",
				))
			}
		} else {
			err := userService.CompareAndUpdate(ctx, user, update.Message.From)
			if err != nil {
				fmt.Println(err)
			}
		}

		// Отправка сообщения
		_, err = ctx.Bot().SendMessage(ctx,
			tu.Message(
				tu.ID(update.Message.Chat.ID),
				"<b>Привет! Это удобный бот для покупки подарков в Telegram</b>\n\n"+
					"С ним ты можешь моментально и автоматически покупать новые подарки и обеспечить себе здоровый сон.",
			).WithParseMode("HTML").WithReplyMarkup(StartKeyboard()),
		)

		if err != nil {
			return fmt.Errorf("failed to send start message: %w", err) // ✅ Возвращаем ошибку
		}

		return nil
	}, th.CommandEqual("start")
}
