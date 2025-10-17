-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS users (
    id            UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    email         TEXT NOT NULL UNIQUE,
    first_name    TEXT NOT NULL,
    last_name     TEXT NOT NULL,
    created_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMP NOT NULL DEFAULT NOW(),
    role TEXT CHECK (role IN ('guest', 'admin')) DEFAULT 'guest',
    password_hash TEXT NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS users;
-- +goose StatementEnd      

