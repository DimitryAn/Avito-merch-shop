CREATE SCHEMA IF NOT EXISTS shop;

-- пользователи
CREATE TABLE IF NOT EXISTS shop.users (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL, -- должно быть уникальным
    password TEXT NOT NULL,
    balance INT NOT NULL DEFAULT 1000 CHECK (balance >= 0)
);


