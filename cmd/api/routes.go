package main

import (
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/handler"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	r := chi.NewRouter()

	r.Use(LoggingMiddleware)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete, http.MethodPatch, http.MethodOptions},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	repo := repository.NewRepositoryUser(db)
	h := handler.NewHandler(&cfg, repo)

	api := chi.NewRouter()

	// --- Auth routes ---
	api.Route("/auth", func(r chi.Router) {
		r.Post("/login", h.Login)
		r.Post("/register", h.Register)
		// protected route returning current user
		r.With(AuthMiddleware).Get("/me", h.Me)
	})

	// --- Properties ---
	api.Route("/properties", func(r chi.Router) {
		// public
		r.Get("/", h.GetAllProperties)
		r.Get("/{id}", h.GetPropertyByID)
		r.Get("/{id}/availability", h.SearchAvailability)

		// protected: only users with appropriate role (e.g. host/admin)
		r.With(AuthMiddleware, RoleMiddleware).Post("/", h.PostProperty)
		r.With(AuthMiddleware, RoleMiddleware).Put("/{id}", h.UpdateProperty)
		r.With(AuthMiddleware, RoleMiddleware).Delete("/{id}", h.DeleteProperty)

		// property images
		r.With(AuthMiddleware, RoleMiddleware).Post("/{id}/images", h.PostImage)
		r.With(AuthMiddleware, RoleMiddleware).Delete("/{id}/images/{imageID}", h.DeletePropertyImage)
	})

	// --- Amenities ---
	api.Route("/amenities", func(r chi.Router) {
		// anyone can list amenities
		r.Get("/", h.GetAmenities)
		// only admins/hosts can create or attach
		r.With(AuthMiddleware, RoleMiddleware).Post("/", h.AddAmenity)
		r.With(AuthMiddleware, RoleMiddleware).Post("/{propertyID}", h.PostPropertyAmenities)
	})

	// --- Bookings ---
	api.Route("/bookings", func(r chi.Router) {
		// user must be authenticated to access bookings
		r.With(AuthMiddleware).Get("/", h.GetBookings)
		r.With(AuthMiddleware).Post("/", h.CreateBooking)
		r.With(AuthMiddleware).Get("/{id}", h.GetBookingByID)
		// partial update for status changes (cancel, check-in, etc.)
		r.With(AuthMiddleware).Patch("/{id}", h.CancelBooking)
	})

	// health check
	api.Get("/healthz", h.Healthz)

	// mount versioned API
	r.Mount("/api/v1", api)

	return r
}
