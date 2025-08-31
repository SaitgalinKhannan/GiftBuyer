package repository

import (
	"GiftBuyer/internal/database"
	. "GiftBuyer/internal/model"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// AccountRepository управляет хранением аккаунтов в БД
type accountRepository struct {
	db *sqlx.DB
}

// NewAccountRepository создает новый репозиторий
func NewAccountRepository(db *database.DB) AccountRepository {
	return &accountRepository{db: db.DB}
}

// Save сохраняет или обновляет аккаунт
func (r *accountRepository) Save(account *Account) error {
	query := `
	INSERT INTO accounts (id, api_id, api_hash, phone, username, first_name, last_name, is_active, updated_at)
	VALUES (:id, :api_id, :api_hash, :phone, :username, :first_name, :last_name, :is_active, CURRENT_TIMESTAMP)
	ON CONFLICT(id) DO UPDATE SET
		api_id = excluded.api_id,
		api_hash = excluded.api_hash,
		phone = excluded.phone,
		username = excluded.username,
		first_name = excluded.first_name,
		last_name = excluded.last_name,
		is_active = excluded.is_active,
		updated_at = CURRENT_TIMESTAMP
	`

	_, err := r.db.NamedExec(query, account)
	return err
}

// GetAll возвращает все активные аккаунты
func (r *accountRepository) GetAll() ([]Account, error) {
	var accounts []Account
	err := r.db.Select(&accounts, "SELECT * FROM accounts WHERE is_active = true ORDER BY id")
	return accounts, err
}

// GetByID возвращает аккаунт по ID
func (r *accountRepository) GetByID(id int64) (*Account, error) {
	var account Account
	err := r.db.Get(&account, "SELECT * FROM accounts WHERE id = ?", id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("account %d not found", id)
	}
	return &account, err
}

// Delete удаляет аккаунт
func (r *accountRepository) Delete(id int64) error {
	_, err := r.db.Exec("DELETE FROM accounts WHERE id = ?", id)
	return err
}

// SetActive устанавливает статус активности аккаунта
func (r *accountRepository) SetActive(id int64, active bool) error {
	_, err := r.db.Exec("UPDATE accounts SET is_active = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?", active, id)
	return err
}

// UpdateUserInfo обновляет информацию о пользователе
func (r *accountRepository) UpdateUserInfo(id int64, username, firstName, lastName, phone string) error {
	_, err := r.db.Exec(`
		UPDATE accounts 
		SET username = ?, first_name = ?, last_name = ?, phone = ?, updated_at = CURRENT_TIMESTAMP 
		WHERE id = ?`,
		username, firstName, lastName, phone, id)
	return err
}

// Close закрывает соединение с БД
func (r *accountRepository) Close() error {
	return r.db.Close()
}
