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
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE INDEX idx_uploaded_at ON orders(uploaded_at);

CREATE TABLE IF NOT EXISTS balances
(
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    current_balance REAL NOT NULL,
    total_withdrawn REAL NOT NULL,
    CONSTRAINT fk_user_balance
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS withdrawals
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL,
    "order" VARCHAR(255) NOT NULL,
    sum REAL NOT NULL,
    processed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_user_withdrawals
        FOREIGN KEY(user_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_order_withdrawals
        FOREIGN KEY("order")
            REFERENCES orders(number)
            ON DELETE CASCADE
);

COMMIT;
