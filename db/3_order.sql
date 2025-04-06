CREATE TYPE order_status AS ENUM ('processing', 'accepted', 'rejected');

CREATE TABLE orders (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    customer_name VARCHAR(64) NOT NULL,
    status order_status NOT NULL DEFAULT 'processing',
    allergens VARCHAR(64) [],
    total DECIMAL(10, 2) NOT NULL CHECK (total >= 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

CREATE TABLE order_items (
    -- id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    product_id INT NOT NULL REFERENCES menu_items (id) ON DELETE CASCADE,
    quantity INT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (order_id, product_id)
);

CREATE TABLE order_status_history (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    order_id INT NOT NULL REFERENCES orders (id) ON DELETE CASCADE,
    status order_status NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--orders
CREATE INDEX idx_orders_customer_name ON orders USING GIN (
    to_tsvector('english', customer_name)
);

CREATE INDEX idx_orders_allergens ON orders USING GIN (allergens);

-- 1. Функция-триггер
CREATE OR REPLACE FUNCTION log_order_status_change()
RETURNS TRIGGER AS $$
BEGIN
  -- Проверяем, изменился ли статус
  IF NEW.status IS DISTINCT FROM OLD.status THEN
    INSERT INTO order_status_history (order_id, status)
    VALUES (NEW.id, NEW.status);
  END IF;
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 2. Сам триггер
CREATE TRIGGER trg_log_order_status_change
AFTER UPDATE ON orders
FOR EACH ROW
EXECUTE FUNCTION log_order_status_change();

INSERT INTO
    orders (
        customer_name,
        status,
        allergens,
        total,
        created_at,
        updated_at
    )
VALUES (
        'Alice Smith',
        'processing',
        ARRAY['dairy'],
        30.00,
        '2024-01-05',
        '2024-01-05'
    ),
    (
        'Bob Johnson',
        'accepted',
        ARRAY['gluten'],
        25.50,
        '2024-01-10',
        '2024-01-11'
    ),
    (
        'Charlie Brown',
        'rejected',
        NULL,
        0.00,
        '2024-01-15',
        '2024-01-16'
    ),
    (
        'David Wilson',
        'processing',
        ARRAY['nuts'],
        28.75,
        '2024-01-20',
        '2024-01-20'
    ),
    (
        'Emma Davis',
        'accepted',
        NULL,
        35.20,
        '2024-01-25',
        '2024-01-26'
    ),
    (
        'Frank Miller',
        'processing',
        ARRAY['dairy'],
        40.00,
        '2024-02-02',
        '2024-02-02'
    ),
    (
        'Grace Lee',
        'accepted',
        ARRAY['gluten'],
        22.00,
        '2024-02-08',
        '2024-02-09'
    ),
    (
        'Henry Moore',
        'rejected',
        NULL,
        0.00,
        '2024-02-14',
        '2024-02-15'
    ),
    (
        'Isla White',
        'processing',
        ARRAY['nuts'],
        45.60,
        '2024-02-20',
        '2024-02-20'
    ),
    (
        'Jack Taylor',
        'accepted',
        NULL,
        50.25,
        '2024-02-25',
        '2024-02-26'
    ),
    (
        'Liam Scott',
        'processing',
        ARRAY['dairy'],
        30.00,
        '2024-03-01',
        '2024-03-01'
    ),
    (
        'Mia Carter',
        'accepted',
        ARRAY['gluten'],
        23.90,
        '2024-03-05',
        '2024-03-06'
    ),
    (
        'Noah Wright',
        'rejected',
        NULL,
        0.00,
        '2024-03-10',
        '2024-03-11'
    ),
    (
        'Olivia Harris',
        'processing',
        ARRAY['nuts'],
        38.40,
        '2024-03-15',
        '2024-03-15'
    ),
    (
        'Paul Adams',
        'accepted',
        NULL,
        27.30,
        '2024-03-20',
        '2024-03-21'
    ),
    (
        'Quinn Baker',
        'processing',
        ARRAY['dairy'],
        31.10,
        '2024-03-25',
        '2024-03-25'
    ),
    (
        'Rachel Clark',
        'accepted',
        ARRAY['gluten'],
        26.50,
        '2024-03-30',
        '2024-03-31'
    ),
    (
        'Samuel Nelson',
        'rejected',
        NULL,
        0.00,
        '2024-04-05',
        '2024-04-06'
    ),
    (
        'Tina Young',
        'processing',
        ARRAY['nuts'],
        33.20,
        '2024-04-10',
        '2024-04-10'
    ),
    (
        'Umar King',
        'accepted',
        NULL,
        29.75,
        '2024-04-15',
        '2024-04-16'
    ),
    (
        'Victor Green',
        'processing',
        ARRAY['dairy'],
        42.00,
        '2024-04-20',
        '2024-04-20'
    ),
    (
        'Wendy Hall',
        'accepted',
        ARRAY['gluten'],
        25.80,
        '2024-04-25',
        '2024-04-26'
    ),
    (
        'Xavier Brown',
        'rejected',
        NULL,
        0.00,
        '2024-04-30',
        '2024-05-01'
    ),
    (
        'Yasmine Perez',
        'processing',
        ARRAY['nuts'],
        38.50,
        '2024-05-05',
        '2024-05-05'
    ),
    (
        'Zachary Reed',
        'accepted',
        NULL,
        30.10,
        '2024-05-10',
        '2024-05-11'
    ),
    (
        'Amy Foster',
        'processing',
        ARRAY['dairy'],
        32.75,
        '2024-05-15',
        '2024-05-15'
    ),
    (
        'Brian Stewart',
        'accepted',
        ARRAY['gluten'],
        24.00,
        '2024-05-20',
        '2024-05-21'
    ),
    (
        'Catherine Lewis',
        'rejected',
        NULL,
        0.00,
        '2024-05-25',
        '2024-05-26'
    ),
    (
        'Daniel Martinez',
        'processing',
        ARRAY['nuts'],
        35.00,
        '2024-05-30',
        '2024-05-30'
    ),
    (
        'Urystem Qabdolla',
        'processing',
        ARRAY['nuts'],
        37.40,
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