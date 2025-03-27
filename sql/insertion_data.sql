-- Insert into orders table
INSERT INTO orders (customer_name, status, total) VALUES
('Alice Johnson', 'pending', 250),
('Bob Smith', 'confirmed', 400),
('Charlie Brown', 'in progress', 150),
('Diana Prince', 'completed', 500),
('Ethan Hunt', 'cancelled', 300),
('Frank Castle', 'pending', 200),
('Grace Hopper', 'confirmed', 350),
('Hank Pym', 'in progress', 275),
('Ivy League', 'completed', 600),
('Jack Sparrow', 'rejected', 450),
('Kim Possible', 'pending', 320),
('Luke Skywalker', 'confirmed', 500),
('Monica Geller', 'in progress', 280),
('Nathan Drake', 'completed', 700),
('Olivia Benson', 'cancelled', 380),
('Peter Parker', 'pending', 150),
('Quentin Tarantino', 'confirmed', 550),
('Rachel Green', 'in progress', 420),
('Steve Rogers', 'completed', 750),
('Tony Stark', 'cancelled', 900);

-- Insert into menu_items table
INSERT INTO menu_items (name, description, details, price, allergens) VALUES
('Margherita Pizza', 'Classic pizza with tomato and mozzarella', '{"ingredients": ["tomato", "mozzarella"], "dietary": "vegetarian"}', 12.99, '{"dairy"}'),
('Veggie Burger', 'Grilled vegetable patty with lettuce and tomato', '{"ingredients": ["lettuce", "tomato", "patty"], "dietary": "vegetarian"}', 9.99, '{"soy", "gluten"}'),
('Grilled Chicken Salad', 'Salad with grilled chicken and vinaigrette', '{"ingredients": ["chicken", "lettuce", "vinaigrette"]}', 11.99, '{"none"}'),
('Pepperoni Pizza', 'Pizza with spicy pepperoni', '{"ingredients": ["tomato", "cheese", "pepperoni"]}', 14.99, '{"dairy", "gluten"}'),
('Caesar Salad', 'Salad with romaine, parmesan, and croutons', '{"ingredients": ["lettuce", "parmesan", "croutons"]}', 10.99, '{"dairy", "gluten"}');

-- Insert into inventory_items table
INSERT INTO inventory_items (name, quantity, unit, allergens, extra_info) VALUES
('Tomato', 100, 'kg', '{}', '{"supplier": "Fresh Farms", "expiry_date": "2025-05-10"}'),
('Mozzarella', 50, 'kg', '{"dairy"}', '{"supplier": "Cheese Co", "expiry_date": "2025-04-20"}'),
('Chicken Breast', 80, 'kg', '{}', '{"supplier": "Poultry Inc", "expiry_date": "2025-03-15"}');

-- Insert into order_items table
INSERT INTO order_items (order_id, menu_item_id, quantity) VALUES
(1, 1, 2),
(2, 2, 1),
(3, 3, 1);

-- Insert into menu_item_ingredients table
INSERT INTO menu_item_ingredients (menu_item_id, inventory_item_id, quantity) VALUES
(1, 1, 5),
(1, 2, 3),
(3, 3, 2);

-- Insert into inventory_transactions table
INSERT INTO inventory_transactions (inventory_item_id, quantity_change, transaction_date) VALUES
(1, -5, '2025-03-01 10:00:00'),
(2, -3, '2025-03-01 11:00:00'),
(3, -2, '2025-03-01 12:00:00');

-- Insert into order_status_history table
INSERT INTO order_status_history (order_id, status, updated_at) VALUES
(1, 'confirmed', '2025-03-02 10:00:00'),
(2, 'in progress', '2025-03-02 11:00:00'),
(3, 'completed', '2025-03-02 12:00:00');

-- Insert into price_history table
INSERT INTO price_history (menu_item_id, old_price, new_price, changed_at) VALUES
(1, 11.99, 12.99, '2025-02-28 10:00:00'),
(2, 8.99, 9.99, '2025-02-28 11:00:00'),
(3, 10.99, 11.99, '2025-02-28 12:00:00');
