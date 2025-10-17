package main

import (
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/handler"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/repository"
	"github.com/go-chi/chi/v5"
)

func routes() http.Handler {
	// mux := chi.NewRouter()
	v1Router := chi.NewRouter()
	repo := repository.NewRepositoryUser(db)
	h := handler.NewHandler(&cfg, repo)

	v1Router.Use(LoggingMiddleware)

	v1Router.Get("/healthzz", h.Healthz)
	v1Router.With(AuthMiddleware).Get("/testing", h.Testing)
	// v1Router.With(AuthMiddleware).Get("/me", h.Me)

	// auth
	v1Router.Post("/auth/login", h.Login)
	v1Router.Post("/auth/register", h.Register)

	//Properties :
	// TODO : Accept as Params - Post, Put Property,
	v1Router.Get("/properties", h.GetAllProperties)
	v1Router.Get("/properties/{id}", h.GetPropertyByID)
	v1Router.With(AuthMiddleware, RoleMiddleware).Post("/properties", h.PostProperty)
	v1Router.With(AuthMiddleware, RoleMiddleware).Delete("/properties/{id}", h.DeleteProperty)
	v1Router.With(AuthMiddleware, RoleMiddleware).Put("/properties", h.UpdateProperty)
	v1Router.Get("/properties/{id}/availability", h.SearchAvailability)

	//Bookings
	v1Router.With(AuthMiddleware).Get("/bookings", h.GetBookings)
	v1Router.With(AuthMiddleware).Post("/bookings", h.CreateBooking)
	v1Router.With(AuthMiddleware).Get("/bookings/{id}", h.GetBookingByID)
	v1Router.With(AuthMiddleware).Patch("/bookings/{id}", h.CancelBooking)

	return v1Router
}
