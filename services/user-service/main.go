package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"user-service/proto"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/keepalive"
)

type Server struct {
	userClient proto.UserServiceClient
	cb         *gobreaker.CircuitBreaker
}

func main() {
	// Get DB Gateway connection address
	dbGatewayAddr := os.Getenv("DB_GATEWAY_ADDR")
	if dbGatewayAddr == "" {
		dbGatewayAddr = "db-gateway-service:8086" // Default address
	}

	// Set up gRPC connection with connection pooling and keepalive
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:                10 * time.Second, // Send pings every 10 seconds if no activity
			Timeout:             3 * time.Second,  // Wait 3 seconds for ping acknowledgement
			PermitWithoutStream: true,              // Send pings even without active streams
		}),
		grpc.WithDefaultCallOptions(
			grpc.MaxCallRecvMsgSize(10*1024*1024), // 10MB max message size
			grpc.MaxCallSendMsgSize(10*1024*1024),
		),
	}

	conn, err := grpc.DialContext(ctx, dbGatewayAddr, opts...)
	if err != nil {
		log.Fatalf("Failed to connect to DB Gateway: %v", err)
	}
	defer conn.Close()

	// Create gRPC client
	userClient := proto.NewUserServiceClient(conn)

	// Set up circuit breaker for resilience
	settings := gobreaker.Settings{
		Name:        "DBGateway",
		MaxRequests: 3,                // Number of requests allowed to pass through when half-open
		Interval:    10 * time.Second, // Time window for failure rate calculation
		Timeout:     30 * time.Second, // Time before attempting to recover
		ReadyToTrip: func(counts gobreaker.Counts) bool {
			failureRatio := float64(counts.TotalFailures) / float64(counts.Requests)
			return counts.Requests >= 3 && failureRatio >= 0.6 // Trip if 60% failure rate
		},
		OnStateChange: func(name string, from gobreaker.State, to gobreaker.State) {
			log.Printf("Circuit breaker %s changed from %v to %v", name, from, to)
		},
	}

	cb := gobreaker.NewCircuitBreaker(settings)

	server := &Server{
		userClient: userClient,
		cb:         cb,
	}

	// Routes
	http.HandleFunc("/health", server.healthHandler)
	http.HandleFunc("/users", server.usersHandler)
	http.HandleFunc("/users/upsert", server.upsertHandler)
	http.HandleFunc("/users/verify", server.verifyHandler)
	http.HandleFunc("/auth/forgot-password", server.forgotPasswordHandler)
	http.HandleFunc("/auth/reset-password", server.resetPasswordHandler)
	http.HandleFunc("/users/", server.userHandler) // Handles /users/{id} CRUD operations

	port := os.Getenv("SERVICE_PORT")
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

func (s *Server) forgotPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.forgotPassword(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (s *Server) resetPasswordHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		s.resetPassword(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

