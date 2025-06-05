package handlers

import (
	"fmt"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"strings"
)

func HandleGifts() (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update telego.Update) error {
		// Получение подарков
		gifts, err := ctx.Bot().GetAvailableGifts(ctx)
		if err != nil {
			return fmt.Errorf("failed to get gifts: %w", err)
		}

		// Формирование сообщения
		var giftList strings.Builder
		giftList.WriteString("🎁 Доступные подарки:\n\n")
		for _, gift := range gifts.Gifts {
			giftList.WriteString(fmt.Sprintf("ID: %s | Цена: %d ⭐️\n", gift.ID, gift.StarCount))
		}

		// Отправка сообщения
		_, err = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			giftList.String()),
		)

		return err
	}, th.CommandEqual("gifts")
}
