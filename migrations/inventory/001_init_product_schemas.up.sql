CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE products (
    id         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    price      NUMERIC(10, 2) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL
);

INSERT INTO products (name, price, created_at, updated_at)
VALUES
    ('Mechanical Keyboard', 89.99, NOW(), NOW()),
    ('Wireless Mouse', 49.50, NOW(), NOW()),
    ('27-inch Monitor', 279.00, NOW(), NOW()),
    ('USB-C Hub', 39.99, NOW(), NOW());
