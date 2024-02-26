CREATE TABLE wallet (
    id  SERIAL PRIMARY KEY,
    account UUID NOT NULL,
    amount INTEGER DEFAULT 0 NOT NULL CONSTRAINT positive_amount CHECK (amount >= 0),
    currency VARCHAR(3) NOT NULL,
    updated_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL
);

CREATE TABLE transaction (
    wallet_id INTEGER NOT NULL,
    transaction_id UUID NOT NULL,
    created_at   TIMESTAMPTZ DEFAULT NOW() NOT NULL,
    PRIMARY KEY (wallet_id, transaction_id),
    CONSTRAINT fk_wallet
      FOREIGN KEY(wallet_id) 
        REFERENCES wallet(id)
);

CREATE INDEX idx_transaction_wallet_transaction_id ON transaction (wallet_id, transaction_id);

CREATE  FUNCTION update_updated_at_trigger()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_wallet_updated_at
    BEFORE UPDATE ON wallet
    FOR EACH ROW
EXECUTE PROCEDURE update_updated_at_trigger();