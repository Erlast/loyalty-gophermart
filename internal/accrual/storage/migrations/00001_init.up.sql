BEGIN TRANSACTION;
CREATE TABLE IF NOT EXISTS a_orders
(
    id          SERIAL PRIMARY KEY,
    uuid        VARCHAR(255) NOT NULL UNIQUE,
    status      VARCHAR(255) NOT NULL,
    accrual     REAL          NOT NULL,
    uploaded_at TIMESTAMP default CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS a_order_items
(
    id          SERIAL PRIMARY KEY,
    order_id    INT          NOT NULL,
    price       REAL         NOT NULL,
    description VARCHAR(255) NOT NULL
);

ALTER TABLE a_order_items
    ADD CONSTRAINT a_orders_order_id
        FOREIGN KEY (order_id)
            REFERENCES a_orders (id);

CREATE TABLE IF NOT EXISTS a_accrual_rules
(
    id          SERIAL PRIMARY KEY,
    match       VARCHAR(255) NOT NULL UNIQUE,
    reward      int          NOT NULL,
    reward_type VARCHAR(255) NOT NULL
);
-- Создание таблицы пользователей
CREATE TABLE IF NOT EXISTS users (
                                     id SERIAL PRIMARY KEY,
                                     login VARCHAR(255) NOT NULL UNIQUE,
                                     password VARCHAR(255) NOT NULL,
                                     created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                                     updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Создание таблицы заказов
CREATE TABLE IF NOT EXISTS orders (
                                      id SERIAL PRIMARY KEY,
                                      user_id INTEGER NOT NULL,  -- Убрано REFERENCES users(id)
                                      number VARCHAR(255) NOT NULL UNIQUE,
                                      status VARCHAR(50) NOT NULL,
                                      accrual REAL,
                                      uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для улучшения производительности запросов по user_id в таблице orders
CREATE INDEX idx_orders_user_id ON orders(user_id);

ALTER TABLE orders
    ADD CONSTRAINT orders_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

-- Создание таблицы балансов
CREATE TABLE IF NOT EXISTS balances (
                                        user_id INTEGER PRIMARY KEY,  -- Убрано REFERENCES users(id)
                                        current_balance REAL NOT NULL,
                                        total_withdrawn REAL NOT NULL
);

-- Индекс для улучшения производительности запросов по user_id в таблице balances
CREATE INDEX idx_balances_user_id ON balances(user_id);

ALTER TABLE balances
    ADD CONSTRAINT balances_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

-- Создание таблицы выводов средств
CREATE TABLE IF NOT EXISTS withdrawals (
                                           id SERIAL PRIMARY KEY,
                                           user_id INTEGER NOT NULL,  -- Убрано REFERENCES users(id)
                                           order_number VARCHAR(255) NOT NULL,  -- Изменено имя столбца для избежания использования зарезервированного слова
                                           sum REAL NOT NULL,
                                           processed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Индекс для улучшения производительности запросов по user_id в таблице withdrawals
CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);

COMMIT;