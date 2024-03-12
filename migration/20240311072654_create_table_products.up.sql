CREATE TYPE product_condition AS ENUM ('new', 'second');

CREATE TABLE
    IF NOT EXISTS products(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    userId UUID NOT NULL,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    imageUrl VARCHAR(255) NOT NULL,
    stock INT NOT NULL,
    condition product_condition NOT NULL,
    tags VARCHAR[] NOT NULL,
    isPurchaseable BOOLEAN NOT NULL DEFAULT TRUE,
    createdAt TIMESTAMP DEFAULT now(),
    updatedAt TIMESTAMP DEFAULT now(),

    CONSTRAINT fk_user_id_users FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE  
);

CREATE TRIGGER update_products_updated_at
  BEFORE UPDATE
  ON products
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_updated();