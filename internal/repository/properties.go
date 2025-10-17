package repository

import (
	"context"
	"errors"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)



func (repo *Repository) GetAllProperties() ([]models.Property, error) {
	query := `
		SELECT id, title, location, max_guests, price_per_night, description, image_url
		FROM properties;
	`

	var properties []models.Property
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, query)

	if err != nil {
		return properties, err
	}

	for rows.Next() {
		var property models.Property

		err := rows.Scan(
			&property.ID,
			&property.Title,
			&property.Location,
			&property.MaxGuests,
			&property.PricePerNight,
			&property.Description,
			&property.ImageURL,
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
	query := `
		SELECT id, title, location, max_guests, price_per_night, description, image_url
		FROM properties WHERE ID = $1;
	`
	var property models.Property

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&property.ID,
		&property.Title,
		&property.Location,
		&property.MaxGuests,
		&property.PricePerNight,
		&property.Description,
		&property.ImageURL,
	)

	if err != nil {
		return property, err
	}

	return property, nil

}

func (repo *Repository) PostProperty(property models.PostProperty) (uuid.UUID, error) {
	query := `
		INSERT INTO properties (title, description, location, price_per_night, max_guests, image_url, user_id)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
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
		property.ImageURL,
		property.UserID,
	).Scan(&id)

	if err != nil {
		return uuid.Nil, errors.New("cannot query to databse")
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
		SET title = $1, description = $2, location = $3, max_guests = $4, price_per_night = $5, image_url = $6, updated_at = NOW()
		WHERE id = $7
	`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := repo.db.ExecContext(ctx, query,
		property.Title,
		property.Description,
		property.Location,
		property.MaxGuests,
		property.PricePerNight,
		property.ImageURL,
		property.ID,
	)

	if err != nil {
		return err
	}

	return nil
}
