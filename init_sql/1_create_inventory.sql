--CREATE TABLE
--inventory
CREATE TYPE uints AS ENUM ('g', 'ml', 'pcs');

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    quantity FLOAT NOT NULL CHECK (quantity >= 0),
    reorder_level FLOAT NOT NULL CHECK (reorder_level > 0),
    unit uints NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);

CREATE TYPE reason_of_inventory_transaction AS ENUM ('restock', 'usage', 'cancelled', 'annul');

CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INT NOT NULL REFERENCES inventory (id) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason reason_of_inventory_transaction NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--INDEXING
CREATE INDEX idx_inventory_name ON inventory USING GIN (to_tsvector ('english', name));

CREATE INDEX idx_inventory_description ON inventory USING GIN (
    to_tsvector ('english', description)
);