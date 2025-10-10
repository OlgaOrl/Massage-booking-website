package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"massage-booking/backend/database"
	"massage-booking/backend/models"
)

// CreateReservation handles POST /api/reservations
func CreateReservation(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow POST method
	if r.Method != "POST" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse request body
	var req models.ReservationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error parsing reservation request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate slot_id
	if req.SlotID <= 0 {
		http.Error(w, "Invalid slot_id", http.StatusBadRequest)
		return
	}

	// Create reservation
	reservationID, expiresAt, err := database.CreateReservation(req.SlotID)
	if err != nil {
		log.Printf("Error creating reservation for slot %d: %v", req.SlotID, err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Slot not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "not available") || strings.Contains(err.Error(), "already reserved") {
			http.Error(w, "Slot is not available", http.StatusConflict)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Calculate expires in seconds
	expiresInSeconds := int(time.Until(expiresAt).Seconds())

	// Create response
	response := models.ReservationResponse{
		ReservationID:    reservationID,
		ExpiresAt:        expiresAt,
		ExpiresInSeconds: expiresInSeconds,
	}

	// Send response
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding reservation response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Created reservation %d for slot %d, expires at %v", reservationID, req.SlotID, expiresAt)
}

// DeleteReservation handles DELETE /api/reservations/:id
func DeleteReservation(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow DELETE method
	if r.Method != "DELETE" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract reservation ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/reservations/")
	if path == "" {
		http.Error(w, "Missing reservation ID", http.StatusBadRequest)
		return
	}

	reservationID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid reservation ID", http.StatusBadRequest)
		return
	}

	// Delete reservation
	if err := database.DeleteReservation(reservationID); err != nil {
		log.Printf("Error deleting reservation %d: %v", reservationID, err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Reservation not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Return 204 No Content
	w.WriteHeader(http.StatusNoContent)
	log.Printf("Deleted reservation %d", reservationID)
}
