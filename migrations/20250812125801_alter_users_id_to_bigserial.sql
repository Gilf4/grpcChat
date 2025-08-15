-- +goose Up
ALTER TABLE users
    ALTER COLUMN id TYPE bigint;

-- +goose Down
ALTER TABLE users
    ALTER COLUMN id TYPE integer;
