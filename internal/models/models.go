package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	SecondName   string    `json:"second_name"`
	Email        string    `json:"email"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	PasswordHash string    `json:"password_hash"`
}

type RegisterUser struct {
	FirstName    string `json:"first_name"`
	SecondName   string `json:"second_name"`
	Email        string `json:"email"`
	PasswordHash string `json:"password_hash"`
}

type LoginUser struct {
	ID           uuid.UUID `json:"id"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"password_hash"`
	Role         string    `json:"role"`
}

type Property struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Location      string    `json:"location"`
	Description   string    `json:"description"`
	PricePerNight float32   `json:"price_per_night"`
	MaxGuests     int       `json:"max_guests"`
	ImageURL      string    `json:"image_url"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

type PostProperty struct {
	Title         string    `json:"title"`
	Location      string    `json:"location"`
	Description   string    `json:"description"`
	PricePerNight float32   `json:"price_per_night"`
	MaxGuests     int       `json:"max_guests"`
	ImageURL      string    `json:"image_url"`
	UserID        uuid.UUID `json:"user_id"`
}

type Booking struct {
	ID       uuid.UUID `json:"id"`
	Property struct {
		Title    string `json:"title"`
		Location string `json:"location"`
	}
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float32       `json:"total_price"`
	Status     string    `json:"status"`
}

type GetBooking struct {
	ID       uuid.UUID `json:"id"`
	Property struct {
		Title    string `json:"title"`
		Location string `json:"location"`
	}
	FirstName  string    `json:"first_name"`
	LastName   string    `json:"last_name"`
	StartDate  time.Time `json:"start_date"`
	EndDate    time.Time `json:"end_date"`
	TotalPrice float32       `json:"total_price"`
	Status     string    `json:"status"`
}
