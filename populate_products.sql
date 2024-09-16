-- Create the products table (if it doesn't exist)
CREATE TABLE IF NOT EXISTS products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    price NUMERIC(10, 2) NOT NULL,
    stock INTEGER NOT NULL
);

-- Insert sample records
INSERT INTO products (name, description, price, stock) VALUES
('Product 1', 'Description for Product 1', 19.99, 100),
('Product 2', 'Description for Product 2', 29.99, 200),
('Product 3', 'Description for Product 3', 39.99, 300),
('Product 4', 'Description for Product 4', 49.99, 400),
('Product 5', 'Description for Product 5', 59.99, 500);
