package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "smartfit")
	dbPassword := getEnv("DB_PASSWORD", "smartfit123")
	dbName := getEnv("DB_NAME", "smartfitgirl")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", connStr)
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
	port := getEnv("PORT", "8080")
	log.Printf("Check-in service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Handler functions (to be implemented)
func handleGetCheckIns(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get check-ins endpoint"))
	}
}

func handleCreateCheckIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create check-in endpoint"))
	}
}

func handleGetCheckIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get check-in by ID endpoint"))
	}
}

func handleUpdateCheckIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update check-in endpoint"))
	}
}

func handleDeleteCheckIn(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delete check-in endpoint"))
	}
}
