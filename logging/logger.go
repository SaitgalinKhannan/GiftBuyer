package logging

import (
	"context"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
	"io"
	"log"
	"os"
)

func InitLogger() (*os.File, error) {
	logFile, err := os.OpenFile("bot.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	// Пишем и в файл, и в консоль
	multiWriter := io.MultiWriter(os.Stderr, logFile)

	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags)

	return logFile, nil
}

func SendLogErrorToTelegram(ctx context.Context, bot *telego.Bot, chatID int64, err error) {
	sendLog(ctx, bot, chatID, "❌ "+err.Error())
}

func SendLogMessageToTelegram(ctx context.Context, bot *telego.Bot, chatID int64, message string) {
	sendLog(ctx, bot, chatID, message)
}

// Вспомогательная функция
func sendLog(ctx context.Context, bot *telego.Bot, chatID int64, text string) {
	_, _ = bot.SendMessage(ctx, tu.Message(
		telego.ChatID{ID: chatID},
		text,
	))
}
