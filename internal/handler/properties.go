package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/helper"
	"github.com/Umesh-Tiruvalluru/BookBnb/internal/models"
	"github.com/google/uuid"
)

func (h *Handler) GetAllProperties(w http.ResponseWriter, r *http.Request) {
	properties, err := h.repo.GetAllProperties()
	if err != nil {
		h.cfg.Logger.Error("Unable to get all properties", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, properties, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetPropertyByID(w http.ResponseWriter, r *http.Request) {
	id := uuid.MustParse(r.PathValue("id"))

	property, err := h.repo.GetPropertyByID(id)
	if err != nil {
		h.cfg.Logger.Error("Unable to get a property", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, property, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) PostProperty(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	roleVal := r.Context().Value("role")

	if userVal == nil && roleVal == nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", "User and role are not in the context")
		http.Error(w, "User and role are not in the context", http.StatusUnauthorized)
		return
	}

	userID, ok := userVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid user id type")
		http.Error(w, "Invalid user id type in context", http.StatusInternalServerError)
		return
	}
	
	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid role type")
		http.Error(w, "invalid role type in context", http.StatusInternalServerError)
		return
	}
	
	if role != "admin" {
		h.cfg.Logger.Error("Forbidden", "Error", "The authenticated user lacks the necessary permission")
		http.Error(w, "User do not have necessary permissions", http.StatusForbidden)
		return
	}

	var property models.PostProperty

	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "Error", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	property.UserID = uuid.MustParse(userID)

	id, err := h.repo.PostProperty(property)
	if err != nil {
		h.cfg.Logger.Error("Unable to post a property", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	message := map[string]any{
		"id":      id,
		"message": "Successfully created a property",
	}

	if err := helper.WriteJSON(w, message, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

//error
func (h *Handler) UpdateProperty(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value("role")
	userVal := r.Context().Value("userID")

	if userVal == nil && roleVal == nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", "User and role are not in the context")
		http.Error(w, "User and role are not in the context", http.StatusUnauthorized)
		return
	}

	userID, ok := userVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid user id type")
		http.Error(w, "Invalid user id type in context", http.StatusInternalServerError)
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid role type")
		http.Error(w, "invalid role type in context", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		h.cfg.Logger.Error("Forbidden", "Error", "The authenticated user lacks the necessary permission")
		http.Error(w, "User do not have necessary permissions", http.StatusForbidden)
		return
	}

	var property models.Property

	if err := json.NewDecoder(r.Body).Decode(&property); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "Error", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	err := h.repo.UpdateProperty(property)
	if err != nil {
		h.cfg.Logger.Error("Unable to update a property", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	message := map[string]any{
 		"message": "Updated successfully", 
		"userID": userID,
	}

	if err := helper.WriteJSON(w, message, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}

}

func (h *Handler) DeleteProperty(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value("role")
	id := uuid.MustParse(r.PathValue("id"))

	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid type", "Error", "Invalid role type")
		http.Error(w, "invalid role type in context", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		h.cfg.Logger.Error("Forbidden", "Error", "The authenticated user lacks the necessary permission")
		http.Error(w, "User do not have necessary permissions", http.StatusForbidden)
		return
	}

	_, err := h.repo.DeleteProperty(id)
	if err != nil {
		h.cfg.Logger.Error("Unable to delete a property", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}
	

	if err := helper.WriteJSON(w, map[string]any{"message": "Successfully Deleted"}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}

}
