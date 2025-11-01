package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/helper"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)

func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
	var req models.User

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "Error", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.PasswordHash == "" || req.FirstName == "" {
		h.cfg.Logger.Error("first_name, email and password are required")
		http.Error(w, "first_name, email and password are required", http.StatusBadRequest)
		return
	}

	pwHash, err := helper.HashPassword(req.PasswordHash)
	if err != nil {
		http.Error(w, "unable to hash password", http.StatusBadRequest)
		return
	}

	user := models.RegisterUser{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		Email:        req.Email,
		PasswordHash: pwHash,
	}

	id, err := h.repo.RegisterUser(&user)
	if err != nil {
		h.cfg.Logger.Error("User registration failed", "Error", err)
		http.Error(w, "User registration failed", http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, envelope{"id": id}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginUser
	var usr models.LoginUser

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "Error", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	usr, err := h.repo.LoginUser(req.Email)
	if err != nil {
		h.cfg.Logger.Error("Internal Server Error", "Error", err)
		http.Error(w, "Internal Server Error: Database connection failed during login attempt.", http.StatusInternalServerError)
		return
	}

	CheckPassword := helper.CheckPasswordHash(req.PasswordHash, usr.PasswordHash)

	if !CheckPassword {
		h.cfg.Logger.Error("Password do not match")
		http.Error(w, "Passwords do not match.", http.StatusBadRequest)
		return
	}

	token, err := helper.CreateToken(usr.ID, usr.Role, h.cfg.JwtSecret)
	if err != nil {
		h.cfg.Logger.Error("Error occurred while creating token", "Error", err)
		http.Error(w, "Error occurred while creating token", http.StatusInternalServerError)
	}

	if err := helper.WriteJSON(w, map[string]any{"token": token}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}

}

func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	if userVal == nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", "User is not in the context")
		http.Error(w, "User is not in the context", http.StatusUnauthorized)
		return
	}

	userID, ok := userVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid user id type")
		http.Error(w, "Invalid user id type in context", http.StatusInternalServerError)
		return
	}

	userDetails, err := h.repo.UserDetails(uuid.MustParse(userID))
	if err != nil {
		h.cfg.Logger.Error("Unable to get user details", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, map[string]any{"user": userDetails}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
		return
	}
}
