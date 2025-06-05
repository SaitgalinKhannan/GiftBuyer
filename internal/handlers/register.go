package handlers

import (
	"GiftBuyer/app"
	th "github.com/mymmrac/telego/telegohandler"
	"log"
)

func RegisterHandlers(bh *th.BotHandler, a *app.App) {
	if bh == nil || a == nil {
		log.Fatal("Bot handler or app is nil")
	}

	bh.Handle(HandleStartCommand())
	bh.Handle(HandleGifts())
}
