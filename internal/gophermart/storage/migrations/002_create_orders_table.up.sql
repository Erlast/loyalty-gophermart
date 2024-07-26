BEGIN;

CREATE TABLE IF NOT EXISTS orders
(
    id SERIAL PRIMARY KEY,
    user_id INTEGER NOT NULL REFERENCES users(id),
    number VARCHAR(255) NOT NULL UNIQUE,
    status VARCHAR(50) NOT NULL,
    accrual REAL,
    uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_uploaded_at ON orders(uploaded_at);

COMMIT;
