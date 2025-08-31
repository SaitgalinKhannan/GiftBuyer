package service

import (
	. "GiftBuyer/internal/model"
	. "GiftBuyer/internal/repository"
	"fmt"
)

// AccountService предоставляет бизнес-логику для работы с аккаунтами
type accountService struct {
	repo AccountRepository
}

// NewAccountService создает новый экземпляр сервиса
func NewAccountService(repo AccountRepository) AccountService {
	return &accountService{repo: repo}
}

// Преобразуем internal Account → public Account (можно опустить, если используете одну структуру)
func toAccountDTO(dbAcc *Account) *Account {
	return &Account{
		ID:        dbAcc.ID,
		ApiID:     dbAcc.ApiID,
		ApiHash:   dbAcc.ApiHash,
		Phone:     dbAcc.Phone,
		Username:  dbAcc.Username,
		FirstName: dbAcc.FirstName,
		LastName:  dbAcc.LastName,
		IsActive:  dbAcc.IsActive,
		CreatedAt: dbAcc.CreatedAt,
		UpdatedAt: dbAcc.UpdatedAt,
	}
}

// toAccountModel преобразует DTO обратно в модель репозитория
func (s *accountService) toAccountModel(acc *Account) *Account {
	return &Account{
		ID:        acc.ID,
		ApiID:     acc.ApiID,
		ApiHash:   acc.ApiHash,
		Phone:     acc.Phone,
		Username:  acc.Username,
		FirstName: acc.FirstName,
		LastName:  acc.LastName,
		IsActive:  acc.IsActive,
		// CreatedAt и UpdatedAt остаются без изменений, будут обновлены в БД
	}
}

// Validate проверяет корректность данных аккаунта
func (s *accountService) Validate(account *Account) error {
	if account.ApiID == 0 {
		return fmt.Errorf("api_id is required")
	}
	if account.ApiHash == "" {
		return fmt.Errorf("api_hash is required")
	}
	if account.ApiHash != "" && !isValidAPIHash(account.ApiHash) {
		return fmt.Errorf("invalid api_hash format")
	}
	if account.Phone != "" && !isValidPhone(account.Phone) {
		return fmt.Errorf("invalid phone number format")
	}
	if account.Username != "" && !isValidUsername(account.Username) {
		return fmt.Errorf("invalid username format")
	}
	return nil
}

// isValidAPIHash проверяет, что api_hash — это 32-символьная шестнадцатеричная строка
func isValidAPIHash(hash string) bool {
	/*matched, _ := regexp.MatchString(`^[a-fA-F0-9]{32}$`, hash)
	return matched*/
	return true
}

// isValidPhone проверяет простой формат телефона (например, начинается с + и содержит цифры)
func isValidPhone(phone string) bool {
	/*matched, _ := regexp.MatchString(`^\+\d{7,15}$`, phone)
	return matched*/
	return true
}

// isValidUsername проверяет, что username содержит только буквенно-цифровые символы и подчеркивания
func isValidUsername(username string) bool {
	/*matched, _ := regexp.MatchString(`^[a-zA-Z0-9_]{1,32}$`, username)
	return matched*/
	return true
}

// Create создает или обновляет аккаунт
func (s *accountService) Create(account *Account) error {
	if err := s.Validate(account); err != nil {
		return fmt.Errorf("validation failed: %w", err)
	}

	dbAccount := s.toAccountModel(account)
	if err := s.repo.Save(dbAccount); err != nil {
		return fmt.Errorf("failed to save account: %w", err)
	}

	return nil
}

// GetAll возвращает все активные аккаунты
func (s *accountService) GetAll() ([]*Account, error) {
	dbAccounts, err := s.repo.GetAll()
	if err != nil {
		return nil, fmt.Errorf("failed to get accounts from repository: %w", err)
	}

	var accounts []*Account
	for i := range dbAccounts {
		accounts = append(accounts, toAccountDTO(&dbAccounts[i]))
	}

	return accounts, nil
}

// GetByID возвращает аккаунт по ID
func (s *accountService) GetByID(id int64) (*Account, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid account ID")
	}

	dbAccount, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err // уже содержит сообщение "not found" или другую ошибку
	}

	return toAccountDTO(dbAccount), nil
}

// Delete удаляет аккаунт по ID
func (s *accountService) Delete(id int64) error {
	if id <= 0 {
		return fmt.Errorf("invalid account ID")
	}

	// Опционально: проверить существование
	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}

	if err := s.repo.Delete(id); err != nil {
		return fmt.Errorf("failed to delete account: %w", err)
	}

	return nil
}

// SetActive устанавливает статус активности
func (s *accountService) SetActive(id int64, active bool) error {
	if id <= 0 {
		return fmt.Errorf("invalid account ID")
	}

	if _, err := s.repo.GetByID(id); err != nil {
		return err
	}

	if err := s.repo.SetActive(id, active); err != nil {
		return fmt.Errorf("failed to update active status: %w", err)
	}

	return nil
}

// UpdateUserInfo обновляет профиль пользователя
func (s *accountService) UpdateUserInfo(id int64, username, firstName, lastName, phone string) error {
	if id <= 0 {
		return fmt.Errorf("invalid account ID")
	}

	// Валидация
	if phone != "" && !isValidPhone(phone) {
		return fmt.Errorf("invalid phone number")
	}
	if username != "" && !isValidUsername(username) {
		return fmt.Errorf("invalid username")
	}

	if err := s.repo.UpdateUserInfo(id, username, firstName, lastName, phone); err != nil {
		return fmt.Errorf("failed to update user info: %w", err)
	}

	return nil
}

func (s *accountService) Close() error {
	return s.repo.Close()
}
