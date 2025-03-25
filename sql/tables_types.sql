ALTER OR CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled', 'rejected');

ALTER OR CREATE TABLE orders (
    orders_id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    total INTEGER NOT NULL
);

CREATE TABLE menu_items (
    menu_items_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price REAL NOT NULL
);

CREATE TABLE inventory_items (
    inventory_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    unit TEXT NOT NULL
);

CREATE TABLE order_items (
    order_items_id SERIAL PRIMARY KEY,
    orders_id INT REFERENCES orders(orders_id) ON DELETE CASCADE,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

CREATE TABLE menu_item_ingredients (
    menu_item_ingredients_id SERIAL PRIMARY KEY,
    menu_items_id INT REFERENCES menu_item(menu_items_id) ON DELETE CASCADE,
    inventory_items_id INT REFERENCES inventory_items(inventory_items_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    inventory_item_id INT REFERENCES inventory_item(inventory_item_id) ON DELETE CASCADE,
    quantity_change INT NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE order_status_history (
    status_history_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    status oreder_status NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_item(menu_item_id) ON DELETE CASCADE,
    old_price INT NOT NULL,
    new_price INT NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
