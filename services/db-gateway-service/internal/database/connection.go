package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Config holds database configuration
type Config struct {
	Host            string
	Port            string
	User            string
	Password        string
	Database        string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

// Pool holds the database connection pool
type Pool struct {
	DB *sqlx.DB
}

// NewConfig creates database configuration from environment variables
func NewConfig() *Config {
	return &Config{
		Host:            getEnv("DB_HOST", "localhost"),
		Port:            getEnv("DB_PORT", "5432"),
		User:            getEnv("DB_USER", "smartfit"),
		Password:        getEnv("DB_PASSWORD", "smartfit123"),
		Database:        getEnv("DB_NAME", "smartfitgirl"),
		SSLMode:         getEnv("DB_SSLMODE", "disable"),
		MaxOpenConns:    25,
		MaxIdleConns:    5,
		ConnMaxLifetime: 5 * time.Minute,
	}
}

// NewPool creates a new database connection pool
func NewPool(config *Config) (*Pool, error) {
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.Database, config.SSLMode)

	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Configure connection pool
	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifetime)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Printf("Database connection pool established successfully")
	log.Printf("Max open connections: %d", config.MaxOpenConns)
	log.Printf("Max idle connections: %d", config.MaxIdleConns)
	log.Printf("Connection max lifetime: %v", config.ConnMaxLifetime)

	return &Pool{DB: db}, nil
}

// Close closes the database connection pool
func (p *Pool) Close() error {
	if p.DB != nil {
		return p.DB.Close()
	}
	return nil
}

// GetDB returns the underlying sqlx.DB instance
func (p *Pool) GetDB() *sqlx.DB {
	return p.DB
}

// HealthCheck checks if the database is healthy
func (p *Pool) HealthCheck() error {
	return p.DB.Ping()
}

// getEnv gets environment variable with fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
