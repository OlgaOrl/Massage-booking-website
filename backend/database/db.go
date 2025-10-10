package database

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
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

	// Create bookings table for future stories (not used in Story #1)
	bookingsTable := `
	CREATE TABLE IF NOT EXISTS bookings (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
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

// generateTimeSlots creates time slots for the next 30 days
func generateTimeSlots() error {
	// Create a new random source for Go 1.20+
	source := rand.NewSource(time.Now().UnixNano())
	rng := rand.New(source)

	// Get all massage type IDs
	rows, err := DB.Query("SELECT id FROM massage_types")
	if err != nil {
		return fmt.Errorf("failed to get massage types: %v", err)
	}
	defer rows.Close()

	var serviceIDs []int
	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			return fmt.Errorf("failed to scan service ID: %v", err)
		}
		serviceIDs = append(serviceIDs, id)
	}

	// Generate slots for next 30 days
	startDate := time.Now()
	for day := 0; day < 30; day++ {
		currentDate := startDate.AddDate(0, 0, day)
		dateStr := currentDate.Format("2006-01-02")

		// Generate time slots from 09:00 to 18:00 with 15-minute intervals
		for hour := 9; hour < 18; hour++ {
			for minute := 0; minute < 60; minute += 15 {
				timeStr := fmt.Sprintf("%02d:%02d", hour, minute)

				// Create slots for each service
				for _, serviceID := range serviceIDs {
					// Randomly make some slots unavailable (about 30% booked)
					available := rng.Float32() > 0.3

					_, err := DB.Exec("INSERT INTO time_slots (date, time, service_id, available) VALUES (?, ?, ?, ?)",
						dateStr, timeStr, serviceID, available)
					if err != nil {
						return fmt.Errorf("failed to insert time slot: %v", err)
					}
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

// GetTimeSlots retrieves time slots for a specific date and service
func GetTimeSlots(date string, serviceID int) ([]models.TimeSlot, error) {
	query := "SELECT id, date, time, service_id, available FROM time_slots WHERE date = ? AND service_id = ? ORDER BY time"
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

// CloseDB closes the database connection
func CloseDB() error {
	if DB != nil {
		return DB.Close()
	}
	return nil
}
