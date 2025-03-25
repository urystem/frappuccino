CREATE OR REPLACE FUNCTION check_inventory_before_order()
RETURNS TRIGGER AS $$
DECLARE
    required_quantity INT;
    available_quantity INT;
BEGIN
    -- Loop through each ingredient required for the ordered menu item
    FOR required_quantity, available_quantity IN
        SELECT mi.quantity, ii.quantity
        FROM menu_item_ingredients mi
        JOIN inventory_items ii ON mi.inventory_item_id = ii.inventory_item_id
        WHERE mi.menu_item_id = NEW.menu_item_id
    LOOP
        -- If there's not enough stock, cancel the order
        IF available_quantity < required_quantity * NEW.quantity THEN
            UPDATE orders SET status = 'cancelled' WHERE order_id = NEW.order_id;
            RAISE EXCEPTION 'Order % has been cancelled due to insufficient inventory', NEW.order_id;
        END IF;
    END LOOP;

    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER enforce_inventory_check
BEFORE INSERT ON order_items
FOR EACH ROW
EXECUTE FUNCTION check_inventory_before_order();


CREATE OR REPLACE FUNCTION check_inventory_exists()
RETURNS TRIGGER AS $$
BEGIN
    IF EXISTS (
        SELECT 1 FROM inventory_items WHERE inventory_name = NEW.inventory_name
    ) THEN
        RAISE EXCEPTION 'Inventory item "%" already exists', NEW.inventory_name;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER before_inventory_insert
BEFORE INSERT ON inventory_items
FOR EACH ROW
EXECUTE FUNCTION check_inventory_exists();
