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
		w.Write([]byte("Survey service is healthy"))
	}).Methods("GET")

	// Survey endpoints
	r.HandleFunc("/surveys", handleGetSurveys(db)).Methods("GET")
	r.HandleFunc("/surveys", handleCreateSurvey(db)).Methods("POST")
	r.HandleFunc("/surveys/{id}", handleGetSurvey(db)).Methods("GET")
	r.HandleFunc("/surveys/{id}", handleUpdateSurvey(db)).Methods("PUT")
	r.HandleFunc("/surveys/{id}", handleDeleteSurvey(db)).Methods("DELETE")

	// Goals endpoints
	r.HandleFunc("/goals", handleGetGoals(db)).Methods("GET")
	r.HandleFunc("/goals", handleCreateGoal(db)).Methods("POST")
	r.HandleFunc("/goals/{id}", handleGetGoal(db)).Methods("GET")
	r.HandleFunc("/goals/{id}", handleUpdateGoal(db)).Methods("PUT")
	r.HandleFunc("/goals/{id}", handleDeleteGoal(db)).Methods("DELETE")

	// Start server
	port := getEnv("PORT", "8080")
	log.Printf("Survey service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Handler functions (to be implemented)
func handleGetSurveys(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get surveys endpoint"))
	}
}

func handleCreateSurvey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create survey endpoint"))
	}
}

func handleGetSurvey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get survey by ID endpoint"))
	}
}

func handleUpdateSurvey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update survey endpoint"))
	}
}

func handleDeleteSurvey(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delete survey endpoint"))
	}
}

func handleGetGoals(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get goals endpoint"))
	}
}

func handleCreateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Create goal endpoint"))
	}
}

func handleGetGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Get goal by ID endpoint"))
	}
}

func handleUpdateGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Update goal endpoint"))
	}
}

func handleDeleteGoal(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Delete goal endpoint"))
	}
}
