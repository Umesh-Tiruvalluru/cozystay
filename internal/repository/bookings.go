package repository

import (
	"context"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)



func (repo *Repository) GetBookings(userID uuid.UUID) ([]models.Booking, error) {
	query := `
		SELECT b.id, b.start_date, b.end_date, p.title, p.location, b.total_price, b.status
		FROM bookings b
		LEFT JOIN properties p ON b.property_id = p.id
		WHERE b.user_id = $1;
	`

	var bookings []models.Booking

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	rows, err := repo.db.QueryContext(ctx, query, userID)

	if err != nil {
		return bookings, err
	}

	for rows.Next() {
		var booking models.Booking

		err := rows.Scan(
			&booking.ID,
			&booking.StartDate,
			&booking.EndDate,
			&booking.Property.Title,
			&booking.Property.Location,
			&booking.TotalPrice,
			&booking.Status,
		)
		if err != nil {
			return bookings, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return bookings, err
	}

	return bookings, nil

}

func (repo *Repository) CreateBooking(userId uuid.UUID, propertyID uuid.UUID, startDate time.Time, endDate time.Time, totalPrice int) (uuid.UUID, error) {
	query := `
		INSERT INTO bookings (user_id, property_id, start_date, end_date, total_price)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id;
	`
	var id uuid.UUID

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, userId, propertyID, startDate, endDate, totalPrice).Scan(&id)

	if err != nil {
		return uuid.Nil, err
	}

	return id, nil
}

func (repo *Repository) GetBookingByID(id uuid.UUID) (models.GetBooking, error) {
	query := `
		SELECT b.id, b.start_date, b.end_date, b.total_price, b.status, p.title, p.location, u.first_name, u.last_name
		FROM bookings b
		LEFT JOIN properties p ON b.property_id = p.id
		LEFT JOIN users u ON b.user_id = u.id
		WHERE b.id = $1;
	`

	var booking models.GetBooking

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.StartDate,
		&booking.EndDate,
		&booking.TotalPrice,
		&booking.Status,
		&booking.Property.Title,
		&booking.Property.Location,
		&booking.FirstName,
		&booking.LastName,
	)

	if err != nil {
		return booking, err
	}

	return booking, nil
}

func (repo *Repository) CancelBooking(id uuid.UUID, userID uuid.UUID) (string, error) {
	query := `
		UPDATE bookings 
		SET status = 'cancelled'
		WHERE user_id = $1 AND id = $2
		RETURNING status;
	`

	var status string
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := repo.db.QueryRowContext(ctx, query, userID, id).Scan(&status)

	if err != nil {
		return "", err
	}

	return status, nil
}
