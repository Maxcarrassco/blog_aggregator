-- +goose Up

BEGIN;
ALTER TABLE users
ADD COLUMN api_Key VARCHAR(64);

UPDATE users
SET api_Key = SHA2(UUID(), 256);

ALTER TABLE users
MODIFY api_Key VARCHAR(64) NOT NULL;
COMMIT;

-- +goose Down
ALTER TABLE users DROP COLUMN api_Key;
