CREATE TABLE IF NOT EXISTS accounts
(
    id         INTEGER PRIMARY KEY,
    api_id     INTEGER NOT NULL,
    api_hash   TEXT    NOT NULL,
    phone      TEXT     DEFAULT '',
    username   TEXT     DEFAULT '',
    first_name TEXT     DEFAULT '',
    last_name  TEXT     DEFAULT '',
    is_active  BOOLEAN  DEFAULT true,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_accounts_active ON accounts (is_active);