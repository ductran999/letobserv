INSERT INTO inventory_stock (product_id, total_qty)
SELECT id,
       CASE name
           WHEN 'Mechanical Keyboard' THEN 10
           WHEN 'Wireless Mouse'      THEN 25
           WHEN '27-inch Monitor'     THEN 5
           WHEN 'USB-C Hub'            THEN 15
       END
FROM products
ON CONFLICT (product_id)
DO UPDATE SET total_qty = EXCLUDED.total_qty,
              updated_at = NOW();
