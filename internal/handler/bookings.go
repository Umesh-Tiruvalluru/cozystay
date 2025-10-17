package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/Umesh-Tiruvalluru/BookBnb/internal/helper"
	"github.com/google/uuid"
)

func (h *Handler) GetBookings(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")

	if userVal == "" {
		h.cfg.Logger.Error("user not found in contet")
		http.Error(w, "User not in the context", http.StatusUnauthorized)
		return
	}

	userID := uuid.MustParse(userVal.(string))

	bookings, err := h.repo.GetBookings(userID)
	if err != nil {
		h.cfg.Logger.Error("Unable to get all the booking", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, bookings, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	bookingID := uuid.MustParse(r.PathValue("id"))

	res, err := h.repo.GetBookingByID(bookingID)
	if err != nil {
		h.cfg.Logger.Error("Unable to get booking", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w, res, http.StatusCreated); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}
}

func (h *Handler) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userVal := r.Context().Value("userID")

	if userVal == nil {
		h.cfg.Logger.Error("Can not find user id", "Error", "User is not in the context")
		http.Error(w, "User is not in the context", http.StatusUnauthorized)
		return
	}

	userID := uuid.MustParse(userVal.(string))

	fmt.Println(userID)

	var req struct {
		PropertyID uuid.UUID `json:"property_id"`
		StartDate  string    `json:"start_date"`
		EndDate    string    `json:"end_date"`
		TotalPrice int       `json:"total_price"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.cfg.Logger.Error("Invalid JSON body", "Error", err)
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}

	startDate, err1 := time.Parse("2006-01-02", req.StartDate)
	endDate, err2 := time.Parse("2006-01-02", req.EndDate)
	if err1 != nil || err2 != nil {
		h.cfg.Logger.Error("Invalid date format", "Error", err1, "Error", err2)
		http.Error(w, "Invalid date format, use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	if !endDate.After(startDate) {
		h.cfg.Logger.Error("End date must be after start date")
		http.Error(w, "End date must be after start date", http.StatusBadRequest)
		return
	}

	res, err := h.repo.CreateBooking(userID, req.PropertyID, startDate, endDate, req.TotalPrice)
	if err != nil {
		h.cfg.Logger.Error("Unable to create booking", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	helper.WriteJSON(w, res, http.StatusAccepted)
}

func (h *Handler) CancelBooking(w http.ResponseWriter, r *http.Request) {
	bookingID := uuid.MustParse(r.PathValue("id"))
	userVal := r.Context().Value("userID")

	if userVal == nil {
		h.cfg.Logger.Error("Can not find user id", "Error", "User is not in the context")
		http.Error(w, "User and role are not in the context", http.StatusUnauthorized)
		return
	}

	userID := uuid.MustParse(userVal.(string))

	res, err := h.repo.CancelBooking(bookingID, userID)
	if err != nil {
		h.cfg.Logger.Error("Failed Cancellation", "Error", err)
		http.Error(w, fmt.Sprintf("Error:%v", err), http.StatusConflict)
		return
	}

	if err := helper.WriteJSON(w,res, http.StatusNoContent); err != nil {
		h.cfg.Logger.Error("Failed to generate a response", "Error", err)
		http.Error(w, "Failed to generate a response", http.StatusInternalServerError)
	}


}
