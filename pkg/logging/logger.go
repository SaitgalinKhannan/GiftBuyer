package logging

import (
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
