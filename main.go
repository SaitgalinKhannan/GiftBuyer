package main

import (
	. "GiftBuyer/app"
	"GiftBuyer/config"
	"GiftBuyer/internal/handler"
	"GiftBuyer/internal/repository"
	"GiftBuyer/internal/service"
	"GiftBuyer/pkg/database"
	"GiftBuyer/pkg/logging"
	"context"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"log"
	"os"
)

func main() {
	logFile, err := logging.InitLogger()
	if err != nil {
		log.Fatal("Ошибка при инициализации логгера:", err)
	}
	defer logFile.Close() // Закрываем файл при завершении

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
	services := &service.Services{
		Payment: service.NewPaymentService(repos.Payment, repos.User),
		User:    service.NewUserService(repos.User),
	}

	bot, err := NewBot(cfg.BotToken, WithDefaultLogger(true, true))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	_ = bot.SetMyCommands(ctx, &SetMyCommandsParams{
		Commands: []BotCommand{
			{
				Command:     "start",
				Description: "Главное меню",
			},
		},
	})

	botApp := &App{
		DB:       db,
		Repos:    repos,
		Services: services,
		Bot:      bot,
		Config:   cfg,
		StateStorage: &StateStorage{
			States: make(map[int64]State),
		},
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
	handler.RegisterHandlers(bh, botApp)

	// Start handling updates
	_ = bh.Start()
}
