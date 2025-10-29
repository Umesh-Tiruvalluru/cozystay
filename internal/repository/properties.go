package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)

// TODO: Handle Null value when no images are provided
func (repo *Repository) GetAllProperties() ([]models.GetProperty, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.location,
			p.price_per_night,
			p.max_guests,
			p.created_at,
			(
				SELECT pi.image_url
				FROM property_images pi
				WHERE pi.property_id = p.id
				ORDER BY pi.display_order ASC
				LIMIT 1
			) AS thumbnail_url
		FROM properties p
		ORDER BY p.created_at DESC;
	`

	var properties []models.GetProperty
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return properties, err
	}

	for rows.Next() {
		var property models.GetProperty

		err := rows.Scan(
			&property.ID,
			&property.Title,
			&property.Location,
			&property.PricePerNight,
			&property.MaxGuests,
			&property.CreatedAt,
			&property.ThumbnailURL,
		)

		if err != nil {
			return properties, err
		}

		properties = append(properties, property)

	}

	if err = rows.Err(); err != nil {
		return properties, err
	}

	return properties, nil
}

func (repo *Repository) GetPropertyByID(id uuid.UUID) (models.Property, error) {
	query1 := `
		SELECT id, title, location, max_guests, price_per_night, description, created_at
		FROM properties WHERE id = $1;
	`

	query2 := `
		SELECT
			id,
			image_url,
			caption,
			display_order
		FROM property_images
		WHERE property_id = $1
		ORDER BY display_order ASC;
	`

	query3 := `
		SELECT
			a.id,
			a.name
		FROM amenities a
		JOIN property_amenities pa ON pa.amenity_id = a.id
		WHERE pa.property_id = $1;
	`

	var property models.Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query1, id).Scan(
		&property.ID,
		&property.Title,
		&property.Location,
		&property.MaxGuests,
		&property.PricePerNight,
		&property.Description,
		&property.CreatedAt,
	)

	if err != nil {
		return property, err
	}

	imgRows, err := repo.db.QueryContext(ctx, query2, id)

	if err != nil {
		return property, err
	}
	defer imgRows.Close()

	for imgRows.Next() {
		var images struct {
			ImageID      uuid.UUID `json:"image_id"`
			ImageURL     string    `json:"image_url"`
			Caption      string    `json:"caption"`
			DisplayOrder int       `json:"display_order"`
		}

		if err := imgRows.Scan(&images.ImageID, &images.ImageURL, &images.Caption, &images.DisplayOrder); err != nil {
			return property, err
		}
		property.Images = append(property.Images, images)
	}

	if err = imgRows.Err(); err != nil {
		return property, err
	}

	amenityRows, err := repo.db.QueryContext(ctx, query3, id)

	if err != nil {
		return property, err
	}
	defer amenityRows.Close()

	for amenityRows.Next() {
		var amenity struct {
			AmenityID uuid.UUID `json:"amenity_id"`
			Name      string    `json:"name"`
		}

		if err := amenityRows.Scan(&amenity.AmenityID, &amenity.Name); err != nil {
			return property, err
		}

		property.Amenities = append(property.Amenities, amenity)
	}

	if err = amenityRows.Err(); err != nil {
		return property, err
	}

	return property, nil
}

func (repo *Repository) PostProperty(property models.PostProperty) (uuid.UUID, error) {
	query := `
		INSERT INTO properties (title, description, location, price_per_night, max_guests, user_id)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id;
	`
	var id uuid.UUID

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query,
		property.Title,
		property.Description,
		property.Location,
		property.PricePerNight,
		property.MaxGuests,
		property.UserID,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *Repository) DeleteProperty(id uuid.UUID) (int64, error) {
	query := `
		DELETE FROM properties
		WHERE id = $1
	`
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	exisitingProperty, err := repo.GetPropertyByID(id)
	if err != nil {
		return 0, err
	}

	if exisitingProperty.ID == uuid.Nil {
		return 0, errors.New("property do not exist")
	}

	result, err := repo.db.ExecContext(ctx, query, id)
	if err != nil {
		return 0, err
	}

	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}

	return count, err
}

func (repo *Repository) UpdateProperty(property models.Property) error {
	query := `
		UPDATE properties
		SET title = $1, description = $2, location = $3, max_guests = $4, price_per_night = $5, updated_at = NOW()
		WHERE id = $6
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx, query,
		property.Title,
		property.Description,
		property.Location,
		property.MaxGuests,
		property.PricePerNight,
		property.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (repo *Repository) PostPropertyImages(data models.AddImagesRequest) error {
	query := `
		INSERT INTO property_images
		(property_id, image_url, caption, display_order)
		VALUES ($1, $2, $3, $4);
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	tx, err := repo.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	for _, img := range data.Images {
		_, err := tx.ExecContext(ctx, query,
			data.PropertyID,
			img.ImageURL,
			img.Caption,
			img.DisplayOrder,
		)
		if err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to insert image: %w", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (repo *Repository) DeletePropertyImage(imageID uuid.UUID) error {
	query := `
		DELETE FROM property_images
		WHERE id = $1;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res, err := repo.db.ExecContext(ctx, query, imageID)
	if err != nil {
		return fmt.Errorf("failed to delete property image: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check delete result: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (repo *Repository) GetAllAmenities() ([]models.Amenity, error) {
	query := `
		SELECT id, name
		FROM amenities
		ORDER BY name ASC;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var amenities []models.Amenity
	for rows.Next() {
		var a models.Amenity
		if err := rows.Scan(&a.AmenityID, &a.Name); err != nil {
			return nil, err
		}
		amenities = append(amenities, a)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return amenities, nil
}

func (repo *Repository) AddAmenity(name string) (*models.Amenity, error) {
	if name == "" {
		return nil, errors.New("amenity name required")
	}

	query := `
		INSERT INTO amenities (name)
		VALUES ($1)
		RETURNING id, name;
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var a models.Amenity
	err := repo.db.QueryRowContext(ctx, query, name).
		Scan(&a.AmenityID, &a.Name)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (repo *Repository) SearchAvailability(searchParams models.SearchPropertyParams) ([]models.GetProperty, error) {
	query := `
		SELECT
			p.id,
			p.title,
			p.location,
			p.price_per_night,
			p.max_guests,
			p.created_at,
			(
				SELECT pi.image_url
				FROM property_images pi
				WHERE pi.property_id = p.id
				ORDER BY pi.display_order ASC
				LIMIT 1
			) AS thumbnail_url
		FROM properties p
		WHERE LOWER(p.location) LIKE LOWER($1)
		AND NOT EXISTS (
				SELECT 1 FROM bookings b
				WHERE b.property_id = p.id
				AND daterange(b.start_date, b.end_date, '[]')
					&& daterange($2, $3, '[]')
		);

	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var props []models.GetProperty

	rows, err := repo.db.QueryContext(ctx, query, searchParams.Location, searchParams.StartDate, searchParams.EndDate)
	if err != nil {
		return props, err
	}
	defer rows.Close()

	for rows.Next() {
		var p models.GetProperty
		if err := rows.Scan(
			&p.ID, &p.Title,
			&p.Location, &p.PricePerNight, &p.MaxGuests, &p.CreatedAt, &p.ThumbnailURL,
		); err != nil {
			return props, err
		}
		props = append(props, p)
	}

	return props, nil
}

func (repo *Repository) PostPropertyAmenity(amenities []models.PostAmenity) error {
	if len(amenities) == 0 {
		return nil
	}

	// Build the query with multiple value placeholders
	valueStrings := make([]string, 0, len(amenities))
	valueArgs := make([]interface{}, 0, len(amenities)*2)

	for i, amenity := range amenities {
		valueStrings = append(valueStrings, fmt.Sprintf("($%d, $%d)", i*2+1, i*2+2))
		valueArgs = append(valueArgs, amenity.PropertyID, amenity.AmenityID)
	}

	query := fmt.Sprintf(`
        INSERT INTO property_amenities (property_id, amenity_id)
        VALUES %s
        ON CONFLICT (property_id, amenity_id) DO NOTHING`,
		strings.Join(valueStrings, ","))

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx, query, valueArgs...)
	return err

}
