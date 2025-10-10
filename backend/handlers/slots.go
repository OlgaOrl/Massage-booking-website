package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"massage-booking/backend/database"
)

// GetSlotsHandler handles GET /api/slots?date=YYYY-MM-DD&service_id=1
func GetSlotsHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get query parameters
	date := r.URL.Query().Get("date")
	serviceIDStr := r.URL.Query().Get("service_id")

	// Validate required parameters
	if date == "" {
		http.Error(w, "Missing required parameter: date", http.StatusBadRequest)
		return
	}

	if serviceIDStr == "" {
		http.Error(w, "Missing required parameter: service_id", http.StatusBadRequest)
		return
	}

	// Parse service_id
	serviceID, err := strconv.Atoi(serviceIDStr)
	if err != nil {
		http.Error(w, "Invalid service_id parameter", http.StatusBadRequest)
		return
	}

	// Get time slots from database
	timeSlots, err := database.GetTimeSlots(date, serviceID)
	if err != nil {
		log.Printf("Error getting time slots for date %s and service %d: %v", date, serviceID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode and send response
	if err := json.NewEncoder(w).Encode(timeSlots); err != nil {
		log.Printf("Error encoding time slots: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned %d time slots for date %s and service %d", len(timeSlots), date, serviceID)
}
