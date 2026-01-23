CREATE TABLE products (
    product_id SERIAL PRIMARY KEY,
    product_name VARCHAR(255),
    product_price NUMERIC,
    product_description TEXT,
    product_quantity INT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);