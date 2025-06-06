package handler

import (
	. "GiftBuyer/app"
	. "GiftBuyer/internal/keyboard"
	"GiftBuyer/internal/service"
	"context"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
)

type PaymentHandler struct {
	paymentService service.PaymentService
	stateStorage   *StateStorage
}

func NewPaymentHandler(paymentService service.PaymentService, storage *StateStorage) *PaymentHandler {
	return &PaymentHandler{paymentService: paymentService, stateStorage: storage}
}

// Предикат для определения платежных обновлений
func isPaymentUpdate(_ context.Context, update Update) bool {
	return update.PreCheckoutQuery != nil ||
		(update.Message != nil && update.Message.SuccessfulPayment != nil)
}

// HandlePayment Handle HandlePayments Обработка pre-checkout запросов и успешной оплаты
func (h *PaymentHandler) HandlePayment() (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.PreCheckoutQuery != nil {
			return h.handlePreCheckoutQuery(ctx, update.PreCheckoutQuery)
		}
		if update.Message != nil && update.Message.SuccessfulPayment != nil {
			return h.handleSuccessfulPayment(ctx, update.Message)
		}

		return nil
	}, isPaymentUpdate
}

func (h *PaymentHandler) HandleTopUpBalanceCallback() (th.Handler, th.Predicate) {
	return func(ctx *th.Context, update Update) error {
		if update.CallbackQuery == nil || update.CallbackQuery.Message == nil {
			return nil
		}

		_ = ctx.Bot().AnswerCallbackQuery(ctx, &AnswerCallbackQueryParams{CallbackQueryID: update.CallbackQuery.ID})

		// Устанавливаем состояние ожидания суммы
		h.stateStorage.SetState(update.CallbackQuery.From.ID, StateWaitingTopUpAmount)

		_, err := ctx.Bot().EditMessageText(ctx, &EditMessageTextParams{
			ChatID:      update.CallbackQuery.Message.GetChat().ChatID(),
			MessageID:   update.CallbackQuery.Message.GetMessageID(),
			Text:        "<b>Отправьте В ЧАТ сумму пополнения в звездах:</b>\n\nЕсли у вас нет звезд, купите их на ваш аккаунт по ссылке",
			ReplyMarkup: BuyStarsKeyboard(),
			ParseMode:   "HTML",
		})
		if err != nil {
			return err
		}
		return nil
	}, th.CallbackDataEqual("top_up_balance")
}

// Обработка pre-checkout запроса
func (h *PaymentHandler) handlePreCheckoutQuery(ctx *th.Context, query *PreCheckoutQuery) error {
	log.Printf("Pre-checkout запрос от пользователя %d, payload: %s, общая сумма: %d STARS",
		query.From.ID, query.InvoicePayload, query.TotalAmount)

	err := h.paymentService.ValidatePreCheckout(ctx, query)
	answerResult := true

	if err != nil {
		answerResult = false
	}

	// Подтверждаем pre-checkout
	answer := &AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: query.ID,
		Ok:                 answerResult,
		ErrorMessage:       "Ошибка, попробуйте позже.",
	}

	err = ctx.Bot().AnswerPreCheckoutQuery(ctx, answer)
	if err != nil {
		return fmt.Errorf("ошибка ответа на pre-checkout пользователя %d: %v", query.From.ID, err)
	}

	return nil
}

// Обработка успешного платежа
func (h *PaymentHandler) handleSuccessfulPayment(ctx *th.Context, message *Message) error {
	payment := message.SuccessfulPayment

	if payment == nil {
		return nil
	}

	log.Printf(
		"[Платёж] Пользователь: %d | Валюта: %s | Сумма: %d | Payload: %s | ChargeID: %s",
		message.From.ID,
		payment.Currency,
		payment.TotalAmount,
		payment.InvoicePayload,
		payment.TelegramPaymentChargeID,
	)

	// Вызов сервиса
	err := h.paymentService.ProcessSuccessfulPayment(ctx, payment, message.From.ID)
	if err != nil {
		log.Printf("Failed to process payment: %v", err)
		return err
	}

	// Отправка подтверждения пользователю
	_, err = ctx.Bot().SendMessage(ctx, &SendMessageParams{
		ChatID:      tu.ID(message.Chat.ID),
		Text:        fmt.Sprintf("✅ Платеж на сумму %d успешно обработан!", payment.TotalAmount),
		ReplyMarkup: GoMainKeyboard(),
	})
	return err
}
