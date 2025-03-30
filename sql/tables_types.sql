-- Order Status Enum
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled', 'rejected');

-- Orders Table
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    total REAL NOT NULL,
    search_vector TSVECTOR
);

-- Menu Items Table (Added allergens array and JSONB details)
CREATE TABLE menu_items (
    menu_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    details JSONB, -- Flexible data storage (e.g., {"ingredients": ["cheese", "tomato"], "dietary": "vegetarian"})
    price REAL NOT NULL,
    allergens TEXT[] DEFAULT '{}', -- Stores multiple allergens like {"nuts", "gluten"}
    search_vector TSVECTOR
);

-- Inventory Items Table (Added allergens array and JSONB extra_info)
CREATE TABLE inventory_items (
    inventory_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    quantity INT NOT NULL,
    unit TEXT NOT NULL,
    allergens TEXT[] DEFAULT '{}', -- Example: {"dairy", "soy"}
    extra_info JSONB -- Flexible metadata (e.g., {"supplier": "FarmFresh", "expiry_date": "2025-01-01"})
);

-- Order Items Table
CREATE TABLE order_items (
    order_item_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

-- Menu Item Ingredients Table
CREATE TABLE menu_item_ingredients (
    menu_item_ingredient_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    inventory_item_id INT REFERENCES inventory_items(inventory_item_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

-- Inventory Transactions Table
CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    inventory_item_id INT REFERENCES inventory_items(inventory_item_id) ON DELETE CASCADE,
    quantity_change INT NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Order Status History Table
CREATE TABLE order_status_history (
    status_history_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Price History Table
CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_items(menu_item_id) ON DELETE CASCADE,
    old_price REAL NOT NULL,
    new_price REAL NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE INDEX idx_order_items_menu_item_id ON order_items (menu_item_id);
CREATE INDEX idx_menu_items_menu_item_id ON menu_items (menu_item_id);
CREATE INDEX idx_menu_items_price ON menu_items(price);
CREATE INDEX idx_orders_total ON orders(total);
CREATE INDEX idx_menu_items_search ON menu_items USING gin(search_vector);
CREATE INDEX idx_orders_search ON orders USING gin(search_vector);

CREATE VIEW popular_menu_items AS
SELECT 
    mi.name, 
    COUNT(oi.menu_item_id) AS order_count, 
    MAX(osh.updated_at) AS last_updated_at
FROM order_items oi
JOIN menu_items mi ON oi.menu_item_id = mi.menu_item_id
JOIN order_status_history osh ON osh.order_id = oi.order_id
JOIN orders o ON oi.order_id = o.order_id
GROUP BY mi.name;


CREATE OR REPLACE FUNCTION update_menu_items_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('english', NEW.name || ' ' || NEW.description);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Trigger to Automatically Update Search Vector
CREATE TRIGGER menu_items_search_vector_trigger
BEFORE INSERT OR UPDATE ON menu_items
FOR EACH ROW EXECUTE FUNCTION update_menu_items_search_vector();