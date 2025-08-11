package main

import (
	"log"
	"net"
	"os"

	"db-gateway-service/internal/database"
	"db-gateway-service/internal/services"
	"db-gateway-service/proto/checkin"
	"db-gateway-service/proto/user"
	"db-gateway-service/sql/check-in-service"
	"db-gateway-service/sql/user-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Environment variables - no hardcoded defaults for sensitive data
	port := getEnv("SERVICE_PORT", "8086") // Non-sensitive default OK
	dbHost := os.Getenv("DB_HOST")
	dbPort := getEnv("DB_PORT", "5432") // Non-sensitive default OK
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	// Validate required environment variables
	if dbHost == "" || dbUser == "" || dbPassword == "" || dbName == "" {
		log.Fatal("Required database environment variables not set: DB_HOST, DB_USER, DB_PASSWORD, DB_NAME")
	}

	// Initialize database connection pool
	dbConfig := &database.Config{
		Host:     dbHost,
		Port:     dbPort,
		User:     dbUser,
		Password: dbPassword,
		Database: dbName,
	}

	dbPool, err := database.NewPool(dbConfig)
	if err != nil {
		log.Fatalf("Failed to create database pool: %v", err)
	}
	defer dbPool.Close()

	// Initialize repositories
	userRepo := users.NewRepository(dbPool.GetDB())
	checkInRepo := checkins.NewRepository(dbPool.GetDB())

	// Create gRPC server
	grpcServer := grpc.NewServer()

	// Initialize and register services
	userService := services.NewUserService(userRepo)
	checkInService := services.NewCheckInService(checkInRepo)

	// Register services with gRPC server
	user.RegisterUserServiceServer(grpcServer, userService)
	checkin.RegisterCheckInServiceServer(grpcServer, checkInService)

	// Enable reflection for development
	reflection.Register(grpcServer)

	// Start gRPC server
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("DB Gateway Service starting on port %s", port)
	log.Printf("Database: %s:%s/%s", dbHost, dbPort, dbName)
	
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

// getEnv gets environment variable with fallback (only for non-sensitive data)
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
