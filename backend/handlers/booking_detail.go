package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"massage-booking/backend/database"
)

// GetBooking handles GET /api/bookings/:id
func GetBooking(w http.ResponseWriter, r *http.Request) {
	// Set CORS headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "application/json")

	// Handle preflight OPTIONS request
	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusOK)
		return
	}

	// Only allow GET method
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Extract booking ID from URL path
	path := strings.TrimPrefix(r.URL.Path, "/api/bookings/")
	if path == "" {
		http.Error(w, "Missing booking ID", http.StatusBadRequest)
		return
	}

	bookingID, err := strconv.Atoi(path)
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	// Get booking from database
	booking, err := database.GetBookingByID(bookingID)
	if err != nil {
		log.Printf("Error getting booking %d: %v", bookingID, err)
		if strings.Contains(err.Error(), "not found") {
			http.Error(w, "Booking not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Send response
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned booking details for ID %d (reference: %s)", bookingID, booking.Reference)
}
