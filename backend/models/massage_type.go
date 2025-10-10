package models

// MassageType represents a massage service offered
type MassageType struct {
	ID       int     `json:"id" db:"id"`
	Name     string  `json:"name" db:"name"`
	Duration int     `json:"duration" db:"duration"` // minutes
	Price    float64 `json:"price" db:"price"`       // euros
}
