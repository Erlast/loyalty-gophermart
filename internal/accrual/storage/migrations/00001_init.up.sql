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

COMMIT;