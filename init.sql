-- Включаем расширение uuid-ossp, если оно ещё не создано
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Создаём таблицу orders, если она не существует
CREATE TABLE IF NOT EXISTS orders (
    order_uid TEXT PRIMARY KEY,
    track_number TEXT,
    entry TEXT,
    locale TEXT,
    internal_signature TEXT,
    customer_id TEXT,
    delivery_service TEXT,
    shardkey TEXT,
    sm_id INT,
    date_created TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP,
    oof_shard TEXT
);
