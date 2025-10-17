package handler

import (
	"fmt"
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/config"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/helper"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/repository"
)

type envelope map[string]any

// type contextKey string

// const userIDKey contextKey = "user_id"

type Handler struct {
	cfg  *config.Config
	repo *repository.Repository
}

func NewHandler(cfg *config.Config, repo *repository.Repository) *Handler {
	return &Handler{
		cfg:  cfg,
		repo: repo,
	}
}

func (h *Handler) Healthz(w http.ResponseWriter, r *http.Request) {
	data := envelope{
		"status": "available",
		"system_info": map[string]string{
			"environment": h.cfg.Env,
			"version":     "v1",
		},
	}

	if err := helper.WriteJSON(w, data, http.StatusOK); err != nil {
		http.Error(w, fmt.Sprintf("write json: %v", err), http.StatusInternalServerError)
	}
}

func (h *Handler) Testing(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	if userVal == nil {
		http.Error(w, "user not found in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userVal.(string)
	if !ok {
		http.Error(w, "invalid user id type in context", http.StatusInternalServerError)
		return
	}

	helper.WriteJSON(w, userID, http.StatusAccepted)
}
