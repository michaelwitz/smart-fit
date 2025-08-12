package main

import (
	"context"
	"log"
	"net"
	"os"
	"time"

	"user-service/proto"

	"github.com/sony/gobreaker"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

// CircuitBreakerAdapter adapts gobreaker.CircuitBreaker to our interface
type CircuitBreakerAdapter struct {
	cb *gobreaker.CircuitBreaker
}

func (c *CircuitBreakerAdapter) Execute(fn func() (interface{}, error)) (interface{}, error) {
	return c.cb.Execute(fn)
}

func main() {
	// Get service configuration
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8082" // Default gRPC port
	}

	// Get DB Gateway connection address
	dbGatewayAddr := os.Getenv("DB_GATEWAY_ADDR")
	if dbGatewayAddr == "" {
		dbGatewayAddr = "db-gateway-service:8086" // Default address
	}

	// Set up gRPC connection to db-gateway with keepalive
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialOpts := []grpc.DialOption{
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

	conn, err := grpc.DialContext(ctx, dbGatewayAddr, dialOpts...)
	if err != nil {
		log.Fatalf("Failed to connect to DB Gateway: %v", err)
	}
	defer conn.Close()

	// Create gRPC client to db-gateway
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
	cbAdapter := &CircuitBreakerAdapter{cb: cb}

	// Create gRPC server
	grpcServer := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    10 * time.Second, // Ping client if no activity for 10 seconds
			Timeout: 5 * time.Second,  // Wait 5 seconds for ping ack before considering connection dead
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Second, // Minimum time between client pings
			PermitWithoutStream: true,             // Allow pings even when there are no active streams
		}),
	)

	// Register the user service
	userService := NewGRPCServer(userClient, cbAdapter)
	proto.RegisterUserServiceServer(grpcServer, userService)

	// Register health service
	healthServer := health.NewServer()
	grpc_health_v1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("user-service", grpc_health_v1.HealthCheckResponse_SERVING)

	// Register reflection service for debugging
	reflection.Register(grpcServer)

	// Start listening
	listener, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", grpcPort, err)
	}

	log.Printf("User service gRPC server starting on port %s", grpcPort)
	log.Printf("Connected to DB Gateway at %s", dbGatewayAddr)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC: %v", err)
	}
}

