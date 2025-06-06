-- noinspection SqlConstantExpressionForFile

-- Удаление индексов
DROP INDEX IF EXISTS idx_payments_user_id;
DROP INDEX IF EXISTS idx_payments_telegram_charge_id;

-- Удаление таблицы
DROP TABLE IF EXISTS payments;