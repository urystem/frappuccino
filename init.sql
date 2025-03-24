CREATE TABLE inventories (
    ingredient_id SERIAL PRIMARY KEY,
    name VARCHAR(100) NOT NULL /*UNIQUE*/,
    quantity FLOAT NOT NULL CHECK (quantity >=0),
    unit VARCHAR(20) NOT NULL,
    price DECIMAL(10,2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);

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
