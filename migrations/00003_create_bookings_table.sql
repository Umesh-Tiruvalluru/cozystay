-- +goose Up
-- +goose StatementBegin
CREATE TABLE bookings (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id     UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    property_id UUID NOT NULL REFERENCES properties(id) ON DELETE CASCADE,
    start_date  DATE NOT NULL,
    end_date    DATE NOT NULL,
    total_price NUMERIC(10, 2) NOT NULL CHECK (total_price >= 0),
    status TEXT CHECK (status IN ('booked', 'cancelled')) DEFAULT 'booked',
    created_at  TIMESTAMP DEFAULT NOW(),
    CONSTRAINT  chk_date_range CHECK (end_date > start_date)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS bookings;
-- +goose StatementEnd
