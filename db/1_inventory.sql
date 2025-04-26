CREATE TYPE uints AS ENUM ('g', 'ml', 'pcs');

CREATE TABLE inventory (
    id SERIAL PRIMARY KEY,
    name VARCHAR(48) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    quantity FLOAT NOT NULL CHECK (quantity >= 0),
    reorder_level FLOAT NOT NULL CHECK (reorder_level > 0),
    unit uints NOT NULL,
    price DECIMAL(10, 2) NOT NULL DEFAULT 0 CHECK (price >= 0)
);

CREATE TYPE reason_of_inventory_transaction AS ENUM ('restock', 'usage', 'cancelled', 'annul');

CREATE TABLE inventory_transactions (
    id SERIAL PRIMARY KEY,
    inventory_id INT NOT NULL REFERENCES inventory (id) ON DELETE CASCADE,
    quantity_change FLOAT NOT NULL,
    reason reason_of_inventory_transaction NOT NULL ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--INDEXING
CREATE INDEX idx_inventory_name ON inventory USING GIN (to_tsvector ('english', name));

CREATE INDEX idx_inventory_description ON inventory USING GIN (
    to_tsvector ('english', description)
);

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
        'Double Chocolate Cake',
        'Rich chocolate layer cake',
        5000.0,
        200.0,
        'g',
        1500.00
    ),
    (
        'Milk',
        'Fresh dairy milk Chocolate',
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
        'All-purpose wheat flour cake',
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
