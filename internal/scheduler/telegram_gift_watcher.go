package scheduler

import (
	"GiftBuyer/app"
	"GiftBuyer/logging"
	"context"
	"fmt"
	"log"
	"runtime/debug"
	"time"
)

func StartGiftWatcher(ctx context.Context, a *app.App) {
	// Дефолтный интервал: 10 секунд
	interval := 10 * time.Second

	// Проверяем, есть ли конфиг и задан ли интервал
	if a == nil || a.Config == nil || a.Config.MonitorInterval <= 0 {
		log.Println("Используется дефолтный интервал (10s), так как MonitorInterval не задан или <= 0")
		interval = 10 * time.Second
	} else {
		interval = time.Duration(a.Config.MonitorInterval) * time.Second
	}

	log.Printf("GiftWatcher запущен с интервалом %v", interval)

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			logging.SendLogMessageToTelegram(ctx, a.Bot, a.Config.LogChatId, "Завершение фоновой задачи")
			log.Println("Завершение фоновой задачи")
			return
		case t := <-ticker.C:
			fmt.Printf("Выполняю задачу в %v\n", t)
			func() {
				defer func() {
					if r := recover(); r != nil {
						logging.SendLogMessageToTelegram(ctx, a.Bot, a.Config.LogChatId, "Восстановлено после паники telegram_gift_watcher.go")
						log.Printf("Восстановлено после паники: %v\n%s", r, debug.Stack())
					}
				}()
				if err := a.Services.Gift.CheckAndProcessNewGifts(ctx, a.Bot); err != nil {
					logging.SendLogErrorToTelegram(ctx, a.Bot, a.Config.LogChatId, err)
					log.Printf("Ошибка проверки подарков: %v", err)
				}
			}()
		}
	}
}
