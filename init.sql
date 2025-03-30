--CREATE TABLE
--inventory
CREATE TYPE uints AS ENUM ('g', 'ml', 'pcs');

CREATE TABLE inventory (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    quantity FLOAT NOT NULL CHECK (quantity >= 0),
    reorder_level FLOAT NOT NULL CHECK (reorder_level > 0),
    unit uints NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);

CREATE TYPE reason_of_inventory_transaction AS ENUM ('restock', 'usage', 'cancelled', 'annul');

CREATE TABLE inventory_transactions (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    inventory_id INT NOT NULL REFERENCES inventory (id) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason reason_of_inventory_transaction NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--menu
CREATE TABLE menu_items (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    tags TEXT [],
    allergen TEXT [],
    price DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0) --inventories TEXT[] NOT NULL CHECK (array_length(allergens, 1) > 0) --cardinality(allergens)>0
);

CREATE TABLE menu_item_ingredients (
    product_id INT NOT NULL REFERENCES menu_items (id) ON DELETE CASCADE,
    inventory_id INT NOT NULL REFERENCES inventory (id),
    quantity FLOAT NOT NULL CHECK (quantity > 0)
);

--need trigger
CREATE TABLE price_history (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id INT NOT NULL REFERENCES menu_items (id) ON DELETE CASCADE,
    old_price DECIMAL(10, 2) NOT NULL CHECK (old_price >= 0),
    new_price DECIMAL(10, 2) NOT NULL CHECK (new_price >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--order
CREATE TYPE order_status AS ENUM ('processing', 'accepted', 'rejected');

CREATE TABLE orders (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_name VARCHAR(64) NOT NULL,
    status order_status NOT NULL,
    allergen TEXT [],
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

CREATE TABLE order_items (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES menu_items (id),
    quantity INT NOT NULL CHECK (quantity > 0)
);

CREATE TABLE order_status_history (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--INDEXING
--inventory
CREATE INDEX idx_inventory_name ON inventory USING GIN (to_tsvector('english', name));

CREATE INDEX idx_inventory_description ON inventory USING GIN (
    to_tsvector('english', description)
);

--menu_items
CREATE INDEX idx_menu_items_name ON menu_items USING GIN (to_tsvector('english', name));

CREATE INDEX idx_menu_items_description ON menu_items USING GIN (
    to_tsvector('english', description)
);

CREATE INDEX idx_menu_items_tags ON menu_items USING GIN (tags);

CREATE INDEX idx_menu_items_allergen ON menu_items USING GIN (allergen);

--orders
CREATE INDEX idx_orders_customer_name ON orders USING GIN (
    to_tsvector('english', customer_name)
);

CREATE INDEX idx_orders_allergen ON orders USING GIN (allergen);

--INSERT TO THE TABLE
INSERT INTO
    inventory (
        name,
        description,
        quantity,
        reorder_level,
        unit,
        price
    )
VALUES (
        'Espresso Beans',
        'High-quality coffee beans',
        5000.0,
        200.0,
        'g',
        1500.00
    ),
    (
        'Milk',
        'Fresh dairy milk',
        100.0,
        10.0,
        'ml',
        50.00
    ),
    (
        'Sugar',
        'Refined white sugar',
        2000.0,
        500.0,
        'g',
        100.00
    ),
    (
        'Flour',
        'All-purpose wheat flour',
        10000.0,
        500.0,
        'g',
        300.00
    ),
    (
        'Butter',
        'Unsalted butter',
        500.0,
        50.0,
        'g',
        800.00
    ),
    (
        'Eggs',
        'Organic chicken eggs',
        30.0,
        10.0,
        'pcs',
        200.00
    ),
    (
        'Vanilla Extract',
        'Pure vanilla flavoring',
        100.0,
        20.0,
        'ml',
        900.00
    ),
    (
        'Chocolate Chips',
        'Dark chocolate pieces',
        2000.0,
        200.0,
        'g',
        1200.00
    ),
    (
        'Honey',
        'Organic wildflower honey',
        500.0,
        50.0,
        'ml',
        700.00
    ),
    (
        'Cinnamon Powder',
        'Ground cinnamon spice',
        300.0,
        30.0,
        'g',
        400.00
    ),
    (
        'Cocoa Powder',
        'Unsweetened cocoa',
        1000.0,
        100.0,
        'g',
        1100.00
    ),
    (
        'Baking Powder',
        'Leavening agent for baking',
        500.0,
        50.0,
        'g',
        250.00
    ),
    (
        'Salt',
        'Fine sea salt',
        2000.0,
        500.0,
        'g',
        100.00
    ),
    (
        'Lemon Juice',
        'Freshly squeezed lemon juice',
        500.0,
        50.0,
        'ml',
        600.00
    ),
    (
        'Olive Oil',
        'Extra virgin olive oil',
        1000.0,
        100.0,
        'ml',
        1200.00
    ),
    (
        'Yeast',
        'Active dry yeast',
        250.0,
        50.0,
        'g',
        350.00
    ),
    (
        'Maple Syrup',
        'Natural maple syrup',
        500.0,
        50.0,
        'ml',
        1500.00
    ),
    (
        'Whipping Cream',
        'Heavy cream for desserts',
        1000.0,
        100.0,
        'ml',
        900.00
    ),
    (
        'Oats',
        'Whole grain rolled oats',
        2000.0,
        500.0,
        'g',
        500.00
    ),
    (
        'Almonds',
        'Raw almonds',
        1000.0,
        200.0,
        'g',
        1400.00
    );

INSERT INTO
    inventory_transactions (
        inventory_id,
        quantity_change,
        reason
    )
VALUES (1, 5000.0, 'restock'),
    (2, 100.0, 'restock'),
    (3, 2000.0, 'restock'),
    (4, 10000.0, 'restock'),
    (5, 500.0, 'restock'),
    (6, 30.0, 'restock'),
    (7, 100.0, 'restock'),
    (8, 2000.0, 'restock'),
    (9, 500.0, 'restock'),
    (10, 300.0, 'restock'),
    (11, 1000.0, 'restock'),
    (12, 500.0, 'restock'),
    (13, 2000.0, 'restock'),
    (14, 500.0, 'restock'),
    (15, 1000.0, 'restock'),
    (16, 250.0, 'restock'),
    (17, 500.0, 'restock'),
    (18, 1000.0, 'restock'),
    (19, 2000.0, 'restock'),
    (20, 1000.0, 'restock');

--MENU
INSERT INTO
    menu_items (
        name,
        description,
        tags,
        allergen,
        price
    )
VALUES (
        'Espresso',
        'Rich and bold espresso shot',
        ARRAY['coffee'],
        ARRAY['none'],
        3.00
    ),
    (
        'Cappuccino',
        'Espresso with steamed milk and foam',
        ARRAY['coffee', 'milk'],
        ARRAY['dairy'],
        4.50
    ),
    (
        'Chocolate Chip Cookies',
        'Homemade cookies with dark chocolate chips',
        ARRAY['dessert', 'sweet'],
        ARRAY['gluten', 'dairy'],
        5.00
    ),
    (
        'Honey Oatmeal',
        'Warm oatmeal sweetened with organic honey',
        ARRAY['breakfast', 'healthy'],
        ARRAY['none'],
        6.00
    ),
    (
        'Lemon Tart',
        'Tangy lemon-flavored tart with a buttery crust',
        ARRAY['dessert'],
        ARRAY['gluten', 'dairy'],
        7.00
    ),
    (
        'Vanilla Ice Cream',
        'Classic vanilla-flavored ice cream',
        ARRAY['dessert'],
        ARRAY['dairy'],
        4.50
    ),
    (
        'Cinnamon Rolls',
        'Soft rolls with cinnamon sugar and icing',
        ARRAY['dessert', 'sweet'],
        ARRAY['gluten', 'dairy'],
        6.50
    ),
    (
        'Chocolate Brownie',
        'Rich and fudgy chocolate brownie',
        ARRAY['dessert', 'chocolate'],
        ARRAY['gluten', 'dairy'],
        5.50
    ),
    (
        'Maple Pancakes',
        'Fluffy pancakes topped with maple syrup',
        ARRAY['breakfast', 'sweet'],
        ARRAY['gluten', 'dairy'],
        7.50
    ),
    (
        'Almond Croissant',
        'Flaky croissant with almond filling',
        ARRAY['pastry'],
        ARRAY['gluten', 'dairy', 'nuts'],
        6.00
    );

INSERT INTO
    menu_item_ingredients (
        product_id,
        inventory_id,
        quantity
    )
VALUES (1, 1, 30), -- Espresso -> Espresso Beans
    (2, 1, 30), -- Cappuccino -> Espresso Beans
    (2, 2, 100), -- Cappuccino -> Milk
    (3, 4, 200), -- Chocolate Chip Cookies -> Flour
    (3, 5, 100), -- Chocolate Chip Cookies -> Butter
    (3, 8, 50), -- Chocolate Chip Cookies -> Chocolate Chips
    (3, 6, 2), -- Chocolate Chip Cookies -> Eggs
    (3, 12, 10), -- Chocolate Chip Cookies -> Baking Powder
    (4, 19, 100), -- Honey Oatmeal -> Oats
    (4, 9, 20), -- Honey Oatmeal -> Honey
    (5, 4, 150), -- Lemon Tart -> Flour
    (5, 5, 100), -- Lemon Tart -> Butter
    (5, 6, 2), -- Lemon Tart -> Eggs
    (5, 14, 30), -- Lemon Tart -> Lemon Juice
    (6, 7, 10), -- Vanilla Ice Cream -> Vanilla Extract
    (6, 2, 150), -- Vanilla Ice Cream -> Milk
    (6, 18, 100), -- Vanilla Ice Cream -> Whipping Cream
    (7, 4, 300), -- Cinnamon Rolls -> Flour
    (7, 5, 50), -- Cinnamon Rolls -> Butter
    (7, 6, 1), -- Cinnamon Rolls -> Eggs
    (7, 10, 5), -- Cinnamon Rolls -> Cinnamon Powder
    (7, 16, 10), -- Cinnamon Rolls -> Yeast
    (8, 4, 200), -- Chocolate Brownie -> Flour
    (8, 5, 100), -- Chocolate Brownie -> Butter
    (8, 6, 2), -- Chocolate Brownie -> Eggs
    (8, 11, 50), -- Chocolate Brownie -> Cocoa Powder
    (8, 8, 50), -- Chocolate Brownie -> Chocolate Chips
    (9, 4, 200), -- Maple Pancakes -> Flour
    (9, 5, 50), -- Maple Pancakes -> Butter
    (9, 6, 2), -- Maple Pancakes -> Eggs
    (9, 17, 50), -- Maple Pancakes -> Maple Syrup
    (10, 4, 150), -- Almond Croissant -> Flour
    (10, 5, 50), -- Almond Croissant -> Butter
    (10, 6, 1), -- Almond Croissant -> Eggs
    (10, 20, 30);

INSERT INTO
    price_history (
        product_id,
        old_price,
        new_price,
        updated_at
    )
VALUES (1, 2.80, 3.00, '2024-01-10'),
    (2, 4.30, 4.50, '2024-01-12'),
    (3, 4.80, 5.00, '2024-02-01'),
    (4, 5.80, 6.00, '2024-03-18'),
    (5, 6.80, 7.00, '2024-04-25'),
    (6, 4.30, 4.50, '2024-05-05'),
    (7, 6.20, 6.50, '2024-03-30'),
    (8, 5.30, 5.50, '2024-03-14'),
    (9, 7.20, 7.50, '2024-04-22'),
    (10, 5.80, 6.00, '2024-11-28');

--ORDER

INSERT INTO
    orders (
        customer_name,
        status,
        allergen,
        created_at,
        updated_at
    )
VALUES (
        'Alice Smith',
        'processing',
        ARRAY['dairy'],
        '2024-01-05',
        '2024-01-05'
    ),
    (
        'Bob Johnson',
        'accepted',
        ARRAY['gluten'],
        '2024-01-10',
        '2024-01-11'
    ),
    (
        'Charlie Brown',
        'rejected',
        ARRAY['none'],
        '2024-01-15',
        '2024-01-16'
    ),
    (
        'David Wilson',
        'processing',
        ARRAY['nuts'],
        '2024-01-20',
        '2024-01-20'
    ),
    (
        'Emma Davis',
        'accepted',
        ARRAY['none'],
        '2024-01-25',
        '2024-01-26'
    ),
    (
        'Frank Miller',
        'processing',
        ARRAY['dairy'],
        '2024-02-02',
        '2024-02-02'
    ),
    (
        'Grace Lee',
        'accepted',
        ARRAY['gluten'],
        '2024-02-08',
        '2024-02-09'
    ),
    (
        'Henry Moore',
        'rejected',
        ARRAY['none'],
        '2024-02-14',
        '2024-02-15'
    ),
    (
        'Isla White',
        'processing',
        ARRAY['nuts'],
        '2024-02-20',
        '2024-02-20'
    ),
    (
        'Jack Taylor',
        'accepted',
        ARRAY['none'],
        '2024-02-25',
        '2024-02-26'
    ),
    (
        'Liam Scott',
        'processing',
        ARRAY['dairy'],
        '2024-03-01',
        '2024-03-01'
    ),
    (
        'Mia Carter',
        'accepted',
        ARRAY['gluten'],
        '2024-03-05',
        '2024-03-06'
    ),
    (
        'Noah Wright',
        'rejected',
        ARRAY['none'],
        '2024-03-10',
        '2024-03-11'
    ),
    (
        'Olivia Harris',
        'processing',
        ARRAY['nuts'],
        '2024-03-15',
        '2024-03-15'
    ),
    (
        'Paul Adams',
        'accepted',
        ARRAY['none'],
        '2024-03-20',
        '2024-03-21'
    ),
    (
        'Quinn Baker',
        'processing',
        ARRAY['dairy'],
        '2024-03-25',
        '2024-03-25'
    ),
    (
        'Rachel Clark',
        'accepted',
        ARRAY['gluten'],
        '2024-03-30',
        '2024-03-31'
    ),
    (
        'Samuel Nelson',
        'rejected',
        ARRAY['none'],
        '2024-04-05',
        '2024-04-06'
    ),
    (
        'Tina Young',
        'processing',
        ARRAY['nuts'],
        '2024-04-10',
        '2024-04-10'
    ),
    (
        'Umar King',
        'accepted',
        ARRAY['none'],
        '2024-04-15',
        '2024-04-16'
    ),
    (
        'Victor Green',
        'processing',
        ARRAY['dairy'],
        '2024-04-20',
        '2024-04-20'
    ),
    (
        'Wendy Hall',
        'accepted',
        ARRAY['gluten'],
        '2024-04-25',
        '2024-04-26'
    ),
    (
        'Xavier Brown',
        'rejected',
        ARRAY['none'],
        '2024-04-30',
        '2024-05-01'
    ),
    (
        'Yasmine Perez',
        'processing',
        ARRAY['nuts'],
        '2024-05-05',
        '2024-05-05'
    ),
    (
        'Zachary Reed',
        'accepted',
        ARRAY['none'],
        '2024-05-10',
        '2024-05-11'
    ),
    (
        'Amy Foster',
        'processing',
        ARRAY['dairy'],
        '2024-05-15',
        '2024-05-15'
    ),
    (
        'Brian Stewart',
        'accepted',
        ARRAY['gluten'],
        '2024-05-20',
        '2024-05-21'
    ),
    (
        'Catherine Lewis',
        'rejected',
        ARRAY['none'],
        '2024-05-25',
        '2024-05-26'
    ),
    (
        'Daniel Martinez',
        'processing',
        ARRAY['nuts'],
        '2024-05-30',
        '2024-05-30'
    ),
    (
        'Urystem Qabdolla',
        'processing',
        ARRAY['nuts'],
        '2024-05-30',
        '2024-05-31'
    );

-- Вставка данных в таблицу order_items
INSERT INTO
    order_items (
        order_id,
        product_id,
        quantity
    )
VALUES (1, 2, 1),
    (2, 3, 3),
    (3, 8, 2),
    (4, 10, 2),
    (5, 4, 3),
    (6, 5, 1),
    (7, 8, 2),
    (8, 3, 1),
    (9, 7, 3),
    (10, 10, 1),
    (11, 2, 1),
    (12, 3, 3),
    (13, 8, 2),
    (14, 10, 2),
    (15, 4, 3),
    (16, 5, 1),
    (17, 8, 2),
    (18, 3, 1),
    (19, 7, 3),
    (20, 10, 1),
    (21, 2, 1),
    (22, 3, 3),
    (23, 8, 2),
    (24, 10, 2),
    (25, 4, 3),
    (26, 5, 1),
    (27, 8, 2),
    (28, 3, 1),
    (29, 7, 3),
    (30, 10, 1);

-- Вставка данных в таблицу order_status_history
INSERT INTO
    order_status_history (order_id, status, updated_at)
VALUES (1, 'processing', '2024-01-05'),
    (2, 'accepted', '2024-01-11'),
    (3, 'rejected', '2024-01-16'),
    (4, 'processing', '2024-01-20'),
    (5, 'accepted', '2024-01-26'),
    (6, 'processing', '2024-02-02'),
    (7, 'accepted', '2024-02-09'),
    (8, 'rejected', '2024-02-15'),
    (9, 'processing', '2024-02-20'),
    (10, 'accepted', '2024-02-26'),
    (
        11,
        'processing',
        '2024-03-01'
    ),
    (12, 'accepted', '2024-03-06'),
    (13, 'rejected', '2024-03-11'),
    (
        14,
        'processing',
        '2024-03-15'
    ),
    (15, 'accepted', '2024-03-21');