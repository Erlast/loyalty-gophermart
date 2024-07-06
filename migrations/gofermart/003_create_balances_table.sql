CREATE TABLE balances (
                          user_id INTEGER PRIMARY KEY REFERENCES users(id),
                          current_balance FLOAT NOT NULL,
                          total_withdrawn FLOAT NOT NULL
);
