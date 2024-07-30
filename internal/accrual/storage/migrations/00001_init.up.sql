BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS a_orders (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    uuid        VARCHAR(255) NOT NULL UNIQUE,
    status      VARCHAR(255) NOT NULL,
    accrual     REAL         NOT NULL,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS a_order_items (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    order_id    INT          NOT NULL,
    price       REAL         NOT NULL,
    description VARCHAR(255) NOT NULL
);

ALTER TABLE a_order_items
    ADD CONSTRAINT a_orders_order_id
        FOREIGN KEY (order_id)
            REFERENCES a_orders (id);

CREATE TABLE IF NOT EXISTS a_accrual_rules (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    match       VARCHAR(255) NOT NULL UNIQUE,
    reward      INT          NOT NULL,
    reward_type VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS users (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    login       VARCHAR(255) NOT NULL UNIQUE,
    password    VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders (
    id          INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id     INTEGER      NOT NULL,
    number      VARCHAR(255) NOT NULL UNIQUE,
    status      VARCHAR(50)  NOT NULL,
    accrual     REAL,
    uploaded_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

ALTER TABLE orders
    ADD CONSTRAINT orders_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

CREATE TABLE IF NOT EXISTS balances (
    user_id         INTEGER PRIMARY KEY,
    current_balance REAL    NOT NULL,
    total_withdrawn REAL    NOT NULL
);

CREATE INDEX idx_balances_user_id ON balances(user_id);

ALTER TABLE balances
    ADD CONSTRAINT balances_user_id
        FOREIGN KEY (user_id)
            REFERENCES users (id);

CREATE TABLE IF NOT EXISTS withdrawals (
    id            INT PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    user_id       INTEGER      NOT NULL,
    order_number  VARCHAR(255) NOT NULL,
    sum           REAL         NOT NULL,
    processed_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);

COMMIT;
