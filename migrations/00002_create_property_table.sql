-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS properties (
    id              UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    title           TEXT NOT NULL,
    created_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP NOT NULL DEFAULT NOW(),
    location        TEXT NOT NULL,
    max_guests      INTEGER NOT NULL,
    price_per_night NUMERIC(10, 2) NOT NULL CHECK (price_per_night >= 0),
    description     TEXT NOT NULL,
    user_id         UUID REFERENCES users(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property;
-- +goose StatementEnd
