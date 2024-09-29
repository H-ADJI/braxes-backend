CREATE TABLE orders (
    id INTEGER PRIMARY KEY,
    platform_id INTEGER UNIQUE NOT NULL,
    is_processed INTEGER DEFAULT 0 NOT NULL,
    processed_date INTEGER,
    creation_date INTEGER NOT NULL
);
