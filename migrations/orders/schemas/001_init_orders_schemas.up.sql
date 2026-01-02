CREATE TABLE orders (
    id VARCHAR(36) PRIMARY KEY,
    order_number VARCHAR(20) UNIQUE NOT NULL,
    user_id VARCHAR(36) NOT NULL,
    
    status VARCHAR(20) NOT NULL, -- PENDING, CONFIRMED, PAID, CANCELLED
    total_amount DECIMAL(10, 2) NOT NULL,
    
    shipping_address TEXT NOT NULL,
    shipping_phone VARCHAR(20) NOT NULL,

    payment_transaction_id VARCHAR(100),
    shipping_tracking_id VARCHAR(100),

    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP
);
