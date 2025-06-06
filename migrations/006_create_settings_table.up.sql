-- noinspection SqlConstantExpressionForFile

-- Создание таблицы user_settings
CREATE TABLE IF NOT EXISTS user_settings
(
    id               SERIAL PRIMARY KEY,
    user_id          BIGINT  NOT NULL REFERENCES users (id) ON DELETE CASCADE,
    auto_buy_enabled BOOLEAN NOT NULL DEFAULT FALSE,
    price_limit_from INT              DEFAULT NULL,
    price_limit_to   INT              DEFAULT NULL,
    supply_limit     INT              DEFAULT NULL,
    auto_buy_cycles  INT              DEFAULT NULL,
    created_at       TIMESTAMP        DEFAULT CURRENT_TIMESTAMP,
    updated_at       TIMESTAMP        DEFAULT CURRENT_TIMESTAMP
);

-- Индекс по user_id для быстрого поиска настроек
CREATE INDEX IF NOT EXISTS idx_user_settings_user_id ON user_settings (user_id);

-- Комментарии для документации
COMMENT ON TABLE user_settings IS 'Настройки пользователей для автопокупки и лимитов';
COMMENT ON COLUMN user_settings.user_id IS 'Ссылка на пользователя';
COMMENT ON COLUMN user_settings.auto_buy_enabled IS 'Включена ли автопокупка';
COMMENT ON COLUMN user_settings.price_limit_from IS 'Минимальная цена подарка для автопокупки';
COMMENT ON COLUMN user_settings.price_limit_to IS 'Максимальная цена подарка для автопокупки';
COMMENT ON COLUMN user_settings.supply_limit IS 'Лимит количества подарков для автопокупки';
COMMENT ON COLUMN user_settings.auto_buy_cycles IS 'Количество циклов автопокупки';