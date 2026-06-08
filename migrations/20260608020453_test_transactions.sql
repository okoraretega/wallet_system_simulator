-- +goose Up
BEGIN;

UPDATE wallets SET wallet_balance = wallet_balance - 100000 WHERE user_id = 'e209e36c-87d1-4601-89bb-f4e161bddc93';
UPDATE wallets SET wallet_balance = wallet_balance + 100000 WHERE user_id = '6b853eca-146f-41a6-b429-8f2af274ae12';

COMMIT;

-- +goose Down
UPDATE wallets SET wallet_balance = wallet_balance + 100000 WHERE user_id = 'e209e36c-87d1-4601-89bb-f4e161bddc93';
UPDATE wallets SET wallet_balance = wallet_balance - 100000 WHERE user_id = '6b853eca-146f-41a6-b429-8f2af274ae12';
