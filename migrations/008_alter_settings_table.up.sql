-- noinspection SqlConstantExpressionForFile

-- Делаем поле NOT NULL, если все записи обновлены
ALTER TABLE user_settings
    ALTER COLUMN auto_buy_cycles SET NOT NULL;