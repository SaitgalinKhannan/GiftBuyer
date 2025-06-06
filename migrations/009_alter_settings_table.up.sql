-- noinspection SqlConstantExpressionForFile

-- Добавляем колонку channels
ALTER TABLE user_settings
    ADD COLUMN IF NOT EXISTS channels TEXT DEFAULT NULL;

-- Комментарий к колонке
COMMENT ON COLUMN user_settings.channels IS 'Список каналов для автопокупки, разделённых запятой';