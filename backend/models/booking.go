package models

import "time"

// Booking represents a confirmed booking
type Booking struct {
	ID          int       `json:"id" db:"id"`
	ClientName  string    `json:"client_name" db:"client_name"`
	Email       string    `json:"email" db:"email"`
	Phone       string    `json:"phone" db:"phone"`
	ServiceID   int       `json:"service_id" db:"service_id"`
	Date        string    `json:"date" db:"date"`
	TimeSlot    string    `json:"time_slot" db:"time_slot"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
}

// BookingRequest represents the request to create a booking
type BookingRequest struct {
	ReservationID int    `json:"reservation_id"`
	ClientName    string `json:"client_name"`
	Email         string `json:"email"`
	Phone         string `json:"phone"`
	ServiceID     int    `json:"service_id"`
	Date          string `json:"date"`
	TimeSlot      string `json:"time_slot"`
}
