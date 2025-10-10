package models

// TimeSlot represents an available time slot for booking
type TimeSlot struct {
	ID        int    `json:"id" db:"id"`
	Date      string `json:"date" db:"date"`           // YYYY-MM-DD
	Time      string `json:"time" db:"time"`           // HH:MM
	ServiceID int    `json:"service_id" db:"service_id"`
	Available bool   `json:"available" db:"available"`
}
