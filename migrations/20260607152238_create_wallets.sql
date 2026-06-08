-- +goose Up
CREATE TABLE wallets(
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    wallet_number VARCHAR(10) NOT NULL,
    wallet_balance NUMERIC(15, 4) NOT NULL DEFAULT 0.0000,
    user_id UUID UNIQUE NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE wallets;
