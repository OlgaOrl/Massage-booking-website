package models

import "time"

// Reservation represents a temporary slot reservation
type Reservation struct {
	ID        int       `json:"id" db:"id"`
	SlotID    int       `json:"slot_id" db:"slot_id"`
	ReservedAt time.Time `json:"reserved_at" db:"reserved_at"`
	ExpiresAt  time.Time `json:"expires_at" db:"expires_at"`
}

// ReservationRequest represents the request to create a reservation
type ReservationRequest struct {
	SlotID int `json:"slot_id"`
}

// ReservationResponse represents the response when creating a reservation
type ReservationResponse struct {
	ReservationID    int   `json:"reservation_id"`
	ExpiresAt        time.Time `json:"expires_at"`
	ExpiresInSeconds int   `json:"expires_in_seconds"`
}
