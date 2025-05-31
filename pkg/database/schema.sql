CREATE TABLE IF NOT EXISTS rummage_items (
    entry TEXT NOT NULL UNIQUE,
    score FLOAT NOT NULL,
    lastaccessed INTEGER NOT NULL
);
