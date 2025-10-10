package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"massage-booking/backend/database"
)

// GetMassageTypesHandler handles GET /api/massage-types
func GetMassageTypesHandler(w http.ResponseWriter, r *http.Request) {
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

	// Get massage types from database
	massageTypes, err := database.GetMassageTypes()
	if err != nil {
		log.Printf("Error getting massage types: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Encode and send response
	if err := json.NewEncoder(w).Encode(massageTypes); err != nil {
		log.Printf("Error encoding massage types: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully returned %d massage types", len(massageTypes))
}
