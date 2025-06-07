package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken              string
	DatabaseURL           string
	AdminIDs              []int64
	LogLevel              string
	MonitorInterval       int
	NotificationChannelId int64
	LogChatId             int64
}

func Load() *Config {
	// Loading .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	adminIDsStr := os.Getenv("ADMIN_IDS")
	var adminIDs []int64

	if adminIDsStr != "" {
		for _, idStr := range strings.Split(adminIDsStr, ",") {
			if id, err := strconv.ParseInt(strings.TrimSpace(idStr), 10, 64); err == nil {
				adminIDs = append(adminIDs, id)
			}
		}
	}

	monitorInterval, _ := strconv.Atoi(getEnv("MONITOR_INTERVAL", "10"))
	notificationChannelId, _ := strconv.ParseInt(os.Getenv("NOTIFICATION_CHANNEL_ID"), 10, 64)
	logChatId, _ := strconv.ParseInt(os.Getenv("LOG_CHAT_ID"), 10, 64)

	return &Config{
		BotToken:              os.Getenv("BOT_TOKEN"),
		DatabaseURL:           getEnv("DATABASE_URL", "postgres://<username>:<password>@localhost:5432/gifts_dev?sslmode=disable"),
		AdminIDs:              adminIDs,
		LogLevel:              getEnv("LOG_LEVEL", "info"),
		MonitorInterval:       monitorInterval,
		NotificationChannelId: notificationChannelId,
		LogChatId:             logChatId,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
