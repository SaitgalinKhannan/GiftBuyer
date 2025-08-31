package client

import (
	. "GiftBuyer/internal/model"
	"GiftBuyer/internal/service"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	pebbledb "github.com/cockroachdb/pebble"
	boltstor "github.com/gotd/contrib/bbolt"
	"github.com/gotd/contrib/middleware/floodwait"
	"github.com/gotd/contrib/middleware/ratelimit"
	"github.com/gotd/contrib/pebble"
	"github.com/gotd/contrib/storage"
	"github.com/gotd/td/telegram"
	"github.com/gotd/td/telegram/auth/qrlogin"
	"github.com/gotd/td/telegram/query"
	"github.com/gotd/td/telegram/updates"
	"github.com/gotd/td/tg"
	"github.com/gotd/td/tgerr"
	"github.com/jedib0t/go-pretty/v6/text"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"github.com/skip2/go-qrcode"
	"go.etcd.io/bbolt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/term"
	"golang.org/x/time/rate"
	lj "gopkg.in/natefinch/lumberjack.v2"
)

// AccountClient хранит клиента и связанные с ним данные
type AccountClient struct {
	Client      *telegram.Client
	API         *tg.Client
	PeerDB      *pebble.PeerStorage
	UpdatesDB   *bbolt.DB
	PebbleDB    *pebbledb.DB
	Logger      *zap.Logger
	Cancel      context.CancelFunc
	IsRunning   bool
	AccountInfo *Account
}

// AccountManager управляет множеством Telegram аккаунтов
type AccountManager struct {
	mu       sync.RWMutex
	accounts map[int64]*AccountClient
	ctx      context.Context
	cancel   context.CancelFunc
	wg       sync.WaitGroup
	service  service.AccountService
}

// NewAccountManager создает новый менеджер аккаунтов
func NewAccountManager(accountService service.AccountService) *AccountManager {
	ctx, cancel := context.WithCancel(context.Background())
	return &AccountManager{
		accounts: make(map[int64]*AccountClient),
		ctx:      ctx,
		cancel:   cancel,
		service:  accountService,
	}
}

// LoadAccounts загружает все аккаунты из БД
func (am *AccountManager) LoadAccounts(bot *telego.Bot, updates <-chan telego.Update) error {
	accounts, err := am.service.GetAll()
	if err != nil {
		return fmt.Errorf("get accounts from db: %w", err)
	}

	log.Printf("Found %d active accounts in database\n", len(accounts))

	for _, acc := range accounts {
		log.Printf("Loading account %d (%s %s @%s)...\n", acc.ID, acc.FirstName, acc.LastName, acc.Username)
		if addErr := am.addAccountFromDB(acc, bot, updates); addErr != nil {
			log.Printf("Failed to load account %d: %v\n", acc.ID, addErr)
			// Деактивируем аккаунт при ошибке загрузки
			_ = am.service.SetActive(acc.ID, false)
			continue
		}
	}

	return nil
}

// AddNewAccount добавляет новый аккаунт и сохраняет в БД
func (am *AccountManager) AddNewAccount(accountID int64, apiID int, apiHash string, bot *telego.Bot, updates <-chan telego.Update) error {
	// Сначала сохраняем в БД
	account := &Account{
		ID:       accountID,
		ApiID:    apiID,
		ApiHash:  apiHash,
		IsActive: true,
	}

	if err := am.service.Create(account); err != nil {
		return fmt.Errorf("save account to db: %w", err)
	}

	// Затем добавляем и авторизуем
	return am.addAccountFromBot(account, true, bot, updates)
}

// addAccountFromDB добавляет аккаунт из БД (уже авторизованный)
func (am *AccountManager) addAccountFromDB(account *Account, bot *telego.Bot, updates <-chan telego.Update) error {
	return am.addAccountFromBot(account, false, bot, updates)
}

func waitForUserMessage(ctx context.Context, userID int64, updates <-chan telego.Update) (string, error) {
	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case u, ok := <-updates:
			if !ok {
				return "", fmt.Errorf("updates channel closed")
			}
			if u.Message == nil {
				continue
			}
			// убедимся, что сообщение от нужного пользователя
			if u.Message.From != nil && u.Message.From.ID == userID {
				// возвращаем текст (trimmed)
				return strings.TrimSpace(u.Message.Text), nil
			}
		}
	}
}

// addAccountFromBot внутренний метод добавления аккаунта
func (am *AccountManager) addAccountFromBot(account *Account, needAuth bool, bot *telego.Bot, u <-chan telego.Update) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Проверяем, не добавлен ли уже этот аккаунт
	if _, exists := am.accounts[account.ID]; exists {
		return fmt.Errorf("account %d already loaded", account.ID)
	}

	// Настройка сессии для каждого аккаунта
	sessionDir := filepath.Join("sessions", fmt.Sprintf("%d", account.ID))
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return fmt.Errorf("create session dir: %w", err)
	}

	logFilePath := filepath.Join(sessionDir, "log.jsonl")
	fmt.Printf("[Account %d] Storing session in %s, logs in %s\n", account.ID, sessionDir, logFilePath)

	// Настройка логирования
	logWriter := zapcore.AddSync(&lj.Logger{
		Filename:   logFilePath,
		MaxBackups: 3,
		MaxSize:    1,
		MaxAge:     7,
	})
	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		logWriter,
		zap.DebugLevel,
	)
	lg := zap.New(logCore).With(zap.Int64("account_id", account.ID))

	// Настройка хранилищ
	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join(sessionDir, "session.json"),
	}

	pebbleDB, err := pebbledb.Open(filepath.Join(sessionDir, "peers.pebble.db"), &pebbledb.Options{})
	if err != nil {
		_ = lg.Sync()
		return fmt.Errorf("create pebble storage: %w", err)
	}
	peerDB := pebble.NewPeerStorage(pebbleDB)

	boltDB, err := bbolt.Open(filepath.Join(sessionDir, "updates.bolt.db"), 0666, nil)
	if err != nil {
		_ = pebbleDB.Close()
		_ = lg.Sync()
		return fmt.Errorf("create bolt storage: %w", err)
	}

	// Настройка обработчиков обновлений
	dispatcher := tg.NewUpdateDispatcher()
	updateHandler := storage.UpdateHook(dispatcher, peerDB)
	updatesRecovery := updates.New(updates.Config{
		Handler: updateHandler,
		Logger:  lg.Named("updates.recovery"),
		Storage: boltstor.NewStateStorage(boltDB),
	})

	// Настройка middleware
	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		lg.Info("Got FLOOD_WAIT", zap.Duration("retry_after", wait.Duration))
	})

	// Создание клиента
	options := telegram.Options{
		Logger:         lg,
		SessionStorage: sessionStorage,
		UpdateHandler:  updatesRecovery,
		Middlewares: []telegram.Middleware{
			waiter,
			ratelimit.New(rate.Every(time.Millisecond*100), 5),
		},
	}
	client := telegram.NewClient(account.ApiID, account.ApiHash, options)

	// Сохраняем клиента
	accountClient := &AccountClient{
		Client:      client,
		API:         client.API(),
		PeerDB:      peerDB,
		UpdatesDB:   boltDB,
		PebbleDB:    pebbleDB,
		Logger:      lg,
		IsRunning:   false,
		AccountInfo: account,
	}
	am.accounts[account.ID] = accountClient

	// Проверяем авторизацию
	authCtx, authCancel := context.WithTimeout(am.ctx, 5*time.Minute)
	defer authCancel()

	err = waiter.Run(authCtx, func(ctx context.Context) error {
		return client.Run(ctx, func(ctx context.Context) error {
			authStatus, authErr := client.Auth().Status(ctx)
			if authErr != nil {
				return fmt.Errorf("get auth status: %w", authErr)
			}

			if !authStatus.Authorized {
				if !needAuth {
					// Если загружаем из БД, но сессия невалидна
					lg.Warn("Session expired, need re-authorization")
				}

				fmt.Printf("[Account %d] Authorization required. Starting QR auth...\n", account.ID)

				_, qrErr := client.QR().Auth(ctx, qrlogin.OnLoginToken(dispatcher), func(ctx context.Context, token qrlogin.Token) error {
					qr, qrcodeCreateErr := qrcode.New(token.URL(), qrcode.Medium)
					if qrcodeCreateErr != nil {
						return qrcodeCreateErr
					}

					// 1) сгенерировать PNG QR из token.URL()
					png, pngErr := qrcode.Encode(token.URL(), qrcode.Medium, 256)
					if pngErr != nil {
						return pngErr
					}

					// 2) отправить PNG картинку через бота
					// Пример с использованием tu.Photo и tu.FileBytes (адаптируйте, если у вас другой API)
					_, sendErr := bot.SendPhoto(ctx, tu.Photo(
						tu.ID(account.ID),
						tu.FileFromBytes(png, "qr.png"),
					).WithCaption("Сканируйте QR код этим аккаунтом в приложении Telegram"))
					if sendErr != nil {
						// если отправка фото не работает, как fallback — отправим ASCII-код в моноширинном формате
						// (тот самый код, который раньше печатался в терминале)
						code := qr.ToSmallString(false) // если у вас есть qr.ToSmallString; иначе пропустите
						//lines := strings.Count(code, "\n")
						_, _ = bot.SendMessage(ctx, tu.Message(
							tu.ID(account.ID),
							"Пожалуйста, отсканируйте этот QR (если картинка не отправилась):\n\n"+code,
						))
						// опционально: логирование
						fmt.Printf("\n[Account %d] QR (ascii):\n%s\n", account.ID, code)
						_ = sendErr // мы вернём исходную ошибку дальше если нужно
					}

					/*qr, qrcodeCreateErr := qrcode.New(token.URL(), qrcode.Medium)
					if qrcodeCreateErr != nil {
						return qrcodeCreateErr
					}

					code := qr.ToSmallString(false)
					lines := strings.Count(code, "\n")

					// Отправка сообщения с просьбой отсканировать qr code
					_, err = bot.SendMessage(ctx, tu.Message(
						tu.ID(account.ID),
						strings.Repeat(text.CursorUp.Sprint(), lines),
					))

					fmt.Printf("\n[Account %d] Scan this QR code:\n", account.ID)
					fmt.Print(code)
					fmt.Print(strings.Repeat(text.CursorUp.Sprint(), lines))
					return nil*/

					return nil
				})

				if qrErr != nil {
					if !tgerr.Is(qrErr, "SESSION_PASSWORD_NEEDED") {
						return fmt.Errorf("qr auth: %w", qrErr)
					}

					_, err = bot.SendMessage(ctx, tu.Message(
						tu.ID(account.ID),
						"Облачный пароль (cloud password) требуется. Пожалуйста, пришлите его в ответном сообщении этому боту.",
					))
					if err != nil {
						return fmt.Errorf("send ask password message: %w", err)
					}

					// ждём пароль от пользователя — используем контекст с таймаутом
					passCtx, passCancel := context.WithTimeout(ctx, 5*time.Minute)
					defer passCancel()

					passwordStr, waitErr := waitForUserMessage(passCtx, account.ID, u)
					if waitErr != nil {
						return fmt.Errorf("waiting for password failed: %w", waitErr)
					}

					if _, passwordErr := client.Auth().Password(ctx, passwordStr); passwordErr != nil {
						return fmt.Errorf("password auth: %w", passwordErr)
					}
				}
			}

			// Получаем информацию о пользователе
			self, selfErr := client.Self(ctx)
			if selfErr != nil {
				return fmt.Errorf("call self: %w", selfErr)
			}

			// Обновляем информацию в БД
			phone := ""
			if self.Phone != "" {
				phone = self.Phone
			}
			_ = am.service.UpdateUserInfo(account.ID, self.Username, self.FirstName, self.LastName, phone)

			name := self.FirstName
			if self.Username != "" {
				name = fmt.Sprintf("%s (@%s)", name, self.Username)
			}
			lg.Info("Authorized successfully", zap.String("user", name), zap.Int64("user_id", self.ID))
			fmt.Printf("[Account %d] Authorized as: %s (ID: %d)\n", account.ID, name, self.ID)

			return nil
		})
	})

	if err != nil {
		// Очистка при ошибке
		delete(am.accounts, account.ID)
		_ = boltDB.Close()
		_ = pebbleDB.Close()
		_ = lg.Sync()
		return fmt.Errorf("authorization failed: %w", err)
	}

	return nil
}

// addAccountInternal внутренний метод добавления аккаунта
func (am *AccountManager) addAccountInternal(account *Account, needAuth bool) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	// Проверяем, не добавлен ли уже этот аккаунт
	if _, exists := am.accounts[account.ID]; exists {
		return fmt.Errorf("account %d already loaded", account.ID)
	}

	// Настройка сессии для каждого аккаунта
	sessionDir := filepath.Join("sessions", fmt.Sprintf("%d", account.ID))
	if err := os.MkdirAll(sessionDir, 0700); err != nil {
		return fmt.Errorf("create session dir: %w", err)
	}

	logFilePath := filepath.Join(sessionDir, "log.jsonl")
	fmt.Printf("[Account %d] Storing session in %s, logs in %s\n", account.ID, sessionDir, logFilePath)

	// Настройка логирования
	logWriter := zapcore.AddSync(&lj.Logger{
		Filename:   logFilePath,
		MaxBackups: 3,
		MaxSize:    1,
		MaxAge:     7,
	})
	logCore := zapcore.NewCore(
		zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig()),
		logWriter,
		zap.DebugLevel,
	)
	lg := zap.New(logCore).With(zap.Int64("account_id", account.ID))

	// Настройка хранилищ
	sessionStorage := &telegram.FileSessionStorage{
		Path: filepath.Join(sessionDir, "session.json"),
	}

	pebbleDB, err := pebbledb.Open(filepath.Join(sessionDir, "peers.pebble.db"), &pebbledb.Options{})
	if err != nil {
		_ = lg.Sync()
		return fmt.Errorf("create pebble storage: %w", err)
	}
	peerDB := pebble.NewPeerStorage(pebbleDB)

	boltDB, err := bbolt.Open(filepath.Join(sessionDir, "updates.bolt.db"), 0666, nil)
	if err != nil {
		_ = pebbleDB.Close()
		_ = lg.Sync()
		return fmt.Errorf("create bolt storage: %w", err)
	}

	// Настройка обработчиков обновлений
	dispatcher := tg.NewUpdateDispatcher()
	updateHandler := storage.UpdateHook(dispatcher, peerDB)
	updatesRecovery := updates.New(updates.Config{
		Handler: updateHandler,
		Logger:  lg.Named("updates.recovery"),
		Storage: boltstor.NewStateStorage(boltDB),
	})

	// Настройка middleware
	waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
		lg.Info("Got FLOOD_WAIT", zap.Duration("retry_after", wait.Duration))
	})

	// Создание клиента
	options := telegram.Options{
		Logger:         lg,
		SessionStorage: sessionStorage,
		UpdateHandler:  updatesRecovery,
		Middlewares: []telegram.Middleware{
			waiter,
			ratelimit.New(rate.Every(time.Millisecond*100), 5),
		},
	}
	client := telegram.NewClient(account.ApiID, account.ApiHash, options)

	// Сохраняем клиента
	accountClient := &AccountClient{
		Client:      client,
		API:         client.API(),
		PeerDB:      peerDB,
		UpdatesDB:   boltDB,
		PebbleDB:    pebbleDB,
		Logger:      lg,
		IsRunning:   false,
		AccountInfo: account,
	}
	am.accounts[account.ID] = accountClient

	// Проверяем авторизацию
	authCtx, authCancel := context.WithTimeout(am.ctx, 5*time.Minute)
	defer authCancel()

	err = waiter.Run(authCtx, func(ctx context.Context) error {
		return client.Run(ctx, func(ctx context.Context) error {
			authStatus, authErr := client.Auth().Status(ctx)
			if authErr != nil {
				return fmt.Errorf("get auth status: %w", authErr)
			}

			if !authStatus.Authorized {
				if !needAuth {
					// Если загружаем из БД, но сессия невалидна
					lg.Warn("Session expired, need re-authorization")
				}

				fmt.Printf("[Account %d] Authorization required. Starting QR auth...\n", account.ID)

				_, qrErr := client.QR().Auth(ctx, qrlogin.OnLoginToken(dispatcher), func(ctx context.Context, token qrlogin.Token) error {
					qr, qrcodeCreateErr := qrcode.New(token.URL(), qrcode.Medium)
					if qrcodeCreateErr != nil {
						return qrcodeCreateErr
					}

					code := qr.ToSmallString(false)
					lines := strings.Count(code, "\n")

					fmt.Printf("\n[Account %d] Scan this QR code:\n", account.ID)
					fmt.Print(code)
					fmt.Print(strings.Repeat(text.CursorUp.Sprint(), lines))
					return nil
				})

				if qrErr != nil {
					if !tgerr.Is(qrErr, "SESSION_PASSWORD_NEEDED") {
						return fmt.Errorf("qr auth: %w", qrErr)
					}

					fmt.Printf("[Account %d] Enter cloud password: ", account.ID)
					password, readPasswordErr := term.ReadPassword(int(os.Stdin.Fd()))
					if readPasswordErr != nil {
						return fmt.Errorf("failed to read password: %w", readPasswordErr)
					}
					fmt.Println()

					passwordStr := strings.TrimSpace(string(password))
					if _, passwordErr := client.Auth().Password(ctx, passwordStr); passwordErr != nil {
						return fmt.Errorf("password auth: %w", passwordErr)
					}
				}
			}

			// Получаем информацию о пользователе
			self, selfErr := client.Self(ctx)
			if selfErr != nil {
				return fmt.Errorf("call self: %w", selfErr)
			}

			// Обновляем информацию в БД
			phone := ""
			if self.Phone != "" {
				phone = self.Phone
			}
			_ = am.service.UpdateUserInfo(account.ID, self.Username, self.FirstName, self.LastName, phone)

			name := self.FirstName
			if self.Username != "" {
				name = fmt.Sprintf("%s (@%s)", name, self.Username)
			}
			lg.Info("Authorized successfully", zap.String("user", name), zap.Int64("user_id", self.ID))
			fmt.Printf("[Account %d] Authorized as: %s (ID: %d)\n", account.ID, name, self.ID)

			return nil
		})
	})

	if err != nil {
		// Очистка при ошибке
		delete(am.accounts, account.ID)
		_ = boltDB.Close()
		_ = pebbleDB.Close()
		_ = lg.Sync()
		return fmt.Errorf("authorization failed: %w", err)
	}

	return nil
}

// GetClient возвращает клиента по ID аккаунта
func (am *AccountManager) GetClient(accountID int64) (*AccountClient, error) {
	am.mu.RLock()
	defer am.mu.RUnlock()

	client, exists := am.accounts[accountID]
	if !exists {
		return nil, fmt.Errorf("account %d not found", accountID)
	}

	return client, nil
}

// StartAll запускает все добавленные аккаунты
func (am *AccountManager) StartAll() error {
	am.mu.Lock()
	defer am.mu.Unlock()

	if len(am.accounts) == 0 {
		return fmt.Errorf("no accounts to start")
	}

	for accountID, accountClient := range am.accounts {
		if accountClient.IsRunning {
			continue
		}

		// Создаем контекст для этого аккаунта
		ctx, cancel := context.WithCancel(am.ctx)
		accountClient.Cancel = cancel

		am.wg.Add(1)
		go func(id int64, ac *AccountClient) {
			defer am.wg.Done()

			ac.Logger.Info("Starting account", zap.Int64("account_id", id))

			// Запускаем клиента с обработкой обновлений
			waiter := floodwait.NewWaiter().WithCallback(func(ctx context.Context, wait floodwait.FloodWait) {
				ac.Logger.Info("Got FLOOD_WAIT", zap.Duration("retry_after", wait.Duration))
			})

			err := waiter.Run(ctx, func(ctx context.Context) error {
				return ac.Client.Run(ctx, func(ctx context.Context) error {
					self, err := ac.Client.Self(ctx)
					if err != nil {
						return fmt.Errorf("get self: %w", err)
					}

					// Заполняем peer storage (опционально)
					if fillPeers := false; fillPeers { // Можно сделать настраиваемым
						ac.Logger.Info("Filling peer storage from dialogs")
						collector := storage.CollectPeers(ac.PeerDB)
						if err := collector.Dialogs(ctx, query.GetDialogs(ac.API).Iter()); err != nil {
							return fmt.Errorf("collect peers: %w", err)
						}
						ac.Logger.Info("Peer storage filled")
					}

					ac.Logger.Info("Listening for updates", zap.Int64("user_id", self.ID))

					// Получаем updates recovery из updateHandler
					dispatcher := tg.NewUpdateDispatcher()
					updateHandler := storage.UpdateHook(dispatcher, ac.PeerDB)
					updatesRecovery := updates.New(updates.Config{
						Handler: updateHandler,
						Logger:  ac.Logger.Named("updates.recovery"),
						Storage: boltstor.NewStateStorage(ac.UpdatesDB),
					})

					return updatesRecovery.Run(ctx, ac.API, self.ID, updates.AuthOptions{
						IsBot: self.Bot,
						OnStart: func(ctx context.Context) {
							ac.Logger.Info("Update recovery initialized and started")
						},
					})
				})
			})

			if err != nil {
				ac.Logger.Error("Account stopped with error", zap.Error(err))
			} else {
				ac.Logger.Info("Account stopped normally")
			}
		}(accountID, accountClient)

		accountClient.IsRunning = true
	}

	fmt.Printf("Started %d accounts\n", len(am.accounts))
	return nil
}

// Stop останавливает все аккаунты
func (am *AccountManager) Stop() {
	am.mu.Lock()
	defer am.mu.Unlock()

	log.Println("Stopping all accounts...")

	// Отменяем основной контекст
	am.cancel()

	// Ждем завершения всех горутин
	am.wg.Wait()

	// Закрываем все ресурсы
	for accountID, client := range am.accounts {
		if client.Cancel != nil {
			client.Cancel()
		}
		if client.UpdatesDB != nil {
			_ = client.UpdatesDB.Close()
		}
		if client.PebbleDB != nil {
			_ = client.PebbleDB.Close()
		}
		if client.Logger != nil {
			_ = client.Logger.Sync()
		}
		client.IsRunning = false
		log.Printf("Account %d stopped and cleaned up\n", accountID)
	}

	// Закрываем репозиторий
	if am.service != nil {
		_ = am.service.Close()
	}

	log.Println("All accounts stopped")
}

// GetAllAccounts возвращает список всех загруженных аккаунтов
func (am *AccountManager) GetAllAccounts() ([]*Account, error) {
	return am.service.GetAll()
}

// RemoveAccount удаляет аккаунт из менеджера и БД
func (am *AccountManager) RemoveAccount(accountID int64, deleteFromDB bool) error {
	am.mu.Lock()
	defer am.mu.Unlock()

	client, exists := am.accounts[accountID]
	if exists {
		// Останавливаем клиента
		if client.Cancel != nil {
			client.Cancel()
		}

		// Закрываем ресурсы
		if client.UpdatesDB != nil {
			_ = client.UpdatesDB.Close()
		}
		if client.PebbleDB != nil {
			_ = client.PebbleDB.Close()
		}
		if client.Logger != nil {
			_ = client.Logger.Sync()
		}

		delete(am.accounts, accountID)
	}

	// Удаляем или деактивируем в БД
	if deleteFromDB {
		if err := am.service.Delete(accountID); err != nil {
			return fmt.Errorf("\\delete from db: %w", err)
		}
	} else {
		if err := am.service.SetActive(accountID, false); err != nil {
			return fmt.Errorf("deactivate in db: %w", err)
		}
	}

	log.Printf("Account %d removed\n", accountID)
	return nil
}

// DeactivateAccount деактивирует аккаунт (остается в БД, но не загружается)
func (am *AccountManager) DeactivateAccount(accountID int64) error {
	return am.RemoveAccount(accountID, false)
}
