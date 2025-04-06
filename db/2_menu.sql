CREATE TABLE menu_items (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(64) NOT NULL UNIQUE,
    description TEXT NOT NULL,
    tags TEXT [] NOT NULL, --DEFAULT '{}'::text [], --::text[] деген '{}' ді массив қалады
    -- tags VARCHAR(128)[],
    allergens VARCHAR(64) [],
    price DECIMAL(10, 2) NOT NULL CHECK (price >= 0) --inventories TEXT[] NOT NULL CHECK (array_length(allergens, 1) > 0) --cardinality(allergens)>0
);

CREATE TABLE menu_item_ingredients (
    product_id INT NOT NULL REFERENCES menu_items (id) ON DELETE CASCADE,
    inventory_id INT NOT NULL REFERENCES inventory (id),
    -- FOREIGN KEY (inventory_id) REFERENCES inventory (id),
    quantity FLOAT NOT NULL CHECK (quantity > 0),
    PRIMARY KEY (product_id, inventory_id)
);

CREATE TABLE price_history (
    id INT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    product_id INT NOT NULL REFERENCES menu_items (id) ON DELETE CASCADE,
    old_price DECIMAL(10, 2) NOT NULL CHECK (old_price >= 0),
    new_price DECIMAL(10, 2) NOT NULL CHECK (new_price >= 0),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP --NOW()
);

--INDEXING
CREATE INDEX idx_menu_items_name ON menu_items USING GIN (to_tsvector('english', name));

CREATE INDEX idx_menu_items_description ON menu_items USING GIN (
    to_tsvector('english', description)
);



CREATE INDEX idx_menu_items_tags ON menu_items USING GIN (tags);

CREATE INDEX idx_menu_items_allergens ON menu_items USING GIN (allergens);

CREATE /*OR REPLACE*/ FUNCTION record_price_change()
RETURNS TRIGGER AS $$
BEGIN
    -- Проверяем, если цена изменилась
    IF NEW.price <> OLD.price THEN
        -- Вставляем запись в таблицу price_history
        INSERT INTO price_history (product_id, old_price, new_price)
        VALUES (NEW.id, OLD.price, NEW.price);
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER price_update_trigger
AFTER UPDATE OF price
ON menu_items
FOR EACH ROW
EXECUTE FUNCTION record_price_change();


--MENU
INSERT INTO
    menu_items (
        name,
        description,
        tags,
        allergens,
        price
    )
VALUES (
        'Espresso',
        'Rich and bold espresso shot',
        ARRAY['coffee'],
        NULL,
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
        NULL,
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

-- INSERT INTO
--     price_history (
--         product_id,
--         old_price,
--         new_price,
--         updated_at
--     )
-- VALUES (1, 2.80, 3.00, '2024-01-10'),
--     (2, 4.30, 4.50, '2024-01-12'),
--     (3, 4.80, 5.00, '2024-02-01'),
--     (4, 5.80, 6.00, '2024-03-18'),
--     (5, 6.80, 7.00, '2024-04-25'),
--     (6, 4.30, 4.50, '2024-05-05'),
--     (7, 6.20, 6.50, '2024-03-30'),
--     (8, 5.30, 5.50, '2024-03-14'),
--     (9, 7.20, 7.50, '2024-04-22'),
--     (10, 5.80, 6.00, '2024-11-28');

