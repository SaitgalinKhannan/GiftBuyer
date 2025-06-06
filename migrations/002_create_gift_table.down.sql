-- noinspection SqlConstantExpressionForFile

-- Удаление таблицы gifts
DROP TABLE IF EXISTS gifts;

-- Удаление индекса (если существует)
DROP INDEX IF EXISTS idx_gifts_id;
DROP INDEX IF EXISTS idx_gifts_created_at;