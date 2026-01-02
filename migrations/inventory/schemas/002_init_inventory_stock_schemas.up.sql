CREATE TABLE inventory_stock (
    product_id    UUID PRIMARY KEY REFERENCES products(id),
    total_qty     INT NOT NULL CHECK (total_qty >= 0),
    reserved_qty  INT NOT NULL DEFAULT 0 CHECK (reserved_qty >= 0),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
