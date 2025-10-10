package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"strings"
	"time"

	"massage-booking/backend/models"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

// InitDB initializes the SQLite database connection and creates tables
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite", "./massage_booking.db")
	if err != nil {
		return fmt.Errorf("failed to open database: %v", err)
	}

	if err = DB.Ping(); err != nil {
		return fmt.Errorf("failed to ping database: %v", err)
	}

	if err = createTables(); err != nil {
		return fmt.Errorf("failed to create tables: %v", err)
	}

	if err = seedData(); err != nil {
		return fmt.Errorf("failed to seed data: %v", err)
	}

	log.Println("Database initialized successfully")
	return nil
}

// createTables creates the necessary database tables
func createTables() error {
	// Create massage_types table
	massageTypesTable := `
	CREATE TABLE IF NOT EXISTS massage_types (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		duration INTEGER NOT NULL,
		price REAL NOT NULL
	);`

	if _, err := DB.Exec(massageTypesTable); err != nil {
		return fmt.Errorf("failed to create massage_types table: %v", err)
	}

	// Create time_slots table
	timeSlotsTable := `
	CREATE TABLE IF NOT EXISTS time_slots (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		date TEXT NOT NULL,
		time TEXT NOT NULL,
		service_id INTEGER NOT NULL,
		available INTEGER NOT NULL DEFAULT 1,
		FOREIGN KEY (service_id) REFERENCES massage_types (id)
	);`

	if _, err := DB.Exec(timeSlotsTable); err != nil {
		return fmt.Errorf("failed to create time_slots table: %v", err)
	}

	// Create bookings table for Story #2 & #3
	bookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		reference TEXT UNIQUE NOT NULL,
		client_name TEXT NOT NULL,
		email TEXT NOT NULL,
		phone TEXT NOT NULL,
		service_id INTEGER NOT NULL,
		date TEXT NOT NULL,
		time_slot TEXT NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY (service_id) REFERENCES massage_types (id)
	);`

	if _, err := DB.Exec(bookingsTable); err != nil {
		return fmt.Errorf("failed to create bookings table: %v", err)
	}

	// Create temporary_reservations table for Story #2
	reservationsTable := `
	CREATE TABLE IF NOT EXISTS temporary_reservations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		slot_id INTEGER NOT NULL,
		reserved_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		expires_at DATETIME NOT NULL,
		FOREIGN KEY (slot_id) REFERENCES time_slots (id)
	);`

	if _, err := DB.Exec(reservationsTable); err != nil {
		return fmt.Errorf("failed to create temporary_reservations table: %v", err)
	}

	// Create index for cleanup queries
	indexQuery := `CREATE INDEX IF NOT EXISTS idx_expires_at ON temporary_reservations(expires_at);`
	if _, err := DB.Exec(indexQuery); err != nil {
		return fmt.Errorf("failed to create index on temporary_reservations: %v", err)
	}

	return nil
}

// seedData populates the database with initial sample data
func seedData() error {
	// Check if data already exists
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM massage_types").Scan(&count)
	if err != nil {
		return fmt.Errorf("failed to check existing data: %v", err)
	}

	if count > 0 {
		log.Println("Sample data already exists, skipping seed")
		return nil
	}

	// Insert massage types
	massageTypes := []models.MassageType{
		{Name: "Swedish Massage", Duration: 60, Price: 50.0},
		{Name: "Deep Tissue", Duration: 90, Price: 70.0},
		{Name: "Hot Stone", Duration: 60, Price: 65.0},
		{Name: "Sports Massage", Duration: 45, Price: 45.0},
	}

	for _, mt := range massageTypes {
		_, err := DB.Exec("INSERT INTO massage_types (name, duration, price) VALUES (?, ?, ?)",
			mt.Name, mt.Duration, mt.Price)
		if err != nil {
			return fmt.Errorf("failed to insert massage type: %v", err)
		}
	}

	// Generate time slots for the next 30 days
	if err := generateTimeSlots(); err != nil {
		return fmt.Errorf("failed to generate time slots: %v", err)
	}

	log.Println("Sample data seeded successfully")
	return nil
}

// generateTimeSlots creates time slots for the next 30 days based on service duration
func generateTimeSlots() error {
	// Create a new random source for Go 1.20+
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Get all massage types with their durations
	rows, err := DB.Query("SELECT id, duration FROM massage_types")
	if err != nil {
		return fmt.Errorf("failed to get massage types: %v", err)
	}
	defer rows.Close()

	var services []struct {
		ID       int
		Duration int
	}
	for rows.Next() {
		var service struct {
			ID       int
			Duration int
		}
		if err := rows.Scan(&service.ID, &service.Duration); err != nil {
			return fmt.Errorf("failed to scan service: %v", err)
		}
		services = append(services, service)
	}

	// Generate slots for next 30 days
	startDate := time.Now()
	for day := 0; day < 30; day++ {
		currentDate := startDate.AddDate(0, 0, day)
		dateStr := currentDate.Format("2006-01-02")

		// Generate time slots for each service based on its duration
		for _, service := range services {
			// Calculate how many slots can fit in a day based on service duration
			// Working hours: 09:00 to 18:00 (9 hours = 540 minutes)
			workingMinutes := 540
			slotsPerDay := workingMinutes / service.Duration

			// Generate slots with proper spacing
			for slot := 0; slot < slotsPerDay; slot++ {
				// Calculate start time for this slot
				startMinutes := 9*60 + slot*service.Duration // Start from 09:00
				if startMinutes >= 18*60 {                   // Don't go past 18:00
					break
				}

				hour := startMinutes / 60
				minute := startMinutes % 60
				timeStr := fmt.Sprintf("%02d:%02d", hour, minute)

				// Randomly make some slots unavailable (about 30% booked)
				available := rng.Float32() > 0.3

				_, err := DB.Exec("INSERT INTO time_slots (date, time, service_id, available) VALUES (?, ?, ?, ?)",
					dateStr, timeStr, service.ID, available)
				if err != nil {
					return fmt.Errorf("failed to insert time slot: %v", err)
				}
			}
		}
	}

	return nil
}

// GetMassageTypes retrieves all massage types from the database
func GetMassageTypes() ([]models.MassageType, error) {
	rows, err := DB.Query("SELECT id, name, duration, price FROM massage_types ORDER BY id")
	if err != nil {
		return nil, fmt.Errorf("failed to query massage types: %v", err)
	}
	defer rows.Close()

	var massageTypes []models.MassageType
	for rows.Next() {
		var mt models.MassageType
		if err := rows.Scan(&mt.ID, &mt.Name, &mt.Duration, &mt.Price); err != nil {
			return nil, fmt.Errorf("failed to scan massage type: %v", err)
		}
		massageTypes = append(massageTypes, mt)
	}

	return massageTypes, nil
}

// GetTimeSlots retrieves time slots for a specific date and service, excluding reserved slots
func GetTimeSlots(date string, serviceID int) ([]models.TimeSlot, error) {
	query := `
		SELECT ts.id, ts.date, ts.time, ts.service_id, ts.available
		FROM time_slots ts
		LEFT JOIN temporary_reservations tr ON ts.id = tr.slot_id AND tr.expires_at > datetime('now')
		WHERE ts.date = ? AND ts.service_id = ? AND tr.id IS NULL
		ORDER BY ts.time
	`
	rows, err := DB.Query(query, date, serviceID)
	if err != nil {
		return nil, fmt.Errorf("failed to query time slots: %v", err)
	}
	defer rows.Close()

	var timeSlots []models.TimeSlot
	for rows.Next() {
		var ts models.TimeSlot
		if err := rows.Scan(&ts.ID, &ts.Date, &ts.Time, &ts.ServiceID, &ts.Available); err != nil {
			return nil, fmt.Errorf("failed to scan time slot: %v", err)
		}
		timeSlots = append(timeSlots, ts)
	}

	return timeSlots, nil
}

// CleanupExpiredReservations removes expired reservations
func CleanupExpiredReservations() error {
	query := "DELETE FROM temporary_reservations WHERE expires_at < datetime('now')"
	result, err := DB.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired reservations: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err == nil && rowsAffected > 0 {
		log.Printf("Cleaned up %d expired reservations", rowsAffected)
	}

	return nil
}

// StartCleanupJob runs cleanup every minute to remove expired reservations
func StartCleanupJob() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		defer ticker.Stop()
		for range ticker.C {
			if err := CleanupExpiredReservations(); err != nil {
				log.Printf("Error during cleanup: %v", err)
			}
		}
	}()
	log.Println("Started cleanup job for expired reservations")
}

// CreateReservation creates a temporary reservation for a slot
func CreateReservation(slotID int) (int, time.Time, error) {
	// Check if slot exists and is available
	var available bool
	err := DB.QueryRow("SELECT available FROM time_slots WHERE id = ?", slotID).Scan(&available)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, time.Time{}, fmt.Errorf("slot not found")
		}
		return 0, time.Time{}, fmt.Errorf("failed to check slot availability: %v", err)
	}

	if !available {
		return 0, time.Time{}, fmt.Errorf("slot is not available")
	}

	// Check if slot is already reserved
	var count int
	err = DB.QueryRow("SELECT COUNT(*) FROM temporary_reservations WHERE slot_id = ? AND expires_at > datetime('now')", slotID).Scan(&count)
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to check existing reservations: %v", err)
	}

	if count > 0 {
		return 0, time.Time{}, fmt.Errorf("slot is already reserved")
	}

	// Create reservation with 10-minute expiration
	expiresAt := time.Now().Add(10 * time.Minute)
	result, err := DB.Exec("INSERT INTO temporary_reservations (slot_id, expires_at) VALUES (?, ?)",
		slotID, expiresAt.Format("2006-01-02 15:04:05"))
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to create reservation: %v", err)
	}

	reservationID, err := result.LastInsertId()
	if err != nil {
		return 0, time.Time{}, fmt.Errorf("failed to get reservation ID: %v", err)
	}

	return int(reservationID), expiresAt, nil
}

// DeleteReservation removes a temporary reservation
func DeleteReservation(reservationID int) error {
	result, err := DB.Exec("DELETE FROM temporary_reservations WHERE id = ?", reservationID)
	if err != nil {
		return fmt.Errorf("failed to delete reservation: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check affected rows: %v", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("reservation not found")
	}

	return nil
}

// IsSlotReserved checks if a slot is temporarily reserved
func IsSlotReserved(slotID int) (bool, error) {
	var count int
	err := DB.QueryRow("SELECT COUNT(*) FROM temporary_reservations WHERE slot_id = ? AND expires_at > datetime('now')", slotID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("failed to check slot reservation: %v", err)
	}

	return count > 0, nil
}

// GenerateBookingReference generates a unique booking reference
func GenerateBookingReference(date string) (string, error) {
	// Get count of bookings for this date
	var count int
	dateOnly := strings.Split(date, " ")[0] // Extract date part if datetime
	err := DB.QueryRow("SELECT COUNT(*) FROM bookings WHERE date = ?", dateOnly).Scan(&count)
	if err != nil {
		return "", fmt.Errorf("failed to get booking count: %v", err)
	}

	// Format: BK-YYYYMMDD-XXX
	dateStr := strings.ReplaceAll(dateOnly, "-", "")
	reference := fmt.Sprintf("BK-%s-%03d", dateStr, count+1)

	return reference, nil
}

// GetBookingByID retrieves a booking by ID with service details
func GetBookingByID(bookingID int) (*models.BookingDetail, error) {
	query := `
		SELECT b.id, b.reference, b.client_name, b.email, b.phone,
		       b.service_id, b.date, b.time_slot, b.created_at,
		       mt.name as service_name, mt.duration, mt.price
		FROM bookings b
		JOIN massage_types mt ON b.service_id = mt.id
		WHERE b.id = ?
	`

	var booking models.BookingDetail
	err := DB.QueryRow(query, bookingID).Scan(
		&booking.ID, &booking.Reference, &booking.ClientName, &booking.Email, &booking.Phone,
		&booking.ServiceID, &booking.Date, &booking.TimeSlot, &booking.CreatedAt,
		&booking.ServiceName, &booking.Duration, &booking.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("booking not found")
		}
		return nil, fmt.Errorf("failed to get booking: %v", err)
	}

	return &booking, nil
}

// CreateBookingWithReference creates a booking with generated reference
func CreateBookingWithReference(clientName, email, phone string, serviceID int, date, timeSlot string) (*models.BookingDetail, error) {
	// Generate reference
	reference, err := GenerateBookingReference(date)
	if err != nil {
		return nil, fmt.Errorf("failed to generate reference: %v", err)
	}

	// Insert booking
	result, err := DB.Exec(`
		INSERT INTO bookings (reference, client_name, email, phone, service_id, date, time_slot)
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`, reference, clientName, email, phone, serviceID, date, timeSlot)

	if err != nil {
		return nil, fmt.Errorf("failed to create booking: %v", err)
	}

	bookingID, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get booking ID: %v", err)
	}

	// Get the created booking with details
	return GetBookingByID(int(bookingID))
}

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
