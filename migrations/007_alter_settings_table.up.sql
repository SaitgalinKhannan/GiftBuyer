-- noinspection SqlConstantExpressionForFile

-- Изменяем значение по умолчанию на -1
ALTER TABLE user_settings
    ALTER COLUMN auto_buy_cycles SET DEFAULT -1;

-- Опционально: обновляем существующие записи, где auto_buy_cycles IS NULL
UPDATE user_settings
SET auto_buy_cycles = -1
WHERE auto_buy_cycles IS NULL;