package main

import (
	. "GiftBuyer/app"
	"GiftBuyer/config"
	"GiftBuyer/internal/client"
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/handler"
	"GiftBuyer/internal/repository"
	"GiftBuyer/internal/scheduler"
	"GiftBuyer/internal/service"
	"GiftBuyer/logging"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	. "github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"go.uber.org/zap"
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
		Account:  repository.NewAccountRepository(db),
	}
	services := &service.Services{
		Payment:  service.NewPaymentService(repos.Payment, repos.User),
		User:     service.NewUserService(repos.User),
		Settings: service.NewSettingsService(repos.Settings),
		Gift:     service.NewGiftService(repos.Gift, repos.User, repos.Settings, cfg),
		Account:  service.NewAccountService(repos.Account),
	}

	// создание менеджера юзер ботов
	accountManager := client.NewAccountManager(services.Account)
	defer accountManager.Stop()

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
		AccountManager: accountManager,
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
	handler.RegisterHandlers(bh, botApp, updates)

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
		if startErr := bh.Start(); startErr != nil {
			log.Printf("Ошибка при запуске BotHandler: %v", startErr)
		}
	}()

	// 5. Загружаем существующие аккаунты из БД
	go func() {
		log.Println("Loading existing accounts from database...")
		if loadAccountsErr := botApp.AccountManager.LoadAccounts(bot, updates); loadAccountsErr != nil {
			log.Printf("Warning: Failed to load accounts: %v\n", loadAccountsErr)
		}

		// Показываем список аккаунтов из БД
		accounts, _ := botApp.AccountManager.GetAllAccounts()
		if accounts == nil || len(accounts) == 0 {
			log.Println("No accounts found in database. Please add accounts first.")
			log.Println("Example: manager.AddNewAccount(accountID, apiID, apiHash)")
			return
		}

		log.Printf("Found %d accounts in database:\n", len(accounts))
		for _, acc := range accounts {
			status := "active"
			if !acc.IsActive {
				status = "inactive"
			}
			log.Printf("  - ID: %d, Username: @%s, Name: %s %s, Status: %s\n", acc.ID, acc.Username, acc.FirstName, acc.LastName, status)
		}

		// Запускаем все загруженные аккаунты
		if accounts != nil && len(accounts) > 0 {
			if startAllErr := accountManager.StartAll(); startAllErr != nil {
				log.Fatalf("Failed to start accounts: %v", startAllErr)
			}

			log.Println("All accounts are running. Press Ctrl+C to stop.")

			// Пример использования клиента для отправки сообщения
			go func() {
				time.Sleep(5 * time.Second) // Ждем пока все запустится

				// Получаем первый аккаунт из списка
				if tgClient, getClientErr := accountManager.GetClient(accounts[0].ID); getClientErr == nil {
					// Здесь можно использовать tgClient.API для работы с Telegram API
					tgClient.Logger.Info("Client is ready for use", zap.String("username", tgClient.AccountInfo.Username))
				}
			}()

			// Ждем завершения
			select {}
		} else {
			log.Println("No active accounts to start.")
		}
	}()

	// 6. Ожидаем сигнал завершения
	<-sigChan
	log.Println("Получен сигнал завершения. Остановка...")

	// 7. Корректно завершаем работу
	if stopErr := bh.Stop(); stopErr != nil {
		log.Printf("Ошибка при остановке бота: %v", stopErr)
	}
	accountManager.Stop()

	cancel()
}
