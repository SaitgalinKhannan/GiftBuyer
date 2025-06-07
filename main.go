package main

import (
	. "GiftBuyer/app"
	"GiftBuyer/config"
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/handler"
	"GiftBuyer/internal/repository"
	"GiftBuyer/internal/scheduler"
	"GiftBuyer/internal/service"
	"GiftBuyer/logging"
	"context"
	"fmt"
	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"log"
	"os"
	"os/signal"
	"syscall"
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
		User:     repository.NewUserRepository(db),
		Gift:     repository.NewGiftRepository(db),
		Payment:  repository.NewPaymentRepository(db),
		Settings: repository.NewSettingsRepository(db),
	}
	services := &service.Services{
		Payment:  service.NewPaymentService(repos.Payment, repos.User),
		User:     service.NewUserService(repos.User),
		Settings: service.NewSettingsService(repos.Settings),
		Gift:     service.NewGiftService(repos.Gift, repos.User, repos.Settings, cfg),
	}

	bot, err := NewBot(cfg.BotToken, WithDefaultLogger(false, true))
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

	if botApp == nil {
		log.Fatal("botApp не инициализирован")
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

	_, err = bot.GetAvailableGifts(ctx)
	if err != nil {
		fmt.Println("Error:", err)
	}

	// 1. Создаём контекст с возможностью отмены
	ctxWithCancel, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 2. Инициализируем и запускаем watcher
	log.Println("Запуск GiftWatcher...")
	go scheduler.StartGiftWatcher(ctxWithCancel, botApp)

	// 3. Подписываемся на сигналы завершения
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	defer func() {
		signal.Stop(sigChan)
		cancel()
	}()

	// 4. Запускаем бота
	go func() {
		if err := bh.Start(); err != nil {
			log.Printf("Ошибка при запуске BotHandler: %v", err)
		}
	}()

	// 5. Ожидаем сигнал завершения
	<-sigChan
	log.Println("Получен сигнал завершения. Остановка...")

	// 6. Корректно завершаем работу
	if err := bh.Stop(); err != nil {
		log.Printf("Ошибка при остановке бота: %v", err)
	}
	cancel()
}
