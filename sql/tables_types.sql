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
    orders_id INT REFERENCES orders(orders_id) ,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ,
    quantity INT NOT NULL
);

-- Menu Item Ingredients Table
CREATE TABLE menu_item_ingredients (
    menu_item_ingredients_id SERIAL PRIMARY KEY,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ,
    inventory_items_id INT REFERENCES inventory_items(inventory_items_id) ,
    quantity INT NOT NULL
);

-- Inventory Transactions Table
CREATE TABLE inventory_transactions (
    transaction_id SERIAL PRIMARY KEY,
    inventory_items_id INT REFERENCES inventory_items(inventory_items_id) ,
    quantity_change INT NOT NULL,
    transaction_date TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Order Status History Table
CREATE TABLE order_status_history (
    status_history_id SERIAL PRIMARY KEY,
    orders_id INT REFERENCES orders(orders_id) ,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL
);

-- Price History Table
CREATE TABLE price_history (
    price_history_id SERIAL PRIMARY KEY,
    menu_items_id INT REFERENCES menu_items(menu_items_id) ,
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
CREATE INDEX idx_order_items_menu_items_id ON order_items(menu_items_id);
CREATE INDEX idx_order_items_orders_id ON order_items(orders_id);
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