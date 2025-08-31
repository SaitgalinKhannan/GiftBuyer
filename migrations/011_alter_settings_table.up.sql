-- noinspection SqlConstantExpressionForFile

-- Добавляем колонку channels
ALTER TABLE user_settings
    ADD COLUMN IF NOT EXISTS only_premium_gift BOOLEAN DEFAULT FALSE;

-- Комментарий к колонке
COMMENT ON COLUMN user_settings.only_premium_gift IS 'Покупка только премиум подарков';