package handler

import (
	"database/sql"
	"encoding/json"
	"errors"
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
		"userID":  userID,
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

func (h *Handler) PostImage(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	roleVal := r.Context().Value("role")

	if userVal == nil || roleVal == nil {
		h.cfg.Logger.Error("Missing user or role in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid role type in context")
		http.Error(w, "invalid role type", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		h.cfg.Logger.Error("Forbidden: insufficient permission", "role", role)
		http.Error(w, "forbidden: only admin can add images", http.StatusForbidden)
		return
	}

	propertyIDParam := r.PathValue("id")
	propertyID, err := uuid.Parse(propertyIDParam)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	var req models.AddImagesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "error", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if len(req.Images) == 0 {
		h.cfg.Logger.Error("No images provided")
		http.Error(w, "no images provided", http.StatusBadRequest)
		return
	}

	req.PropertyID = propertyID.String()

	if err := h.repo.PostPropertyImages(req); err != nil {
		h.cfg.Logger.Error("Failed to add property images", "error", err)
		http.Error(w, "failed to add images", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, map[string]any{"message": "Successfully image added"}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}

}

func (h *Handler) DeletePropertyImage(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	roleVal := r.Context().Value("role")

	if userVal == nil || roleVal == nil {
		h.cfg.Logger.Error("Missing user or role in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid role type in context")
		http.Error(w, "invalid role type", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		h.cfg.Logger.Error("Forbidden: insufficient permission", "role", role)
		http.Error(w, "forbidden: only admin can delete images", http.StatusForbidden)
		return
	}

	imageIDParam := r.PathValue("id")

	imageID, err := uuid.Parse(imageIDParam)
	if err != nil {
		http.Error(w, "invalid image id", http.StatusBadRequest)
		return
	}

	err = h.repo.DeletePropertyImage(imageID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "image not found", http.StatusNotFound)
			return
		}

		h.cfg.Logger.Error("Failed to delete property image", "error", err)
		http.Error(w, "failed to delete image", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, map[string]any{"message": "Successfully image deleted"}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetAmenities(w http.ResponseWriter, r *http.Request) {
	amenities, err := h.repo.GetAllAmenities()
	if err != nil {
		h.cfg.Logger.Error("Failed to get amenities", "error", err)
		http.Error(w, "failed to fetch amenities", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, amenities, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) AddAmenity(w http.ResponseWriter, r *http.Request) {
	roleVal := r.Context().Value("role")

	if roleVal == nil {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := roleVal.(string)
	if !ok || role != "admin" {
		http.Error(w, "forbidden", http.StatusForbidden)
		return
	}

	var req models.AddAmenityRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if req.Name == "" {
		http.Error(w, "name is required", http.StatusBadRequest)
		return
	}

	newAmenity, err := h.repo.AddAmenity(req.Name)
	if err != nil {
		h.cfg.Logger.Error("Failed to add amenity", "error", err)
		http.Error(w, "failed to add amenity", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, newAmenity, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) SearchAvailability(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	location := q.Get("location")
	start_date := q.Get("startDate")
	end_date := q.Get("endDate")

	if start_date == "" || end_date == "" {
		h.cfg.Logger.Error("start date and end date are required")
		http.Error(w, "start date and end date are required", http.StatusBadRequest)
		return
	}

	searchParams := models.SearchPropertyParams{
		Location:  location,
		StartDate: start_date,
		EndDate:   end_date,
	}

	data, err := h.repo.SearchAvailability(searchParams)

	if err != nil {
		h.cfg.Logger.Error("Failed to search properties", "error", err)
		http.Error(w, "failed to search properties", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, data, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}

}

func (h *Handler) PostPropertyAmenities(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")
	roleVal := r.Context().Value("role")

	if userVal == nil || roleVal == nil {
		h.cfg.Logger.Error("Missing user or role in context")
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}

	role, ok := roleVal.(string)
	if !ok {
		h.cfg.Logger.Error("Invalid role type in context")
		http.Error(w, "invalid role type", http.StatusInternalServerError)
		return
	}

	if role != "admin" {
		h.cfg.Logger.Error("Forbidden: insufficient permission", "role", role)
		http.Error(w, "forbidden: only admin can add amenities", http.StatusForbidden)
		return
	}

	propertyIDParam := r.PathValue("propertyID")
	propertyID, err := uuid.Parse(propertyIDParam)
	if err != nil {
		http.Error(w, "invalid property id", http.StatusBadRequest)
		return
	}

	var req models.AddAmenitiesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "error", err)
		http.Error(w, "invalid JSON body", http.StatusBadRequest)
		return
	}

	if len(req.AmenityID) == 0 {
		h.cfg.Logger.Error("No amenities provided")
		http.Error(w, "no amenities provided", http.StatusBadRequest)
		return
	}

	// Convert to the repository model format
	amenities := make([]models.PostAmenity, len(req.AmenityID))
	for i, amenityID := range req.AmenityID {
		amenities[i] = models.PostAmenity{
			PropertyID: propertyID,
			AmenityID:  amenityID,
		}
	}

	if err := h.repo.PostPropertyAmenity(amenities); err != nil {
		h.cfg.Logger.Error("Failed to add property amenities", "error", err)
		http.Error(w, "failed to add amenities", http.StatusInternalServerError)
		return
	}

	if err := helper.WriteJSON(w, map[string]any{"message": "Successfully added amenities"}, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}
