CREATE TABLE withdrawals (
                             id SERIAL PRIMARY KEY,
                             user_id INTEGER NOT NULL REFERENCES users(id),
                             order VARCHAR(255) NOT NULL,
                             sum FLOAT NOT NULL,
                             processed_at TIMESTAMP NOT NULL DEFAULT NOW()
);

