BEGIN TRANSACTION;
-- Сначала удаляем таблицу withdrawals, так как она ссылается на orders и users
DROP TABLE IF EXISTS withdrawals;

-- Затем удаляем таблицу balances, так как она ссылается на users
DROP TABLE IF EXISTS balances;

-- Затем удаляем таблицу orders, так как она ссылается на users
DROP TABLE IF EXISTS orders;

-- И, наконец, удаляем таблицу users
DROP TABLE IF EXISTS users;

COMMIT;