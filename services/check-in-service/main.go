package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/jmoiron/sqlx"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sqlx.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Error connecting to database:", err)
	}
	defer db.Close()

	// Test database connection
	err = db.Ping()
	if err != nil {
		log.Fatal("Error pinging database:", err)
	}
	log.Println("Successfully connected to database")

	// Initialize router
	r := mux.NewRouter()

	// Health check endpoint
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Check-in service is healthy"))
	}).Methods("GET")

	// Check-in endpoints
	r.HandleFunc("/check-ins", handleGetCheckIns(db)).Methods("GET")
	r.HandleFunc("/check-ins", handleCreateCheckIn(db)).Methods("POST")
	r.HandleFunc("/check-ins/{id}", handleGetCheckIn(db)).Methods("GET")
	r.HandleFunc("/check-ins/{id}", handleUpdateCheckIn(db)).Methods("PUT")
	r.HandleFunc("/check-ins/{id}", handleDeleteCheckIn(db)).Methods("DELETE")

	// Start server
	port := os.Getenv("PORT")
	log.Printf("Check-in service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

// Handler functions (to be implemented)
func handleGetCheckIns(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get check-ins endpoint"))
	}
}

func handleCreateCheckIn(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create check-in endpoint"))
	}
}

func handleGetCheckIn(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get check-in by ID endpoint"))
	}
}

func handleUpdateCheckIn(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update check-in endpoint"))
	}
}

func handleDeleteCheckIn(db *sqlx.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delete check-in endpoint"))
	}
}
