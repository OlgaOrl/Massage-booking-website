package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"massage-booking/backend/database"
	"massage-booking/backend/handlers"
)

func main() {
	// Initialize database
	if err := database.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Start cleanup job for expired reservations
	database.StartCleanupJob()

	// Set up graceful shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		log.Println("Shutting down gracefully...")
		if err := database.CloseDB(); err != nil {
			log.Printf("Error closing database: %v", err)
		}
		os.Exit(0)
	}()

	// Set up routes
	setupRoutes()

	// Start server
	port := ":8080"
	log.Printf("Server starting on port %s", port)
	log.Printf("Static files served from: /static")
	log.Printf("API endpoints available at: /api")

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// setupRoutes configures all HTTP routes
func setupRoutes() {
	// API routes
	http.HandleFunc("/api/massage-types", handlers.GetMassageTypesHandler)
	http.HandleFunc("/api/slots", handlers.GetSlotsHandler)

	// Story #2 routes
	http.HandleFunc("/api/reservations", handlers.CreateReservation)
	http.HandleFunc("/api/reservations/", handlers.DeleteReservation)
	http.HandleFunc("/api/bookings", handlers.CreateBooking)

	// Static file server for frontend
	fs := http.FileServer(http.Dir("./backend/static/"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve index.html at root
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.ServeFile(w, r, "./backend/static/index.html")
		} else {
			http.NotFound(w, r)
		}
	})

	log.Println("Routes configured successfully")
}
