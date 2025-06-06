package app

import (
	"GiftBuyer/config"
	"GiftBuyer/internal/database"
	"GiftBuyer/internal/repository"
	"GiftBuyer/internal/service"
	"github.com/mymmrac/telego"
	"sync"
)

type State uint

const (
	_ State = iota
	StateWaitingTopUpAmount
	StateWaitingChannelUsername
)

type StateStorage struct {
	sync.RWMutex
	States map[int64]State
}

type App struct {
	DB           *database.DB
	Repos        *repository.Repositories
	Services     *service.Services
	Bot          *telego.Bot
	Config       *config.Config
	StateStorage *StateStorage
}

// GetState Получить текущее состояние пользователя
func (storage *StateStorage) GetState(userID int64) State {
	storage.RLock()
	defer storage.RUnlock()
	return storage.States[userID]
}

// SetState Установить новое состояние для пользователя
func (storage *StateStorage) SetState(userID int64, state State) {
	storage.Lock()
	defer storage.Unlock()
	storage.States[userID] = state
}

// ClearState Удалить состояние пользователя
func (storage *StateStorage) ClearState(userID int64) {
	storage.Lock()
	defer storage.Unlock()
	delete(storage.States, userID)
}
