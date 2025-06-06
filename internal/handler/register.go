package handler

import (
	"GiftBuyer/app"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"runtime/debug"
)

func RegisterHandlers(bh *th.BotHandler, a *app.App) {
	if bh == nil || a == nil {
		log.Fatal("Bot handler or app is nil")
	}

	paymentHandler := NewPaymentHandler(a.Services.Payment, a.StateStorage)

	// Middleware для логирования ошибок
	bh.Use(func(ctx *th.Context, update telego.Update) error {
		err := ctx.Next(update)
		if err != nil {
			//log.Printf("Global error handler: %v\n", err)
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
				log.Printf("Recovered from panic: %v\nStack: %s", r, debug.Stack())
			}
		}()
		return ctx.Next(update)
	})

	/*bh.Handle(func(ctx *th.Context, update telego.Update) error {
		err := ctx.Bot().SendGift(ctx, &telego.SendGiftParams{
			UserID: 538321015,
			GiftID: "5170233102089322756",
			Text:   "🎁 Специальный подарок для тебя!",
		})

		if err != nil {
			fmt.Println("Ошибка отправки подарка:", err)
			return err
		}

		return nil
	}, th.CommandEqual("gift"))*/

	bh.Handle(paymentHandler.HandlePayment())
	bh.Handle(HandleStartCommand(a.Services.User))
	bh.Handle(paymentHandler.HandleTopUpBalanceCallback())
	bh.Handle(HandleCallback(a))
	bh.Handle(HandleGifts())
	bh.Handle(StateHandler(a))
	bh.Handle(HandleSettingsCallback(a))
	bh.Handle(HandlePriceLimitUpdateCallback(a))
	bh.Handle(HandleSupplyLimitUpdateCallback(a))
	bh.Handle(HandleAutoBuyCyclesUpdateCallback(a))
	bh.Handle(HandleChannelSettingsCallback(a))
}
