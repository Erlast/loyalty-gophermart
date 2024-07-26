CREATE TABLE balances
(
    user_id INTEGER PRIMARY KEY REFERENCES users(id),
    current_balance REAL NOT NULL,
    total_withdrawn REAL NOT NULL
);
