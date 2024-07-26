BEGIN;

CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY,
    login      VARCHAR(255) NOT NULL UNIQUE,
    password   VARCHAR(255) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    accrual REAL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_user_id ON orders(user_id);

CREATE TABLE balances
(
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    current_balance REAL NOT NULL,
    total_withdrawn REAL NOT NULL
);

CREATE INDEX idx_balances_user_id ON balances(user_id);

CREATE TABLE IF NOT EXISTS withdrawals
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    "order" VARCHAR(255) NOT NULL,
    sum REAL NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_withdrawals_user_id ON withdrawals(user_id);

COMMIT;
