CREATE TABLE
    IF NOT EXISTS bank_accounts(
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        userId UUID NOT NULL,
        bankName VARCHAR(255) NOT NULL,
        accountNumber VARCHAR(255) NOT NULL,
        accountHolder VARCHAR(255) NOT NULL,
        createdAt TIMESTAMP DEFAULT now(),
        updatedAt TIMESTAMP DEFAULT now(),

        CONSTRAINT fk_user_id_users FOREIGN KEY (userId) REFERENCES users(id) ON DELETE CASCADE  
    );

CREATE TRIGGER update_bank_accounts_updated_at
  BEFORE UPDATE
  ON bank_accounts
  FOR EACH ROW
  EXECUTE PROCEDURE trigger_set_updated();