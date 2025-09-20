-- +goose up
ALTER TABLE users
ADD COLUMN hashed_password VARCHAR(256) NOT NULL DEFAULT 'unset';

-- +goose down
ALTER TABLE users
DROP COLUMN hashed_password;
