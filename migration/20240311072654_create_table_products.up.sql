CREATE TABLE
    IF NOT EXISTS products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    imageUrl VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    condition product_condition NOT NULL,
    tags VARCHAR[] NOT NULL,
    isPurchaseable BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP DEFAULT now(),
    updated_at TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_user_id_users FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE  
);