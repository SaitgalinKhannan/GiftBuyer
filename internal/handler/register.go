package handler

import (
	"GiftBuyer/app"
	"GiftBuyer/logging"
	"log"
	"runtime/debug"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func RegisterHandlers(bh *th.BotHandler, a *app.App, updates <-chan telego.Update) {
	if bh == nil || a == nil {
		log.Fatal("Bot handler or app is nil")
	}

	paymentHandler := NewPaymentHandler(a.Services.Payment, a)

	// Middleware для логирования ошибок
	bh.Use(func(ctx *th.Context, update telego.Update) error {
		err := ctx.Next(update)
		if err != nil {
			//log.Printf("Global error handler: %v\n", err)
			logging.SendLogErrorToTelegram(ctx, ctx.Bot(), a.Config.LogChatId, err)
			if update.Message != nil {
				_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
					update.Message.Chat.ChatID(),
					"Ошибка, свяжитесь с поддержкой!",
				))
			}
		}
		return err
	})

	bh.Use(func(ctx *th.Context, update telego.Update) error {
		defer func() {
			if r := recover(); r != nil {
				logging.SendLogMessageToTelegram(ctx, ctx.Bot(), a.Config.LogChatId, "Recovered from panic in register.go")
				log.Printf("Recovered from panic: %v\nStack: %s", r, debug.Stack())
			}
		}()
		return ctx.Next(update)
	})

	bh.Handle(paymentHandler.HandlePayment())
	bh.Handle(HandleStartCommand(a))
	bh.Handle(paymentHandler.HandleTopUpBalanceCallback())
	bh.Handle(HandleCallback(a))
	bh.Handle(HandleGifts())
	bh.Handle(StateHandler(a))
	bh.Handle(HandleSettingsCallback(a, updates))
	bh.Handle(HandlePriceLimitUpdateCallback(a))
	bh.Handle(HandleSupplyLimitUpdateCallback(a))
	bh.Handle(HandleAutoBuyCyclesUpdateCallback(a))
	bh.Handle(HandleChannelSettingsCallback(a))
}
