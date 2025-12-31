CREATE TABLE products (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    sku        VARCHAR(100) UNIQUE NOT NULL,
    price      NUMERIC(10, 2) NOT NULL,
    stock      INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

INSERT INTO products (name, sku, price, stock, created_at, updated_at)
VALUES
    ('Mechanical Keyboard', 'KB-MECH-001', 89.99, 3, NOW(), null),
    ('Wireless Mouse', 'MS-WLS-002', 49.50, 50, NOW(), null),
    ('27-inch Monitor', 'MN-27-003', 279.00, 10, NOW(), null),
    ('USB-C Hub', 'HB-USBC-004', 39.99, 100, NOW(), null);