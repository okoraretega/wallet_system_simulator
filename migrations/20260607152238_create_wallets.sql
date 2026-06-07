-- +goose Up
CREATE TABLE wallets(
    id UUID PRIMARY KEY,
    wallet_number VARCHAR(10),
    user_id UUID UNIQUE REFERENCES users(id),
    created_at TIMESTAMP DEFAULT NOW()
);

-- +goose Down
DROP TABLE wallets;
