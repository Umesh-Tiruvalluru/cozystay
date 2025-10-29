package models

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	FirstName    string    `json:"first_name"`
	LastName   string      `json:"last_name"`
	Email        string    `json:"email"`
	CreatedAt    string    `json:"created_at"`
	UpdatedAt    string    `json:"updated_at"`
	PasswordHash string    `json:"password_hash"`
}


type UserDetails struct {
	ID uuid.UUID `json:"id"`
	FirstName string `json:"first_name"`
	LastName string `json:"last_name"`
	Email string `json:"email"`
	Role string `json:"role"`
}

type RegisterUser struct {
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
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
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	Images []struct {
		ImageID 	 uuid.UUID `json:"image_id"`
		ImageURL     string    `json:"image_url"`
		Caption      string    `json:"caption"`
		DisplayOrder int       `json:"display_order"`
	} `json:"images"`
	Amenities []struct {
		AmenityID uuid.UUID `json:"amenity_id"`
		Name 	  string    `json:"name"`
	} `json:"amenities"`
}

type GetProperty struct {
	ID            uuid.UUID `json:"id"`
	Title         string    `json:"title"`
	Location      string    `json:"location"`
	PricePerNight float32   `json:"price_per_night"`
	MaxGuests     int       `json:"max_guests"`
	CreatedAt     time.Time `json:"created_at"`
	ThumbnailURL  sql.NullString    `json:"thumbnail_url"`
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


type AddImagesRequest struct {
	PropertyID string `json:"-"` // filled from path param
	Images     []struct {
		ImageURL     string `json:"image_url"`
		Caption      string `json:"caption"`
		DisplayOrder int    `json:"display_order"`
	} `json:"images"`
}

type Amenity struct {
	AmenityID uuid.UUID `json:"amenity_id"`
	Name      string    `json:"name"`
}

type PostAmenity struct {
	AmenityID  uuid.UUID `json:"amenity_id"`
	PropertyID uuid.UUID `json:"property_id"`
}

type AddAmenityRequest struct {
	Name string `json:"name"`
}

type AddAmenitiesRequest struct {
	AmenityID []uuid.UUID `json:"amenity_id"`
}


type SearchPropertyParams struct {
	Location  string   `json:"location"`
	StartDate string   `json:"start_date"`
	EndDate   string   `json:"end_date"`
	MinPrice  *float64 `json:"min_price,omitempty"`
	MaxPrice  *float64 `json:"max_price,omitempty"`
	Guests    *int     `json:"guests,omitempty"`
}