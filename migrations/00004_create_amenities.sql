-- +goose Up
-- +goose StatementBegin
CREATE TABLE amenities (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE property_amenities (
    id UUID PRIMARY KEY DEFAULT GEN_RANDOM_UUID(),
    property_id UUID REFERENCES properties(id) ON DELETE CASCADE,
    amenity_id UUID REFERENCES amenities(id) ON DELETE CASCADE,
    UNIQUE(property_id, amenity_id)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS property_amenities;
DROP TABLE IF EXISTS amenities;
-- +goose StatementEnd
