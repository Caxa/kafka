-- Создание пользователя и БД (выполняется от имени postgres)
CREATE USER order_user WITH PASSWORD 'password';
CREATE DATABASE ordersdb OWNER order_user;
GRANT ALL PRIVILEGES ON DATABASE ordersdb TO order_user;

-- Подключение к БД ordersdb и создание таблиц
\c ordersdb
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE orders (
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
CREATE TABLE delivery (
        order_uid TEXT PRIMARY KEY REFERENCES orders(order_uid) ON DELETE CASCADE,
        name TEXT NOT NULL,
        phone TEXT NOT NULL,
        zip TEXT NOT NULL,
        city TEXT NOT NULL,
        address TEXT NOT NULL,
        region TEXT,
        email TEXT
    );

    CREATE TABLE payment (
        transaction TEXT PRIMARY KEY,
        order_uid TEXT UNIQUE REFERENCES orders(order_uid) ON DELETE CASCADE,
        request_id TEXT,
        currency TEXT NOT NULL,
        provider TEXT NOT NULL,
        amount INT NOT NULL CHECK (amount > 0),
        payment_dt BIGINT NOT NULL,
        bank TEXT NOT NULL,
        delivery_cost INT DEFAULT 0,
        goods_total INT NOT NULL,
        custom_fee INT DEFAULT 0
    );

    CREATE TABLE items (
        id SERIAL PRIMARY KEY,
        order_uid TEXT REFERENCES orders(order_uid) ON DELETE CASCADE,
        chrt_id INT NOT NULL,
        track_number TEXT NOT NULL,
        price INT NOT NULL CHECK (price > 0),
        rid TEXT NOT NULL,
        name TEXT NOT NULL,
        sale INT DEFAULT 0,
        size TEXT NOT NULL,
        total_price INT NOT NULL,
        nm_id INT NOT NULL,
        brand TEXT NOT NULL,
        status INT NOT NULL
    );

    GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO order_user;
    GRANT USAGE, SELECT ON ALL SEQUENCES IN SCHEMA public TO order_user;