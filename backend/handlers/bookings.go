package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"massage-booking/backend/database"
	"massage-booking/backend/models"
)

// CreateBooking handles POST /api/bookings
func CreateBooking(w http.ResponseWriter, r *http.Request) {
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
	var req models.BookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error parsing booking request: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate request fields
	if err := validateBookingRequest(req); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Check if reservation exists and is not expired
	var slotID int
	var expiresAt string
	err := database.DB.QueryRow(`
		SELECT slot_id, expires_at 
		FROM temporary_reservations 
		WHERE id = ? AND expires_at > datetime('now')
	`, req.ReservationID).Scan(&slotID, &expiresAt)
	
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Reservation not found or expired", http.StatusNotFound)
			return
		}
		log.Printf("Error checking reservation %d: %v", req.ReservationID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Start database transaction
	tx, err := database.DB.Begin()
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Create booking
	result, err := tx.Exec(`
		INSERT INTO bookings (client_name, email, phone, service_id, date, time_slot) 
		VALUES (?, ?, ?, ?, ?, ?)
	`, req.ClientName, req.Email, req.Phone, req.ServiceID, req.Date, req.TimeSlot)
	
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		log.Printf("Error getting booking ID: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Mark slot as unavailable
	_, err = tx.Exec("UPDATE time_slots SET available = 0 WHERE id = ?", slotID)
	if err != nil {
		log.Printf("Error marking slot %d as unavailable: %v", slotID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Delete reservation
	_, err = tx.Exec("DELETE FROM temporary_reservations WHERE id = ?", req.ReservationID)
	if err != nil {
		log.Printf("Error deleting reservation %d: %v", req.ReservationID, err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		log.Printf("Error committing transaction: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Create response
	booking := models.Booking{
		ID:         int(bookingID),
		ClientName: req.ClientName,
		Email:      req.Email,
		Phone:      req.Phone,
		ServiceID:  req.ServiceID,
		Date:       req.Date,
		TimeSlot:   req.TimeSlot,
		CreatedAt:  time.Now(),
	}

	// Send response
	if err := json.NewEncoder(w).Encode(booking); err != nil {
		log.Printf("Error encoding booking response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Created booking %d for %s (%s) on %s at %s", 
		bookingID, req.ClientName, req.Email, req.Date, req.TimeSlot)
}

// validateBookingRequest validates the booking request fields
func validateBookingRequest(req models.BookingRequest) error {
	// Validate name
	if strings.TrimSpace(req.ClientName) == "" {
		return &ValidationError{Field: "client_name", Message: "Name is required"}
	}
	if len(strings.TrimSpace(req.ClientName)) < 2 {
		return &ValidationError{Field: "client_name", Message: "Name must be at least 2 characters"}
	}
	if !regexp.MustCompile(`^[a-zA-Z\s]+$`).MatchString(req.ClientName) {
		return &ValidationError{Field: "client_name", Message: "Name should contain only letters and spaces"}
	}

	// Validate email
	if strings.TrimSpace(req.Email) == "" {
		return &ValidationError{Field: "email", Message: "Email is required"}
	}
	if !regexp.MustCompile(`^[^\s@]+@[^\s@]+\.[^\s@]+$`).MatchString(req.Email) {
		return &ValidationError{Field: "email", Message: "Please enter a valid email"}
	}

	// Validate phone
	if strings.TrimSpace(req.Phone) == "" {
		return &ValidationError{Field: "phone", Message: "Phone is required"}
	}
	if !regexp.MustCompile(`^\+?[\d\s\-()]{8,}$`).MatchString(req.Phone) {
		return &ValidationError{Field: "phone", Message: "Please enter a valid phone number"}
	}

	// Validate other required fields
	if req.ReservationID <= 0 {
		return &ValidationError{Field: "reservation_id", Message: "Invalid reservation ID"}
	}
	if req.ServiceID <= 0 {
		return &ValidationError{Field: "service_id", Message: "Invalid service ID"}
	}
	if req.Date == "" {
		return &ValidationError{Field: "date", Message: "Date is required"}
	}
	if req.TimeSlot == "" {
		return &ValidationError{Field: "time_slot", Message: "Time slot is required"}
	}

	return nil
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

func (e *ValidationError) Error() string {
	return e.Message
}
