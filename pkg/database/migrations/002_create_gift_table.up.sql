-- Создание таблицы gifts
CREATE TABLE IF NOT EXISTS gifts
(
    id                 VARCHAR PRIMARY KEY,
    star_count         INT NOT NULL,
    upgrade_star_count INT,
    total_count        INT,
    remaining_count    INT,
    created_at         TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Индекс по id для поиска
CREATE INDEX IF NOT EXISTS idx_gifts_id ON gifts (id);
-- Индекс по created_at для сортировки
CREATE INDEX IF NOT EXISTS idx_gifts_created_at ON gifts (created_at);