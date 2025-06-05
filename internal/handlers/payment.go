package handlers

import (
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func (a *App) HandlePayments() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		if update.PreCheckoutQuery != nil {
			return a.handlePreCheckoutQuery(ctx, update.PreCheckoutQuery)
		}
		if update.Message != nil && update.Message.SuccessfulPayment != nil {
			return a.handleSuccessfulPayment(update.Message)
		}
		return nil
	}
}
