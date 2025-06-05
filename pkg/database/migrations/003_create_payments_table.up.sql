-- Создание таблицы payments
CREATE TABLE IF NOT EXISTS payments
(
    id                         SERIAL PRIMARY KEY,
    user_id                    BIGINT      NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    currency                   VARCHAR(10) NOT NULL,
    amount                     INT         NOT NULL,
    payload                    TEXT        NOT NULL,
    telegram_payment_charge_id TEXT UNIQUE NOT NULL,
    created_at                 TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс по user_id для быстрого поиска платежей пользователя
CREATE INDEX IF NOT EXISTS idx_payments_user_id ON payments (user_id);

-- Индекс по telegram_payment_charge_id (уникальный)
CREATE INDEX IF NOT EXISTS idx_payments_telegram_charge_id ON payments (telegram_payment_charge_id);