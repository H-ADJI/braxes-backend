-- CREATE TABLE orders (
--     id INTEGER PRIMARY KEY,
--     platform_id TEXT UNIQUE NOT NULL,
--     order_number INTEGER UNIQUE NOT NULL,
--     is_processed INTEGER DEFAULT 0 NOT NULL,
--     total_price REAL NOT NULL,
--     customer_name TEXT,
--     processed_date INTEGER,
--     creation_date INTEGER NOT NULL
-- );
-- name: AddOrder :one
INSERT INTO
    orders (
        order_number,
        platform_id,
        creation_date,
        total_price,
        customer_name
    )
VALUES
    (?, ?, ?, ?, ?) RETURNING *;


-- name: GetOrder :one
SELECT
    *
FROM
    orders
WHERE
    id = ?
LIMIT
    1;


-- name: GetAllOrdersDescDate :many
SELECT
    *
FROM
    orders
ORDER BY
    creation_date DESC;


-- name: GetAllOrders :many
SELECT
    *
FROM
    orders;


-- name: GetProcessedOrders :many
-- sorting latest
SELECT
    *
FROM
    orders
WHERE
    is_processed = 1
ORDER BY
    processed_date DESC;


-- name: GetUnProcessedOrders :many
-- sorting oldest
SELECT
    *
FROM
    orders
WHERE
    is_processed = 0
ORDER BY
    processed_date ASC;


-- name: ProcessOrder :one
UPDATE orders
SET
    is_processed = 1,
    processed_date = ?
WHERE
    platform_id = ? RETURNING is_processed;


-- name: UnProcessOrder :one
UPDATE orders
SET
    is_processed = 0,
    processed_date = 0
WHERE
    platform_id = ? RETURNING is_processed;


-- name: DeleteAuthor :exec
DELETE FROM orders
WHERE
    id = ?;
