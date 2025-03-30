-- Order Status Enum
CREATE TYPE order_status AS ENUM ('pending', 'confirmed', 'in progress', 'completed', 'cancelled', 'rejected');

-- Orders Table
CREATE TABLE orders (
    orders_id SERIAL PRIMARY KEY,
    customer_name TEXT NOT NULL,
    status order_status NOT NULL DEFAULT 'pending',
    total REAL NOT NULL
);

-- Menu Items Table (Added allergens array and JSONB details)
CREATE TABLE menu_items (
    menu_items_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL,
    details JSONB, -- Flexible data storage (e.g., {"ingredients": ["cheese", "tomato"], "dietary": "vegetarian"})
    price REAL NOT NULL,
    allergens TEXT[] DEFAULT '{}' -- Stores multiple allergens like {"nuts", "gluten"}
);

-- Inventory Items Table (Added allergens array and JSONB extra_info)
CREATE TABLE inventory_items (
    inventory_items_id SERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    quantity INT NOT NULL,
    unit TEXT NOT NULL,
    allergens TEXT[] DEFAULT '{}', -- Example: {"dairy", "soy"}
    extra_info JSONB -- Flexible metadata (e.g., {"supplier": "FarmFresh", "expiry_date": "2025-01-01"})
);

-- Order Items Table
CREATE TABLE order_items (
    order_items_id SERIAL PRIMARY KEY,
    orders_id INT REFERENCES orders(orders_id) ON DELETE CASCADE,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

-- Menu Item Ingredients Table
CREATE TABLE menu_item_ingredients (
    menu_item_ingredients_id SERIAL PRIMARY KEY,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ON DELETE CASCADE,
    inventory_items_id INT REFERENCES inventory_items(inventory_items_id) ON DELETE CASCADE,
    quantity INT NOT NULL
);

-- Inventory Transactions Table
CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    inventory_items_id INT REFERENCES inventory_items(inventory_items_id) ON DELETE CASCADE,
    quantity_change INT NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Order Status History Table
CREATE TABLE order_status_history (
    status_history_id SERIAL PRIMARY KEY,
    orders_id INT REFERENCES orders(orders_id) ON DELETE CASCADE,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Price History Table
CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ON DELETE CASCADE,
    old_price REAL NOT NULL,
    new_price REAL NOT NULL,
    changed_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE VIEW popular_menu_items AS
SELECT 
    mi.name, 
    COUNT(oi.menu_items_id) AS order_count, 
    MAX(osh.updated_at) AS last_updated_at
FROM order_items oi
JOIN menu_items mi ON oi.menu_items_id = mi.menu_items_id
JOIN order_status_history osh ON osh.orders_id = oi.orders_id
JOIN orders o ON oi.orders_id = o.orders_id
GROUP BY mi.name;

-- Indexes
CREATE INDEX search_idx ON menu_items USING gin(to_tsvector('english', name || ' ' || description || ' ' || details::text));

CREATE OR REPLACE FUNCTION search_all(
    query_text TEXT,
    search_filter TEXT DEFAULT 'all',
    min_price REAL DEFAULT NULL,
    max_price REAL DEFAULT NULL
)
RETURNS JSONB AS $$
DECLARE
    result JSONB;
    menu_items_result JSONB;
    orders_result JSONB;
BEGIN
    -- Initialize result object
    result := '{}'::JSONB;
    
    -- Search menu items if requested
    IF search_filter = 'all' OR search_filter LIKE '%menu%' THEN
        WITH menu_search AS (
            SELECT 
                mi.menu_items_id as id,
                mi.name,
                mi.description,
                mi.price,
                ts_rank_cd(
                    to_tsvector('english', mi.name || ' ' || mi.description || ' ' || COALESCE(mi.details::text, '')),
                    plainto_tsquery('english', query_text)
                ) as relevance
            FROM menu_items mi
            WHERE 
                (to_tsvector('english', mi.name || ' ' || mi.description || ' ' || COALESCE(mi.details::text, '')) 
                @@ plainto_tsquery('english', query_text))
                AND (min_price IS NULL OR mi.price >= min_price)
                AND (max_price IS NULL OR mi.price <= max_price)
            ORDER BY relevance DESC
        )
        SELECT jsonb_agg(
            jsonb_build_object(
                'id', id,
                'name', name,
                'description', description,
                'price', price,
                'relevance', relevance
            )
        ) INTO menu_items_result
        FROM menu_search;
        
        IF menu_items_result IS NULL THEN
            menu_items_result := '[]'::JSONB;
        END IF;
        
        result := result || jsonb_build_object('menu_items', menu_items_result);
    END IF;
    
    -- Search orders if requested
    IF search_filter = 'all' OR search_filter LIKE '%orders%' THEN
        WITH order_items_agg AS (
            SELECT 
                oi.orders_id,
                jsonb_agg(mi.name) AS items
            FROM order_items oi
            JOIN menu_items mi ON oi.menu_items_id = mi.menu_items_id
            GROUP BY oi.orders_id
        ),
        order_search AS (
            SELECT 
                o.orders_id as id,
                o.customer_name,
                oia.items,
                o.total,
                ts_rank_cd(
                    to_tsvector('english', o.customer_name || ' ' || COALESCE(oia.items::text, '')),
                    plainto_tsquery('english', query_text)
                ) as relevance
            FROM orders o
            LEFT JOIN order_items_agg oia ON o.orders_id = oia.orders_id
            WHERE 
                (to_tsvector('english', o.customer_name || ' ' || COALESCE(oia.items::text, '')) 
                @@ plainto_tsquery('english', query_text))
                AND (min_price IS NULL OR o.total >= min_price)
                AND (max_price IS NULL OR o.total <= max_price)
            ORDER BY relevance DESC
        )
        SELECT jsonb_agg(
            jsonb_build_object(
                'id', id,
                'customer_name', customer_name,
                'items', items,
                'total', total,
                'relevance', relevance
            )
        ) INTO orders_result
        FROM order_search;
        
        IF orders_result IS NULL THEN
            orders_result := '[]'::JSONB;
        END IF;
        
        result := result || jsonb_build_object('orders', orders_result);
    END IF;
    
    -- Add total matches count
    result := result || jsonb_build_object(
        'total_matches',
        COALESCE((SELECT COUNT(*) FROM jsonb_array_elements(result->'menu_items')), 0) +
        COALESCE((SELECT COUNT(*) FROM jsonb_array_elements(result->'orders')), 0)
    );
    
    RETURN result;
END;
$$ LANGUAGE plpgsql;

INSERT INTO orders (customer_name, status, total) VALUES
('John Smith', 'completed', 12.50),
('Emily Johnson', 'in progress', 8.75),
('Michael Brown', 'confirmed', 15.20),
('Sarah Davis', 'pending', 9.95),
('David Wilson', 'completed', 22.30),
('Jessica Miller', 'cancelled', 7.50),
('Robert Taylor', 'completed', 18.45),
('Jennifer Anderson', 'in progress', 11.25),
('Thomas Martinez', 'completed', 14.80),
('Lisa Robinson', 'confirmed', 6.95),
('Daniel White', 'pending', 19.60),
('Amanda Clark', 'completed', 13.75),
('James Rodriguez', 'rejected', 10.50),
('Michelle Lewis', 'completed', 16.90),
('Christopher Lee', 'in progress', 8.25),
('Ashley Walker', 'completed', 21.40),
('Matthew Hall', 'confirmed', 7.80),
('Elizabeth Young', 'completed', 12.95),
('Joshua Allen', 'pending', 9.30),
('Stephanie King', 'completed', 17.65),
('Andrew Wright', 'in progress', 10.75),
('Nicole Scott', 'completed', 14.20),
('Kevin Green', 'confirmed', 6.50),
('Rebecca Baker', 'completed', 20.15),
('Ryan Adams', 'pending', 11.90),
('Laura Nelson', 'completed', 15.75),
('Jason Hill', 'in progress', 9.45),
('Megan Carter', 'completed', 13.30),
('Justin Mitchell', 'confirmed', 7.20),
('Samantha Perez', 'completed', 18.95);

INSERT INTO menu_items (name, description, details, price, allergens) VALUES
('Espresso', 'Strong black coffee', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 2.50, '{}'),
('Cappuccino', 'Espresso with steamed milk foam', '{"ingredients": ["coffee beans", "milk", "foam"], "dietary": "vegetarian"}', 3.75, '{"dairy"}'),
('Latte', 'Espresso with steamed milk', '{"ingredients": ["coffee beans", "milk"], "dietary": "vegetarian"}', 4.00, '{"dairy"}'),
('Americano', 'Espresso with hot water', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 3.00, '{}'),
('Mocha', 'Espresso with chocolate and steamed milk', '{"ingredients": ["coffee beans", "milk", "chocolate"], "dietary": "vegetarian"}', 4.50, '{"dairy"}'),
('Flat White', 'Espresso with microfoam', '{"ingredients": ["coffee beans", "milk"], "dietary": "vegetarian"}', 4.25, '{"dairy"}'),
('Macchiato', 'Espresso with a dollop of foam', '{"ingredients": ["coffee beans", "milk foam"], "dietary": "vegetarian"}', 3.50, '{"dairy"}'),
('Cold Brew', 'Slow-steeped cold coffee', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 4.00, '{}'),
('Iced Latte', 'Latte served over ice', '{"ingredients": ["coffee beans", "milk", "ice"], "dietary": "vegetarian"}', 4.50, '{"dairy"}'),
('Chai Latte', 'Spiced tea with steamed milk', '{"ingredients": ["tea", "milk", "spices"], "dietary": "vegetarian"}', 4.25, '{"dairy"}'),
('Hot Chocolate', 'Steamed milk with chocolate', '{"ingredients": ["milk", "chocolate"], "dietary": "vegetarian"}', 3.75, '{"dairy"}'),
('Matcha Latte', 'Green tea powder with steamed milk', '{"ingredients": ["matcha", "milk"], "dietary": "vegetarian"}', 4.75, '{"dairy"}'),
('Turkish Coffee', 'Strong unfiltered coffee', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 3.25, '{}'),
('Affogato', 'Espresso poured over vanilla ice cream', '{"ingredients": ["coffee beans", "ice cream"], "dietary": "vegetarian"}', 5.00, '{"dairy"}'),
('Red Eye', 'Drip coffee with a shot of espresso', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 3.75, '{}'),
('Cortado', 'Espresso with equal parts warm milk', '{"ingredients": ["coffee beans", "milk"], "dietary": "vegetarian"}', 3.50, '{"dairy"}'),
('Irish Coffee', 'Coffee with whiskey and cream', '{"ingredients": ["coffee beans", "whiskey", "cream"], "dietary": "vegetarian"}', 6.50, '{"dairy"}'),
('Vietnamese Iced Coffee', 'Strong coffee with condensed milk', '{"ingredients": ["coffee beans", "condensed milk"], "dietary": "vegetarian"}', 4.75, '{"dairy"}'),
('Frappuccino', 'Blended iced coffee drink', '{"ingredients": ["coffee beans", "milk", "ice"], "dietary": "vegetarian"}', 5.25, '{"dairy"}'),
('Cafe au Lait', 'Coffee with steamed milk', '{"ingredients": ["coffee beans", "milk"], "dietary": "vegetarian"}', 3.75, '{"dairy"}'),
('Caramel Macchiato', 'Espresso with vanilla, milk and caramel', '{"ingredients": ["coffee beans", "milk", "vanilla", "caramel"], "dietary": "vegetarian"}', 5.00, '{"dairy"}'),
('Pumpkin Spice Latte', 'Espresso with pumpkin spice and milk', '{"ingredients": ["coffee beans", "milk", "pumpkin spice"], "dietary": "vegetarian"}', 5.50, '{"dairy"}'),
('Hazelnut Latte', 'Espresso with hazelnut syrup and milk', '{"ingredients": ["coffee beans", "milk", "hazelnut syrup"], "dietary": "vegetarian"}', 4.75, '{"dairy", "nuts"}'),
('Vanilla Latte', 'Espresso with vanilla syrup and milk', '{"ingredients": ["coffee beans", "milk", "vanilla syrup"], "dietary": "vegetarian"}', 4.50, '{"dairy"}'),
('Almond Milk Latte', 'Espresso with almond milk', '{"ingredients": ["coffee beans", "almond milk"], "dietary": "vegan"}', 4.75, '{"nuts"}'),
('Soy Latte', 'Espresso with soy milk', '{"ingredients": ["coffee beans", "soy milk"], "dietary": "vegan"}', 4.50, '{"soy"}'),
('Oat Milk Latte', 'Espresso with oat milk', '{"ingredients": ["coffee beans", "oat milk"], "dietary": "vegan"}', 4.50, '{}'),
('Peppermint Mocha', 'Espresso with peppermint, chocolate and milk', '{"ingredients": ["coffee beans", "milk", "chocolate", "peppermint"], "dietary": "vegetarian"}', 5.25, '{"dairy"}'),
('Coconut Latte', 'Espresso with coconut milk', '{"ingredients": ["coffee beans", "coconut milk"], "dietary": "vegan"}', 4.75, '{}'),
('Ristretto', 'Short shot of espresso', '{"ingredients": ["coffee beans", "water"], "dietary": "vegan"}', 2.75, '{}');

INSERT INTO inventory_items (name, quantity, unit, allergens, extra_info) VALUES
('Coffee Beans', 5000, 'g', '{}', '{"supplier": "Bean Masters", "roast_date": "2023-11-01"}'),
('Milk', 20000, 'ml', '{"dairy"}', '{"supplier": "Fresh Farms Dairy", "expiry_date": "2023-11-15"}'),
('Sugar', 10000, 'g', '{}', '{"supplier": "Sweet Co", "package_size": "10kg"}'),
('Vanilla Syrup', 5000, 'ml', '{}', '{"supplier": "Flavor Art", "batch": "VA2301"}'),
('Chocolate Syrup', 3000, 'ml', '{"dairy"}', '{"supplier": "Cocoa Delight", "expiry_date": "2024-05-01"}'),
('Caramel Syrup', 2500, 'ml', '{}', '{"supplier": "Golden Sweet", "batch": "CR2302"}'),
('Whipped Cream', 1500, 'ml', '{"dairy"}', '{"supplier": "Cloudy Cream", "expiry_date": "2023-11-10"}'),
('Ice', 0, 'cubes', '{}', '{"note": "Made in-house"}'),
('Cinnamon', 500, 'g', '{}', '{"supplier": "Spice World", "origin": "Sri Lanka"}'),
('Nutmeg', 300, 'g', '{}', '{"supplier": "Spice World", "origin": "Indonesia"}'),
('Almond Milk', 8000, 'ml', '{"nuts"}', '{"supplier": "Nutty Goodness", "expiry_date": "2023-12-01"}'),
('Soy Milk', 7500, 'ml', '{"soy"}', '{"supplier": "Bean Liquid", "expiry_date": "2023-12-15"}'),
('Oat Milk', 9000, 'ml', '{}', '{"supplier": "Grainy Liquid", "expiry_date": "2023-12-10"}'),
('Coconut Milk', 6000, 'ml', '{}', '{"supplier": "Tropical Taste", "expiry_date": "2024-01-05"}'),
('Matcha Powder', 1000, 'g', '{}', '{"supplier": "Green Tea Co", "grade": "Ceremonial"}'),
('Chai Tea', 1200, 'g', '{}', '{"supplier": "Spice Route", "blend": "Masala"}'),
('Espresso Cups', 100, 'pieces', '{}', '{"supplier": "Cup Masters", "material": "Porcelain"}'),
('Coffee Mugs', 80, 'pieces', '{}', '{"supplier": "Mug World", "capacity": "350ml"}'),
('Lids', 500, 'pieces', '{}', '{"supplier": "Cover It", "size": "Medium"}'),
('Straws', 1000, 'pieces', '{}', '{"supplier": "Sip Well", "material": "Paper"}'),
('Napkins', 5000, 'pieces', '{}', '{"supplier": "Clean Touch", "ply": "2"}'),
('Paper Cups', 400, 'pieces', '{}', '{"supplier": "Eco Serve", "size": "12oz"}'),
('Plastic Spoons', 300, 'pieces', '{}', '{"supplier": "Utensil Pro", "material": "BPA-free"}'),
('Whiskey', 5000, 'ml', '{}', '{"supplier": "Spirit Co", "alcohol_content": "40%"}'),
('Condensed Milk', 4000, 'ml', '{"dairy"}', '{"supplier": "Sweet Cream", "expiry_date": "2024-03-01"}'),
('Pumpkin Spice', 800, 'g', '{}', '{"supplier": "Autumn Flavors", "composition": "cinnamon, nutmeg, ginger, cloves"}'),
('Hazelnut Syrup', 2000, 'ml', '{"nuts"}', '{"supplier": "Nutty Flavors", "batch": "HZ2303"}'),
('Peppermint Syrup', 1800, 'ml', '{}', '{"supplier": "Minty Fresh", "seasonal": "Winter"}'),
('Ice Cream', 5000, 'ml', '{"dairy"}', '{"supplier": "Creamy Delight", "flavor": "Vanilla", "expiry_date": "2023-12-20"}'),
('Water', 50000, 'ml', '{}', '{"supplier": "Pure Springs", "type": "Filtered"}');

INSERT INTO order_items (orders_id, menu_items_id, quantity) VALUES
(1, 3, 1), (1, 5, 1), (1, 12, 1),
(2, 1, 2), (2, 8, 1),
(3, 4, 1), (3, 6, 1), (3, 20, 1),
(4, 2, 1), (4, 11, 1),
(5, 7, 2), (5, 9, 1), (5, 15, 1),
(6, 10, 1),
(7, 13, 1), (7, 14, 1),
(8, 16, 2),
(9, 17, 1), (9, 18, 1),
(10, 19, 1),
(11, 21, 1), (11, 22, 1),
(12, 23, 1), (12, 24, 1),
(13, 25, 1),
(14, 26, 1), (14, 27, 1),
(15, 28, 1),
(16, 29, 1), (16, 30, 1),
(17, 1, 1), (17, 3, 1),
(18, 5, 2),
(19, 7, 1), (19, 9, 1),
(20, 11, 1), (20, 13, 1),
(21, 15, 1), (21, 17, 1),
(22, 19, 1), (22, 21, 1),
(23, 23, 1),
(24, 25, 1), (24, 27, 1),
(25, 29, 1),
(26, 2, 1), (26, 4, 1),
(27, 6, 1), (27, 8, 1),
(28, 10, 1), (28, 12, 1),
(29, 14, 1),
(30, 16, 1), (30, 18, 1);


INSERT INTO menu_item_ingredients (menu_items_id, inventory_items_id, quantity) VALUES
(1, 1, 10), (1, 30, 30),
(2, 1, 10), (2, 2, 150), (2, 30, 20),
(3, 1, 10), (3, 2, 200), (3, 30, 10),
(4, 1, 10), (4, 30, 200),
(5, 1, 10), (5, 2, 150), (5, 5, 20), (5, 30, 10),
(6, 1, 10), (6, 2, 180), (6, 30, 10),
(7, 1, 10), (7, 2, 30), (7, 30, 10),
(8, 1, 15), (8, 30, 200), (8, 8, 10),
(9, 1, 10), (9, 2, 200), (9, 8, 10),
(10, 16, 5), (10, 2, 200), (10, 30, 10),
(11, 2, 250), (11, 5, 30), (11, 30, 10),
(12, 15, 5), (12, 2, 200), (12, 30, 10),
(13, 1, 15), (13, 30, 100),
(14, 1, 10), (14, 29, 50),
(15, 1, 10), (15, 30, 200),
(16, 1, 10), (16, 2, 100), (16, 30, 10),
(17, 1, 10), (17, 24, 30), (17, 7, 30), (17, 30, 150),
(18, 1, 15), (18, 25, 50), (18, 30, 50),
(19, 1, 15), (19, 2, 150), (19, 8, 10), (19, 5, 20),
(20, 1, 10), (20, 2, 200), (20, 30, 10),
(21, 1, 10), (21, 2, 150), (21, 4, 15), (21, 6, 15), (21, 30, 10),
(22, 1, 10), (22, 2, 150), (22, 26, 10), (22, 30, 10),
(23, 1, 10), (23, 2, 150), (23, 27, 15), (23, 30, 10),
(24, 1, 10), (24, 2, 150), (24, 4, 15), (24, 30, 10),
(25, 1, 10), (25, 11, 200), (25, 30, 10),
(26, 1, 10), (26, 12, 200), (26, 30, 10),
(27, 1, 10), (27, 13, 200), (27, 30, 10),
(28, 1, 10), (28, 14, 200), (28, 30, 10),
(29, 1, 10), (29, 2, 150), (29, 5, 15), (29, 28, 15), (29, 30, 10),
(30, 1, 7), (30, 30, 30);

INSERT INTO inventory_transactions (inventory_items_id, quantity_change, transaction_date) VALUES
(1, -50, '2023-10-01 08:15:00'), (1, 1000, '2023-10-02 09:00:00'),
(2, -500, '2023-10-01 08:20:00'), (2, 5000, '2023-10-03 10:00:00'),
(3, -100, '2023-10-01 08:25:00'), (3, 2000, '2023-10-04 11:00:00'),
(4, -50, '2023-10-01 08:30:00'), (4, 1000, '2023-10-05 12:00:00'),
(5, -30, '2023-10-01 08:35:00'), (5, 500, '2023-10-06 13:00:00'),
(6, -25, '2023-10-01 08:40:00'), (6, 500, '2023-10-07 14:00:00'),
(7, -15, '2023-10-01 08:45:00'), (7, 300, '2023-10-08 15:00:00'),
(8, -100, '2023-10-01 08:50:00'), (8, 0, '2023-10-09 16:00:00'),
(9, -5, '2023-10-01 08:55:00'), (9, 100, '2023-10-10 17:00:00'),
(10, -3, '2023-10-01 09:00:00'), (10, 50, '2023-10-11 18:00:00'),
(11, -80, '2023-10-01 09:05:00'), (11, 1000, '2023-10-12 19:00:00'),
(12, -75, '2023-10-01 09:10:00'), (12, 1000, '2023-10-13 20:00:00'),
(13, -90, '2023-10-01 09:15:00'), (13, 1000, '2023-10-14 21:00:00'),
(14, -60, '2023-10-01 09:20:00'), (14, 800, '2023-10-15 22:00:00'),
(15, -10, '2023-10-01 09:25:00'), (15, 200, '2023-10-16 23:00:00');

INSERT INTO order_status_history (orders_id, status, updated_at) VALUES
(1, 'pending', '2023-10-01 08:00:00'), (1, 'confirmed', '2023-10-01 08:02:00'), (1, 'in progress', '2023-10-01 08:05:00'), (1, 'completed', '2023-10-01 08:15:00'),
(2, 'pending', '2023-10-01 08:10:00'), (2, 'confirmed', '2023-10-01 08:12:00'), (2, 'in progress', '2023-10-01 08:15:00'),
(3, 'pending', '2023-10-01 08:20:00'), (3, 'confirmed', '2023-10-01 08:22:00'),
(4, 'pending', '2023-10-01 08:30:00'),
(5, 'pending', '2023-10-01 08:40:00'), (5, 'confirmed', '2023-10-01 08:42:00'), (5, 'in progress', '2023-10-01 08:45:00'), (5, 'completed', '2023-10-01 08:55:00'),
(6, 'pending', '2023-10-01 09:00:00'), (6, 'cancelled', '2023-10-01 09:05:00'),
(7, 'pending', '2023-10-01 09:10:00'), (7, 'confirmed', '2023-10-01 09:12:00'), (7, 'in progress', '2023-10-01 09:15:00'), (7, 'completed', '2023-10-01 09:25:00'),
(8, 'pending', '2023-10-01 09:30:00'), (8, 'confirmed', '2023-10-01 09:32:00'), (8, 'in progress', '2023-10-01 09:35:00'),
(9, 'pending', '2023-10-01 09:40:00'), (9, 'confirmed', '2023-10-01 09:42:00'), (9, 'in progress', '2023-10-01 09:45:00'), (9, 'completed', '2023-10-01 09:55:00'),
(10, 'pending', '2023-10-01 10:00:00'), (10, 'confirmed', '2023-10-01 10:02:00'),
(11, 'pending', '2023-10-01 10:10:00'),
(12, 'pending', '2023-10-01 10:20:00'), (12, 'confirmed', '2023-10-01 10:22:00'), (12, 'in progress', '2023-10-01 10:25:00'), (12, 'completed', '2023-10-01 10:35:00'),
(13, 'pending', '2023-10-01 10:40:00'), (13, 'rejected', '2023-10-01 10:45:00'),
(14, 'pending', '2023-10-01 10:50:00'), (14, 'confirmed', '2023-10-01 10:52:00'), (14, 'in progress', '2023-10-01 10:55:00'), (14, 'completed', '2023-10-01 11:05:00'),
(15, 'pending', '2023-10-01 11:10:00'), (15, 'confirmed', '2023-10-01 11:12:00'), (15, 'in progress', '2023-10-01 11:15:00'),
(16, 'pending', '2023-10-01 11:20:00'), (16, 'confirmed', '2023-10-01 11:22:00'), (16, 'in progress', '2023-10-01 11:25:00'), (16, 'completed', '2023-10-01 11:35:00'),
(17, 'pending', '2023-10-01 11:40:00'), (17, 'confirmed', '2023-10-01 11:42:00'),
(18, 'pending', '2023-10-01 11:50:00'), (18, 'confirmed', '2023-10-01 11:52:00'), (18, 'in progress', '2023-10-01 11:55:00'), (18, 'completed', '2023-10-01 12:05:00'),
(19, 'pending', '2023-10-01 12:10:00'), (19, 'confirmed', '2023-10-01 12:12:00'), (19, 'in progress', '2023-10-01 12:15:00'), (19, 'completed', '2023-10-01 12:25:00'),
(20, 'pending', '2023-10-01 12:30:00'), (20, 'confirmed', '2023-10-01 12:32:00'), (20, 'in progress', '2023-10-01 12:35:00'), (20, 'completed', '2023-10-01 12:45:00'),
(21, 'pending', '2023-10-01 12:50:00'), (21, 'confirmed', '2023-10-01 12:52:00'), (21, 'in progress', '2023-10-01 12:55:00'),
(22, 'pending', '2023-10-01 13:00:00'), (22, 'confirmed', '2023-10-01 13:02:00'), (22, 'in progress', '2023-10-01 13:05:00'), (22, 'completed', '2023-10-01 13:15:00'),
(23, 'pending', '2023-10-01 13:20:00'), (23, 'confirmed', '2023-10-01 13:22:00'), (23, 'in progress', '2023-10-01 13:25:00'), (23, 'completed', '2023-10-01 13:35:00'),
(24, 'pending', '2023-10-01 13:40:00'), (24, 'confirmed', '2023-10-01 13:42:00'), (24, 'in progress', '2023-10-01 13:45:00'), (24, 'completed', '2023-10-01 13:55:00'),
(25, 'pending', '2023-10-01 14:00:00'),
(26, 'pending', '2023-10-01 14:10:00'), (26, 'confirmed', '2023-10-01 14:12:00'), (26, 'in progress', '2023-10-01 14:15:00'), (26, 'completed', '2023-10-01 14:25:00'),
(27, 'pending', '2023-10-01 14:30:00'), (27, 'confirmed', '2023-10-01 14:32:00'), (27, 'in progress', '2023-10-01 14:35:00'),
(28, 'pending', '2023-10-01 14:40:00'), (28, 'confirmed', '2023-10-01 14:42:00'), (28, 'in progress', '2023-10-01 14:45:00'), (28, 'completed', '2023-10-01 14:55:00'),
(29, 'pending', '2023-10-01 15:00:00'), (29, 'confirmed', '2023-10-01 15:02:00'),
(30, 'pending', '2023-10-01 15:10:00'), (30, 'confirmed', '2023-10-01 15:12:00'), (30, 'in progress', '2023-10-01 15:15:00'), (30, 'completed', '2023-10-01 15:25:00');

INSERT INTO price_history (menu_items_id, old_price, new_price, changed_at) VALUES
(1, 2.25, 2.50, '2023-09-15 00:00:00'),
(2, 3.50, 3.75, '2023-09-15 00:00:00'),
(3, 3.75, 4.00, '2023-09-15 00:00:00'),
(4, 2.75, 3.00, '2023-09-15 00:00:00'),
(5, 4.25, 4.50, '2023-09-15 00:00:00'),
(6, 4.00, 4.25, '2023-09-15 00:00:00'),
(7, 3.25, 3.50, '2023-09-15 00:00:00'),
(8, 3.75, 4.00, '2023-09-15 00:00:00'),
(9, 4.25, 4.50, '2023-09-15 00:00:00'),
(10, 4.00, 4.25, '2023-09-15 00:00:00'),
(11, 3.50, 3.75, '2023-09-15 00:00:00'),
(12, 4.50, 4.75, '2023-09-15 00:00:00'),
(13, 3.00, 3.25, '2023-09-15 00:00:00'),
(14, 4.75, 5.00, '2023-09-15 00:00:00'),
(15, 3.50, 3.75, '2023-09-15 00:00:00'),
(16, 3.25, 3.50, '2023-09-15 00:00:00'),
(17, 6.25, 6.50, '2023-09-15 00:00:00'),
(18, 4.50, 4.75, '2023-09-15 00:00:00'),
(19, 5.00, 5.25, '2023-09-15 00:00:00'),
(20, 3.50, 3.75, '2023-09-15 00:00:00'),
(21, 4.75, 5.00, '2023-09-15 00:00:00'),
(22, 5.25, 5.50, '2023-09-15 00:00:00'),
(23, 4.50, 4.75, '2023-09-15 00:00:00'),
(24, 4.25, 4.50, '2023-09-15 00:00:00'),
(25, 4.50, 4.75, '2023-09-15 00:00:00'),
(26, 4.25, 4.50, '2023-09-15 00:00:00'),
(27, 4.25, 4.50, '2023-09-15 00:00:00'),
(28, 4.50, 4.75, '2023-09-15 00:00:00'),
(29, 5.00, 5.25, '2023-09-15 00:00:00'),
(30, 2.50, 2.75, '2023-09-15 00:00:00');