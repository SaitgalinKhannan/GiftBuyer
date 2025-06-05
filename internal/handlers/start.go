package handlers

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStartCommand() (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update telego.Update) error {
		// Отправка сообщения
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			"",
		))

		if err != nil {
			fmt.Printf("Error sending start command: %v\n", err)
		}

		return err
	}, th.CommandEqual("start")
}
