INSERT INTO customers (name, email, sex) VALUES
('John Doe', 'john@example.com', 'man'),
('Jane Smith', 'jane@example.com', 'woman'),
('Alex Johnson', 'alex@example.com', 'undefined');

INSERT INTO menu_item (name, description, price) VALUES
('Margherita Pizza', 'Classic pizza with tomato sauce and mozzarella', 1200),
('Cheeseburger', 'Beef patty with cheese, lettuce, and tomato', 900),
('Caesar Salad', 'Romaine lettuce, croutons, parmesan, Caesar dressing', 700);

INSERT INTO inventory_item (name, quantity, unit) VALUES
('Flour', 50, 'kg'),
('Cheese', 30, 'kg'),
('Lettuce', 20, 'kg'),
('Beef Patty', 40, 'pcs'),
('Tomato', 25, 'kg');

INSERT INTO orders (customer_id, status, total) VALUES
(1, 'confirmed', 2100),  -- John Doe's order
(2, 'pending', 900),     -- Jane Smith's order
(3, 'completed', 1200);  -- Alex Johnson's order

INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES
(1, 1, 1),  -- 1 Margherita Pizza for John
(1, 2, 1),  -- 1 Cheeseburger for John
(2, 2, 1),  -- 1 Cheeseburger for Jane
(3, 1, 1);  -- 1 Margherita Pizza for Alex

INSERT INTO menu_item_ingredients (menu_item_id, inventory_item_id, quantity) VALUES
(1, 1, 0.3),  -- Margherita Pizza uses 0.3 kg of Flour
(1, 2, 0.2),  -- Margherita Pizza uses 0.2 kg of Cheese
(2, 4, 1),    -- Cheeseburger uses 1 Beef Patty
(2, 5, 0.1);  -- Cheeseburger uses 0.1 kg of Tomato

INSERT INTO inventory_transactions (inventory_item_id, quantity_change) VALUES
(1, -1),  -- Used 1 kg of Flour
(2, -0.5), -- Used 0.5 kg of Cheese
(4, -2);  -- Used 2 Beef Patties

INSERT INTO order_status_history (order_id, status) VALUES
(1, 'pending'),
(1, 'confirmed'),
(2, 'pending'),
(3, 'completed');

INSERT INTO price_history (menu_item_id, old_price, new_price) VALUES
(1, 1100, 1200),  -- Price increased for Margherita Pizza
(2, 850, 900);    -- Price increased for Cheeseburger
