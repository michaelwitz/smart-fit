package main

import (
	"context"
	"log"
	"strings"

	"user-service/proto"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCServer implements the UserService gRPC server
type GRPCServer struct {
	proto.UnimplementedUserServiceServer
	userClient proto.UserServiceClient
	cb         CallCircuitBreaker
}

// CallCircuitBreaker interface for circuit breaker functionality
type CallCircuitBreaker interface {
	Execute(func() (interface{}, error)) (interface{}, error)
}

// NewGRPCServer creates a new gRPC server instance
func NewGRPCServer(userClient proto.UserServiceClient, cb CallCircuitBreaker) *GRPCServer {
	return &GRPCServer{
		userClient: userClient,
		cb:         cb,
	}
}

// GetAllUsers implements the GetAllUsers RPC
func (s *GRPCServer) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.GetAllUsers(ctx, req)
	})

	if err != nil {
		log.Printf("Failed to get all users: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get users: %v", err)
	}

	return result.(*proto.GetAllUsersResponse), nil
}

// GetUserByID implements the GetUserByID RPC
func (s *GRPCServer) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.GetUserByID(ctx, req)
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("Failed to get user by ID: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	resp := result.(*proto.GetUserByIDResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, resp.Error)
	}

	return resp, nil
}

// CreateUser implements the CreateUser RPC
func (s *GRPCServer) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	// Validate required fields
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name, email, and password are required")
	}

	// Hash password before sending to db-gateway
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("Failed to hash password: %v", err)
		return nil, status.Error(codes.Internal, "failed to hash password")
	}

	// Update request with hashed password
	req.Password = string(hashedPassword)

	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.CreateUser(ctx, req)
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		log.Printf("Failed to create user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to create user: %v", err)
	}

	resp := result.(*proto.CreateUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "duplicate") || strings.Contains(resp.Error, "already exists") {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, resp.Error)
	}

	return resp, nil
}

// UpdateUser implements the UpdateUser RPC
func (s *GRPCServer) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.UpdateUser(ctx, req)
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if strings.Contains(err.Error(), "duplicate") {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		log.Printf("Failed to update user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

	resp := result.(*proto.UpdateUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		if strings.Contains(resp.Error, "duplicate") {
			return nil, status.Error(codes.AlreadyExists, "user with this email already exists")
		}
		return nil, status.Errorf(codes.Internal, resp.Error)
	}

	return resp, nil
}

// DeleteUser implements the DeleteUser RPC
func (s *GRPCServer) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.DeleteUser(ctx, req)
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		log.Printf("Failed to delete user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to delete user: %v", err)
	}

	resp := result.(*proto.DeleteUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			return nil, status.Error(codes.NotFound, "user not found")
		}
		return nil, status.Errorf(codes.Internal, resp.Error)
	}

	return resp, nil
}

// VerifyUser implements the VerifyUser RPC
func (s *GRPCServer) VerifyUser(ctx context.Context, req *proto.VerifyUserRequest) (*proto.VerifyUserResponse, error) {
	// Validate required fields
	if req.Email == "" || req.Password == "" {
		return &proto.VerifyUserResponse{
			Valid: false,
			Error: "email and password are required",
		}, nil
	}

	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.VerifyUser(ctx, req)
	})

	if err != nil {
		// Don't reveal specific errors for security
		log.Printf("VerifyUser error: %v", err)
		return &proto.VerifyUserResponse{
			Valid: false,
		}, nil
	}

	resp := result.(*proto.VerifyUserResponse)
	return resp, nil
}

// UpsertUser implements the UpsertUser RPC
func (s *GRPCServer) UpsertUser(ctx context.Context, req *proto.UpsertUserRequest) (*proto.UpsertUserResponse, error) {
	// Validate required fields
	if req.FullName == "" || req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "full_name and email are required")
	}

	// Hash password if provided
	if req.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password: %v", err)
			return nil, status.Error(codes.Internal, "failed to hash password")
		}
		req.Password = string(hashedPassword)
	}

	// Forward the request to db-gateway
	result, err := s.cb.Execute(func() (interface{}, error) {
		return s.userClient.UpsertUser(ctx, req)
	})

	if err != nil {
		log.Printf("Failed to upsert user: %v", err)
		return nil, status.Errorf(codes.Internal, "failed to upsert user: %v", err)
	}

	resp := result.(*proto.UpsertUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "password is required") {
			return nil, status.Error(codes.InvalidArgument, "password is required for new users")
		}
		return nil, status.Errorf(codes.Internal, resp.Error)
	}

	return resp, nil
}
