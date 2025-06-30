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
CREATE TABLE deliveries (
    order_uid TEXT PRIMARY KEY,
    name TEXT,
    phone TEXT,
    zip TEXT,
    city TEXT,
    address TEXT,
    region TEXT,
    email TEXT
);

CREATE TABLE payments (
    order_uid TEXT PRIMARY KEY,
    transaction TEXT,
    request_id TEXT,
    currency TEXT,
    provider TEXT,
    amount INT,
    payment_dt BIGINT,
    bank TEXT,
    delivery_cost INT,
    goods_total INT,
    custom_fee INT
);

CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid TEXT,
    chrt_id INT,
    track_number TEXT,
    price INT,
    rid TEXT,
    name TEXT,
    sale INT,
    size TEXT,
    total_price INT,
    nm_id INT,
    brand TEXT,
    status INT
);
