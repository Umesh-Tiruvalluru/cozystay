package main

import (
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/handler"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
)

func routes() http.Handler {
	// mux := chi.NewRouter()
	v1Router := chi.NewRouter()
	repo := repository.NewRepositoryUser(db)
	h := handler.NewHandler(&cfg, repo)

	v1Router.Use(LoggingMiddleware)

	v1Router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://*", "https://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router.Get("/healthzz", h.Healthz)
	// v1Router.With(AuthMiddleware).Get("/testing", h.Testing)

	// auth
	v1Router.Post("/auth/login", h.Login)
	v1Router.Post("/auth/register", h.Register)
	v1Router.With(AuthMiddleware).Get("/auth/me", h.Me)

	//Properties :
	v1Router.Get("/properties", h.GetAllProperties)
	v1Router.Get("/properties/{id}", h.GetPropertyByID)
	v1Router.With(AuthMiddleware, RoleMiddleware).Post("/properties", h.PostProperty)
	v1Router.With(AuthMiddleware, RoleMiddleware).Delete("/properties/{id}", h.DeleteProperty)
	v1Router.With(AuthMiddleware, RoleMiddleware).Put("/properties", h.UpdateProperty)
	v1Router.Get("/properties/{id}/availability", h.SearchAvailability)

	//property image
	v1Router.With(AuthMiddleware, RoleMiddleware).Post("/property/image/{id}", h.PostImage)
	v1Router.With(AuthMiddleware, RoleMiddleware).Delete("/property/image/{id}", h.DeletePropertyImage)

	//Amenity:
	v1Router.With(AuthMiddleware, RoleMiddleware).Post("/amenities", h.AddAmenity)
	v1Router.With(AuthMiddleware, RoleMiddleware).Get("/amenities", h.GetAmenities)
	v1Router.With(AuthMiddleware, RoleMiddleware).Post("/properties/{id}/amenities", h.PostPropertyAmenities)
	// v1Router.With(AuthMiddleware, RoleMiddleware).Delete("/amenities", )

	//Bookings
	v1Router.With(AuthMiddleware).Get("/bookings", h.GetBookings)
	v1Router.With(AuthMiddleware).Post("/bookings", h.CreateBooking)
	v1Router.With(AuthMiddleware).Get("/bookings/{id}", h.GetBookingByID)
	v1Router.With(AuthMiddleware).Patch("/bookings/{id}", h.CancelBooking)

	return v1Router
}
