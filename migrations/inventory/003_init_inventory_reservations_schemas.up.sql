CREATE TABLE inventory_reservations (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    order_id    UUID NOT NULL,
    product_id  UUID NOT NULL REFERENCES products(id),
    quantity    INT NOT NULL CHECK (quantity > 0),
    status      TEXT NOT NULL CHECK (status IN ('RESERVED', 'RELEASED')),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),

    UNIQUE (order_id, product_id)
);
