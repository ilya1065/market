
CREATE TABLE IF NOT EXISTS product(
    id SERIAL PRIMARY KEY,
    price NUMERIC(10,2),
    name VARCHAR(255),
    description TEXT,
    stock INT NOT NULL DEFAULT 0,
    image_url UUID,
    category_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now()
    );