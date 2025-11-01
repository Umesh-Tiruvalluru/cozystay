package main

import (
	"context"
	"net/http"
	"strings"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/helper"
	"github.com/go-chi/chi/v5/middleware"
)

func RoleMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := helper.VerifyToken(parts[1], cfg.JwtSecret)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "role", token["role"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")

		if authHeader == "" {
			http.Error(w, "Missing authorization header", http.StatusUnauthorized)
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			http.Error(w, "Invalid Authorization header format", http.StatusUnauthorized)
			return
		}

		token, err := helper.VerifyToken(parts[1], cfg.JwtSecret)
		if err != nil {
			http.Error(w, "Invalid Authorization", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", token["userID"])

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Wrap response writer to capture status code
		ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

		// Process request
		next.ServeHTTP(ww, r)

		// Log the request
		duration := time.Since(start)

		userID := "anonymous"
		if ctxUserID := r.Context().Value("userID"); ctxUserID != nil {
			if id, ok := ctxUserID.(string); ok {
				userID = id
			}
		}

		cfg.Logger.RequestInfo(
			r.Method,
			r.URL.Path,
			r.RemoteAddr,
			ww.Status(),
			duration,
			userID,
		)
	})
}
