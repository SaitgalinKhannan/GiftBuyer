package scheduler

import (
	"GiftBuyer/app"
	"context"
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
			log.Println("Завершение фоновой задачи")
			return
		case t := <-ticker.C:
			log.Printf("Выполняю задачу в %v", t)
			func() {
				defer func() {
					if r := recover(); r != nil {
						log.Printf("Восстановлено после паники: %v\n%s", r, debug.Stack())
					}
				}()
				if err := a.Services.Gift.CheckAndProcessNewGifts(ctx, a.Bot); err != nil {
					log.Printf("Ошибка проверки подарков: %v", err)
				}
			}()
		}
	}
}
