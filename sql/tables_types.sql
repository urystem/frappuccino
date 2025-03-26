-- Order Status Enum
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled', 'rejected');

-- Orders Table
CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    total INTEGER NOT NULL
);

-- Menu Items Table (Added allergens array and JSONB details)
CREATE TABLE menu_items (
    menu_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    details JSONB, -- Flexible data storage (e.g., {"ingredients": ["cheese", "tomato"], "dietary": "vegetarian"})
    price REAL NOT NULL,
    allergens TEXT[] DEFAULT '{}' -- Stores multiple allergens like {"nuts", "gluten"}
);

-- Inventory Items Table (Added allergens array and JSONB extra_info)
CREATE TABLE inventory_items (
    inventory_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
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
