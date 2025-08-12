package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"user-service/proto"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// getAllUsers retrieves all users via gRPC
func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := createContext()
	defer cancel()

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.GetAllUsers(ctx, &proto.GetAllUsersRequest{})
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get users: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.GetAllUsersResponse)
	if resp.Error != "" {
		http.Error(w, resp.Error, http.StatusInternalServerError)
		return
	}

	// Convert proto users to JSON users
	var users []User
	for _, protoUser := range resp.Users {
		if user := protoToUser(protoUser); user != nil {
			users = append(users, *user)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// getUserByID retrieves a specific user by ID via gRPC
func (s *Server) getUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 2 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := createContext()
	defer cancel()

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.GetUserByID(ctx, &proto.GetUserByIDRequest{Id: int32(id)})
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to get user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.GetUserByIDResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	user := protoToUser(resp.User)
	if user == nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// createUser creates a new user via gRPC
func (s *Server) createUser(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FullName == "" || req.Email == "" || req.Password == "" {
		http.Error(w, "full_name, email, and password are required", http.StatusBadRequest)
		return
	}

	// Hash password before sending to db-gateway
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	ctx, cancel := createContext()
	defer cancel()

	// Create proto request
	protoReq := &proto.CreateUserRequest{
		FullName:      req.FullName,
		Email:         req.Email,
		Password:      string(hashedPassword),
		PhoneNumber:   derefString(req.PhoneNumber),
		Sex:           derefString(req.Sex),
		City:          derefString(req.City),
		StateProvince: derefString(req.StateProvince),
		PostalCode:    derefString(req.PostalCode),
		CountryCode:   derefString(req.CountryCode),
		Locale:        derefString(req.Locale),
		Timezone:      derefString(req.Timezone),
		UtcOffset:     derefInt(req.UtcOffset),
	}

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.CreateUser(ctx, protoReq)
	})

	if err != nil {
		if strings.Contains(err.Error(), "duplicate") || strings.Contains(err.Error(), "already exists") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.CreateUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "duplicate") || strings.Contains(resp.Error, "already exists") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	user := protoToUser(resp.User)
	if user != nil {
		user.Password = "" // Clear password for security
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

// updateUser updates an existing user via gRPC
func (s *Server) updateUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 2 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Check if there are any fields to update
	if req.FullName == nil && req.Email == nil && req.PhoneNumber == nil &&
		req.Sex == nil && req.City == nil && req.StateProvince == nil &&
		req.PostalCode == nil && req.CountryCode == nil && req.Locale == nil &&
		req.Timezone == nil && req.UtcOffset == nil {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	ctx, cancel := createContext()
	defer cancel()

	// Create proto request
	protoReq := &proto.UpdateUserRequest{
		Id:            int32(id),
		FullName:      derefString(req.FullName),
		PhoneNumber:   derefString(req.PhoneNumber),
		Sex:           derefString(req.Sex),
		City:          derefString(req.City),
		StateProvince: derefString(req.StateProvince),
		PostalCode:    derefString(req.PostalCode),
		CountryCode:   derefString(req.CountryCode),
		Locale:        derefString(req.Locale),
		Timezone:      derefString(req.Timezone),
		UtcOffset:     derefInt(req.UtcOffset),
	}

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.UpdateUser(ctx, protoReq)
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		if strings.Contains(err.Error(), "duplicate") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.UpdateUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else if strings.Contains(resp.Error, "duplicate") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	user := protoToUser(resp.User)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// deleteUser deletes a user via gRPC
func (s *Server) deleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 2 {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	id, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	ctx, cancel := createContext()
	defer cancel()

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.DeleteUser(ctx, &proto.DeleteUserRequest{Id: int32(id)})
	})

	if err != nil {
		if st, ok := status.FromError(err); ok && st.Code() == codes.NotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.DeleteUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "not found") {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "%s"}`, resp.Message)
}

// upsertUser creates or updates a user via gRPC
func (s *Server) upsertUser(w http.ResponseWriter, r *http.Request) {
	var req UpsertUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.FullName == "" || req.Email == "" {
		http.Error(w, "full_name and email are required", http.StatusBadRequest)
		return
	}

	// Hash password if provided
	var hashedPassword string
	if req.Password != nil && *req.Password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}
		hashedPassword = string(hash)
	}

	ctx, cancel := createContext()
	defer cancel()

	// Create proto request
	protoReq := &proto.UpsertUserRequest{
		FullName:      req.FullName,
		Email:         req.Email,
		Password:      hashedPassword,
		PhoneNumber:   derefString(req.PhoneNumber),
		Sex:           derefString(req.Sex),
		City:          derefString(req.City),
		StateProvince: derefString(req.StateProvince),
		PostalCode:    derefString(req.PostalCode),
		CountryCode:   derefString(req.CountryCode),
		Locale:        derefString(req.Locale),
		Timezone:      derefString(req.Timezone),
		UtcOffset:     derefInt(req.UtcOffset),
	}

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.UpsertUser(ctx, protoReq)
	})

	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to upsert user: %v", err), http.StatusInternalServerError)
		return
	}

	resp := result.(*proto.UpsertUserResponse)
	if resp.Error != "" {
		if strings.Contains(resp.Error, "password is required") {
			http.Error(w, "password is required for new users", http.StatusBadRequest)
		} else {
			http.Error(w, resp.Error, http.StatusInternalServerError)
		}
		return
	}

	user := protoToUser(resp.User)
	if user != nil {
		user.Password = "" // Clear password for security
	}

	// For upsert, we'll use StatusOK since we can't tell if it was created or updated
	// TODO: Update proto to include a Created field in UpsertUserResponse
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}

// verifyUser verifies user credentials via gRPC
func (s *Server) verifyUser(w http.ResponseWriter, r *http.Request) {
	var req VerifyUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		http.Error(w, "email and password are required", http.StatusBadRequest)
		return
	}

	ctx, cancel := createContext()
	defer cancel()

	// Call gRPC service with circuit breaker
	result, err := s.callWithCircuitBreaker(ctx, func() (interface{}, error) {
		return s.userClient.VerifyUser(ctx, &proto.VerifyUserRequest{
			Email:    req.Email,
			Password: req.Password,
		})
	})

	if err != nil {
		// Don't reveal specific errors for security
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyUserResponse{Valid: false})
		return
	}

	resp := result.(*proto.VerifyUserResponse)
	
	// Build response
	verifyResp := VerifyUserResponse{
		Valid: resp.Valid,
	}
	
	if resp.Valid && resp.User != nil {
		user := protoToUser(resp.User)
		if user != nil {
			user.Password = "" // Clear password for security
			verifyResp.User = user
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(verifyResp)
}

// updateUserActivity updates the user's last active time
// This is now a fire-and-forget operation to the db-gateway
func (s *Server) updateUserActivity(userID int) {
	// This would need a new gRPC method in db-gateway to update last_active
	// For now, we'll skip this as it's not critical
	// TODO: Add UpdateUserActivity RPC to db-gateway
}
