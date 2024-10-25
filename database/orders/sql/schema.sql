CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    platform_id TEXT UNIQUE NOT NULL,
    order_number INTEGER UNIQUE NOT NULL,
    is_processed INTEGER DEFAULT 0 NOT NULL,
    total_price REAL NOT NULL,
    customer_name TEXT NOT NULL,
    processed_date INTEGER,
    creation_date INTEGER NOT NULL
);
