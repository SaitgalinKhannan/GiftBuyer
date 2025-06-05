package app

import (
	"GiftBuyer/config"
	"GiftBuyer/internal/repository"
	"GiftBuyer/pkg/database"
	"github.com/mymmrac/telego"
)

type App struct {
	DB     *database.DB
	Repos  *repository.Repositories
	Bot    *telego.Bot
	Config *config.Config
}
