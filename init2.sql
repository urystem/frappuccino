/*CREATE TABLE inventory_categories (
    id SERIAL PRIMARY KEY,       -- Уникальный ID категории
    name VARCHAR(48) NOT NULL,  -- Название категории (например, "Овощи")
    description TEXT             -- Описание категории (может быть NULL)
);*/

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    --category_id INT NOT NULL REFERENCES inventory_categories(id) ON DELETE CASCADE,
    description TEXT,
    quantity FLOAT NOT NULL CHECK (quantity >=0),
    unit VARCHAR(16) NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0),
);

/*
CREATE TABLE IF NOT EXISTS  inventories (
    ingredient_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    category_id INT REFERENCES inventory_categories(id) ON DELETE CASCADE,
    quantity FLOAT NOT NULL CHECK (quantity >=0),
    unit VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);
*/

/*
DROP TABLE IF EXISTS inventories;
CREATE TABLE inventories (
    ingredient_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL /*UNIQUE*/,
    category_id INT REFERENCES inventory_categories(id) ON DELETE CASCADE,
    quantity FLOAT NOT NULL CHECK (quantity >=0),
    unit VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);
*/
/*
INSERT INTO inventories (name, quantity, unit, price) VALUES
('Sugar', 100.0, 'kg', 1.50),
('Flour', 50.0, 'kg', 0.80),
('Salt', 20.0, 'kg', 0.50),
('Olive Oil', 10.0, 'L', 5.00),
('Butter', 30.0, 'kg', 3.20),
('Eggs', 200.0, 'pcs', 0.10),
('Milk', 50.0, 'L', 1.20),
('Chicken Breast', 10.0, 'kg', 7.50),
('Rice', 100.0, 'kg', 1.10),
('Cheese', 20.0, 'kg', 4.00);
*/

CREATE TABLE inventory_transactions(--history
    id SERIAL PRIMARY KEY, 
    ingredient_id INT NOT NULL REFERENCES inventories(id) ON DELETE CASCADE, 
    quantity_change FLOAT NOT NULL, 
    reason TEXT,
    changed_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
/*
CREATE TABLE menu_categories (
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT
);
*/
CREATE TABLE menu_items(
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0),
   --category_id INT NOT NULL REFERENCES menu_categories(id) ON DELETE CASCADE,
    allergens TEXT[]  
);


CREATE TABLE menu_item_ingredients (
    ingredient_id INT NOT NULL CHECK (ingredient_id > 0),
    quantity FLOAT NOT NULL CHECK (quantity > 0)
);