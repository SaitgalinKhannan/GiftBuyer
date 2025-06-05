-- Добавляем новые поля
ALTER TABLE users
    ADD COLUMN IF NOT EXISTS first_name VARCHAR(100),
    ADD COLUMN IF NOT EXISTS last_name  VARCHAR(100);

-- Изменяем тип balance на INTEGER с округлением
ALTER TABLE users
    ALTER COLUMN balance TYPE INTEGER USING ROUND(balance);

-- Создаем индексы (если планируете фильтровать по имени/фамилии)
-- CREATE INDEX IF NOT EXISTS idx_users_first_name ON users (first_name);
-- CREATE INDEX IF NOT EXISTS idx_users_last_name ON users (last_name);