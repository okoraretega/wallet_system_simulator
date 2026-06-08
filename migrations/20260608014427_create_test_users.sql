-- +goose Up
WITH inserted_users AS (
    INSERT INTO users (first_name, last_name, email)
    VALUES
    ('Tega', 'Okorare', 'tega@gmail.com'),
    ('Jackson', 'Henry', 'henry@gmail.com')
    RETURNING id
)

INSERT INTO wallets (wallet_number, user_id)
SELECT
    LPAD(ROW_NUMBER() OVER ()::TEXT, 10, '0'),
    id
FROM inserted_users;

-- +goose Down
DELETE FROM users
WHERE email IN ('tega@gmail.com', 'henry@gmail.com');
