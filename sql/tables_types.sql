ALTER TYPE sex AS ENUM ('man', 'woman', 'undefined');
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled', 'refused');

CREATE TABLE orders (
    order_id SERIAL PRIMARY KEY,
    customer_id INT NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'pending',
    order_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    total INTEGER NOT NULL
);

CREATE TABLE menu_item (
    menu_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT,
    price INT NOT NULL
);

CREATE TABLE inventory_item (
    inventory_item_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    quantity INT NOT NULL,
    unit TEXT NOT NULL
);

CREATE TABLE customers (
    customer_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    sex sex NOT NULL DEFAULT 'undefined',
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);


CREATE TABLE order_items (
    order_items_id SERIAL PRIMARY KEY,
    order_id INT REFERENCES orders(order_id) ON DELETE CASCADE,
    menu_item_id INT REFERENCES menu_item(menu_item_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

CREATE TABLE menu_item_ingredients (
    menu_item_ingredients_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_item(menu_item_id) ON DELETE CASCADE,
    inventory_item_id INT REFERENCES inventory_item(inventory_item_id) ON DELETE CASCADE,
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
    status VARCHAR(50) NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_item_id INT REFERENCES menu_item(menu_item_id) ON DELETE CASCADE,
    old_price INT NOT NULL,
    new_price INT NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);
