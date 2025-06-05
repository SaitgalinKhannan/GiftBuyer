-- 1. Создаём резервную копию данных
CREATE TABLE IF NOT EXISTS users_backup AS
SELECT id, first_name, last_name, balance FROM users;

-- 2. Переименовываем поля вместо удаления
ALTER TABLE users
    RENAME COLUMN first_name TO first_name_old,
RENAME COLUMN last_name TO last_name_old;