package main

import (
	app "GiftBuyer/app"
	"GiftBuyer/config"
	"GiftBuyer/internal/handlers"
	"GiftBuyer/internal/repository"
	"GiftBuyer/pkg/database"
	"context"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
	"log"
	"os"
)

func main() {
	cfg := config.Load()
	ctx := context.Background()

	// Подключение к БД
	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	defer func(db *database.DB) {
		err := db.Close()
		if err != nil {
			log.Fatal("Failed to close database connection", err)
		}
	}(db)

	// Создание репозиториев
	repos := &repository.Repositories{
		User:    repository.NewUserRepository(db),
		Gift:    repository.NewGiftRepository(db),
		Payment: repository.NewPaymentRepository(db),
	}

	bot, err := NewBot(cfg.BotToken, WithDefaultDebugLogger())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	botApp := &app.App{
		DB:     db,
		Repos:  repos,
		Bot:    bot,
		Config: cfg,
	}

	botUser, err := bot.GetMe(context.Background())
	if err != nil {
		fmt.Println("Error:", err)
	}

	fmt.Printf("Bot user: %+v\n", botUser)

	// Get updates channel
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)

	// Create bot handler and specify from where to get updates
	bh, _ := th.NewBotHandler(bot, updates)
	// Stop handling updates
	defer func() { _ = bh.Stop() }()
	handlers.RegisterHandlers(bh, botApp)

	// Get gifts list
	/*bh.Handle(func(ctx *th.Context, update Update) error {
		availableGifts, _ := bot.GetAvailableGifts(ctx)

		// Print a list of gifts
		for _, gift := range availableGifts.Gifts {
			fmt.Printf("Gist ID: %s Price: %d ⭐️\n", gift.ID, gift.StarCount)
		}

		// Making a list of gifts
		var gifts strings.Builder
		gifts.WriteString("🎁 Доступные подарки:\n\n")

		for _, gift := range availableGifts.Gifts {
			gifts.WriteString(fmt.Sprintf("ID: %s | Цена: %d ⭐️\n", gift.ID, gift.StarCount))
		}

		// Send message with gifts list
		_, _ = ctx.Bot().SendMessage(ctx, tu.Message(
			tu.ID(update.Message.Chat.ID),
			gifts.String(),
		))

		return nil
	}, th.CommandEqual("gifts"))*/

	bh.Handle(func(ctx *th.Context, update Update) error {
		// Send message
		err := ctx.Bot().SendGift(ctx, &SendGiftParams{
			ChatID: ChatID{Username: "@soskiblya"},
			GiftID: "5170233102089322756",
			Text:   "🎁 Специальный подарок для тебя!",
		})

		if err != nil {
			fmt.Println("Ошибка отправки подарка:", err)
		}

		return nil
	}, th.CommandEqual("sendgift"))

	bh.Handle(func(ctx *th.Context, update Update) error {
		price := 15
		link, err := ctx.Bot().CreateInvoiceLink(ctx, &CreateInvoiceLinkParams{
			Title:         "Пополнение баланса",
			Description:   fmt.Sprintf("%d STARS", price),
			Payload:       fmt.Sprintf("Invoice ID:%d", update.Message.Chat.ID),
			ProviderToken: "",
			Currency:      "XTR",
			Prices: []LabeledPrice{
				{
					Label:  "STARS",
					Amount: price,
				},
			},
		})

		if err != nil {
			fmt.Println("Ошибка при создании инвойса:", err)
		}

		// Creating keyboard
		keyboard := tu.InlineKeyboard(
			tu.InlineKeyboardRow( // Row 1
				tu.InlineKeyboardButton(fmt.Sprintf("⭐️ Оплатить %d STARS", price)).WithURL(*link),
			),
			tu.InlineKeyboardRow(
				tu.InlineKeyboardButton("⭐️ Купить звезды дёшево").WithURL("https://split.tg/?ref=UQAEBELEbfrTtfyaT1ny28DCdQzSu34-mLv7gY-1czOlZMWL"),
			),
		)

		// Creating message
		msg := tu.Message(
			update.Message.Chat.ChatID(),
			"*Пополните бота по ссылке ниже\\.*",
		).WithReplyMarkup(keyboard).WithParseMode("MarkdownV2")

		_, err = ctx.Bot().SendMessage(ctx, msg)

		if err != nil {
			fmt.Println("Ошибка при отправке сообщения:", err)
		}

		return nil
	}, th.CommandEqual("top_up_balance"))

	// Обработка pre-checkout запросов
	bh.Handle(func(ctx *th.Context, update Update) error {
		if update.PreCheckoutQuery != nil {
			return handlePreCheckoutQuery(ctx, update.PreCheckoutQuery)
		}

		if update.Message != nil && update.Message.SuccessfulPayment != nil {
			return handleSuccessfulPayment(update.Message)
		}

		return nil
	}, th.AnyMessage())

	// Start handling updates
	_ = bh.Start()

}

// Обработка pre-checkout запроса
func handlePreCheckoutQuery(ctx *th.Context, query *PreCheckoutQuery) error {
	log.Printf("Pre-checkout запрос от пользователя %d, payload: %s, общая сумма: %d STARS",
		query.From.ID, query.InvoicePayload, query.TotalAmount)

	// Подтверждаем pre-checkout
	answer := &AnswerPreCheckoutQueryParams{
		PreCheckoutQueryID: query.ID,
		Ok:                 true,
		ErrorMessage:       "Ошибка, просим прощения, обратитесь позже.",
	}

	err := ctx.Bot().AnswerPreCheckoutQuery(ctx, answer)
	if err != nil {
		log.Printf("Ошибка ответа на pre-checkout: %v", err)
		return err
	}

	return nil
}

// Обработка успешного платежа
func handleSuccessfulPayment(message *Message) error {
	payment := message.SuccessfulPayment

	log.Printf("Успешный платеж от пользователя %d:", message.From.ID)
	log.Printf("- Валюта: %s", payment.Currency)
	log.Printf("- Сумма: %d", payment.TotalAmount)
	log.Printf("- Payload: %s", payment.InvoicePayload)
	log.Printf("- Telegram Payment Charge ID: %s", payment.TelegramPaymentChargeID)

	// Здесь добавляете свою бизнес-логику:
	// - Обновление базы данных
	// - Активация премиум-доступа
	// - Отправка подтверждения пользователю

	return nil
}
