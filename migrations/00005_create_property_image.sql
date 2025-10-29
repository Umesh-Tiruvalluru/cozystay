-- +goose Up
-- +goose StatementBegin
CREATE TABLE property_images (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    caption TEXT NOT NULL,
    display_order INT DEFAULT 0 NOT NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE property_images IF EXISTS;
-- +goose StatementEnd
