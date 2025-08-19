-- Создание базы данных
CREATE DATABASE orders_db;

-- Создание пользователя
CREATE USER orders_user WITH PASSWORD 'orders_password';

-- Выдача прав
GRANT ALL PRIVILEGES ON DATABASE orders_db TO orders_user;

-- Подключение к базе и создание таблиц
\c orders_db

-- Таблица заказов
CREATE TABLE orders (
    order_uid VARCHAR(255) PRIMARY KEY,
    track_number VARCHAR(255),
    entry VARCHAR(50),
    locale VARCHAR(10),
    internal_signature VARCHAR(255),
    customer_id VARCHAR(255),
    delivery_service VARCHAR(255),
    shardkey VARCHAR(10),
    sm_id INT,
    date_created TIMESTAMP,
    oof_shard VARCHAR(10)
);

-- Таблица доставки
CREATE TABLE deliveries (
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    name VARCHAR(255),
    phone VARCHAR(50),
    zip VARCHAR(50),
    city VARCHAR(255),
    address VARCHAR(255),
    region VARCHAR(255),
    email VARCHAR(255),
    PRIMARY KEY (order_uid)
);

-- Таблица оплаты
CREATE TABLE payments (
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    transaction VARCHAR(255),
    request_id VARCHAR(255),
    currency VARCHAR(10),
    provider VARCHAR(50),
    amount INT,
    payment_dt BIGINT,
    bank VARCHAR(100),
    delivery_cost INT,
    goods_total INT,
    custom_fee INT,
    PRIMARY KEY (order_uid)
);

-- Таблица товаров
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    order_uid VARCHAR(255) REFERENCES orders(order_uid),
    chrt_id INT,
    track_number VARCHAR(255),
    price INT,
    rid VARCHAR(255),
    name VARCHAR(255),
    sale INT,
    size VARCHAR(50),
    total_price INT,
    nm_id INT,
    brand VARCHAR(255),
    status INT
);

-- Выдача прав пользователю на таблицы
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO orders_user;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO orders_user;