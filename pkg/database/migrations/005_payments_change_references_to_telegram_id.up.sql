-- 1. Удаление старого внешнего ключа
ALTER TABLE payments
    DROP CONSTRAINT IF EXISTS payments_user_id_fkey;

-- 2. Обновление данных (если user_id хранит users.id)
UPDATE payments p
SET user_id = u.telegram_id
FROM users u
WHERE p.user_id = u.id;

-- 3. Добавление нового внешнего ключа
ALTER TABLE payments
    ADD CONSTRAINT fk_payments_user_telegram_id
        FOREIGN KEY (user_id)
            REFERENCES users (telegram_id)
            ON DELETE CASCADE;