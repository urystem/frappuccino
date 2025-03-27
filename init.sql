--CREATE TABLE
--inventory

CREATE TYPE uints AS ENUM ('g','ml','pcs');
CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    quantity FLOAT NOT NULL CHECK (quantity >=0),
    reorder_level FLOAT NOT NULL CHECK (reorder_level > 0),
    unit uints NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);

CREATE TYPE reason_of_inventory_transaction AS ENUM ('restock', 'usage', 'cancelled', 'annul');

CREATE TABLE inventory_transactions(--history
    id SERIAL PRIMARY KEY,
    inventory_id INT NOT NULL REFERENCES inventory(id) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason reason_of_inventory_transaction NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP--NOW()
);

--menu
CREATE TABLE menu_items(
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    tags TEXT[],
    allergen TEXT[],
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0)
    --inventories TEXT[] NOT NULL CHECK (array_length(allergens, 1) > 0) --cardinality(allergens)>0
);

CREATE TABLE menu_item_ingredients(
    product_id INT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    inventory_id INT NOT NULL REFERENCES inventory(id),
    quantity FLOAT NOT NULL CHECK (quantity >0)
);

CREATE TABLE price_history (--need trigger
    id SERIAL PRIMARY KEY, 
    product_id INT NOT NULL REFERENCES menu_items(id) ON DELETE CASCADE,
    old_price DECIMAL(10,2) NOT NULL CHECK (old_price >= 0),
    new_price DECIMAL(10,2) NOT NULL CHECK (new_price >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP--NOW()
);


--order
CREATE TYPE order_status AS ENUM ('processing', 'accepted', 'rejected');

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    customer_name VARCHAR(100) NOT NULL,
    status order_status NOT NULL,
    allergen TEXT[],
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,--NOW()
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP--NOW()
);

CREATE TABLE order_items (
    order_id INT NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES menu_items(id),
    quantity INT NOT NULL CHECK (quantity > 0)
);

CREATE TABLE order_status_history(
    id SERIAL PRIMARY KEY, 
    order_id INT  NOT NULL REFERENCES orders(id) ON DELETE CASCADE,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP--NOW()
);


--INSERT TO THE TABLE
