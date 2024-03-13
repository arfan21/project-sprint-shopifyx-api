CREATE TABLE
    IF NOT EXISTS payments(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        userId UUID NOT NULL,
        bankAccountId UUID NOT NULL,
        productId UUID NOT NULL,
        paymentProofImageUrl VARCHAR(255) NOT NULL,
        quantity INT NOT NULL,
        totalPrice DECIMAL(10, 2) NOT NULL,
        createdAt TIMESTAMP DEFAULT now(),
        updatedAt TIMESTAMP DEFAULT now(),

        CONSTRAINT fk_user_id_users FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE,
        CONSTRAINT fk_bank_account_id_bank_accounts FOREIGN KEY (bankAccountId) REFERENCES bank_accounts(id) ON DELETE CASCADE,
        CONSTRAINT fk_product_id_products FOREIGN KEY (productId) REFERENCES products(id) ON DELETE CASCADE
    );

CREATE TRIGGER update_payments_updated_at
    BEFORE UPDATE
    ON payments
    FOR EACH ROW
    EXECUTE PROCEDURE trigger_set_updated();