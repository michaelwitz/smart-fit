package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	_ "github.com/lib/pq"
)

type Server struct {
	db *sql.DB
}

func main() {
	// Database connection
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbName := getEnv("DB_NAME", "smartfitgirl")
	dbUser := getEnv("DB_USER", "smartfit")
	dbPassword := getEnv("DB_PASSWORD", "smartfit123")

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Test database connection
	if err := db.Ping(); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}

	server := &Server{db: db}

	// Routes
	http.HandleFunc("/health", server.healthHandler)
	http.HandleFunc("/users", server.usersHandler)
	http.HandleFunc("/users/", server.userHandler)
	http.HandleFunc("/users/upsert", server.upsertHandler)
	http.HandleFunc("/users/verify", server.verifyHandler)

	port := getEnv("SERVICE_PORT", "8080")
	log.Printf("User service starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"status": "healthy", "service": "user-service"}`)
}

func (s *Server) usersHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getAllUsers(w, r)
	case http.MethodPost:
		s.createUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) userHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.getUserByID(w, r)
	case http.MethodPut:
		s.updateUser(w, r)
	case http.MethodDelete:
		s.deleteUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) upsertHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost, http.MethodPut:
		s.upsertUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) verifyHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.verifyUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
