 CREATE TABLE orders (
                        id SERIAL PRIMARY KEY,
                        user_id INTEGER NOT NULL REFERENCES users(id),
                        number VARCHAR(255) NOT NULL UNIQUE,
                        status VARCHAR(50) NOT NULL,
                        accrual FLOAT,
                        uploaded_at TIMESTAMP NOT NULL DEFAULT NOW()
);

