package services

import (
	"context"
	"fmt"
	"log"

	"db-gateway-service/proto"
	users "db-gateway-service/sql/user-service"

	"google.golang.org/protobuf/types/known/timestamppb"
)

// UserService implements the gRPC UserService server
type UserService struct {
	proto.UnimplementedUserServiceServer
	repo *users.Repository
}

// NewUserService creates a new UserService instance
func NewUserService(repo *users.Repository) *UserService {
	return &UserService{
		repo: repo,
	}
}

// CreateUser creates a new user
func (s *UserService) CreateUser(ctx context.Context, req *proto.CreateUserRequest) (*proto.CreateUserResponse, error) {
	log.Printf("CreateUser called with email: %s", req.Email)
	
	// Convert request to repository model
	dbUser := &users.User{
		FullName:      req.FullName,
		Email:         req.Email,
		Password:      req.Password,
		PhoneNumber:   stringToPtr(req.PhoneNumber),
		Sex:           stringToPtr(req.Sex),
		City:          stringToPtr(req.City),
		StateProvince: stringToPtr(req.StateProvince),
		PostalCode:    stringToPtr(req.PostalCode),
		CountryCode:   stringToPtr(req.CountryCode),
		Locale:        stringToPtr(req.Locale),
		Timezone:      stringToPtr(req.Timezone),
		UtcOffset:     intToPtr(int(req.UtcOffset)),
	}

	// Create user in database
	if err := s.repo.CreateUser(dbUser); err != nil {
		log.Printf("Failed to create user: %v", err)
		return &proto.CreateUserResponse{
			Error: fmt.Sprintf("Failed to create user: %v", err),
		}, nil
	}

	// Convert to protobuf response
	return &proto.CreateUserResponse{
		User: convertToProtoUser(dbUser),
	}, nil
}

// GetUserByID retrieves a user by ID
func (s *UserService) GetUserByID(ctx context.Context, req *proto.GetUserByIDRequest) (*proto.GetUserByIDResponse, error) {
	log.Printf("GetUserByID called with ID: %d", req.Id)
	
	dbUser, err := s.repo.GetUserByID(int(req.Id))
	if err != nil {
		log.Printf("Failed to get user by ID: %v", err)
		return &proto.GetUserByIDResponse{
			Error: fmt.Sprintf("Failed to get user: %v", err),
		}, nil
	}

	return &proto.GetUserByIDResponse{
		User: convertToProtoUser(dbUser),
	}, nil
}

// GetAllUsers retrieves all users
func (s *UserService) GetAllUsers(ctx context.Context, req *proto.GetAllUsersRequest) (*proto.GetAllUsersResponse, error) {
	log.Printf("GetAllUsers called")
	
	dbUsers, err := s.repo.GetAllUsers()
	if err != nil {
		log.Printf("Failed to get all users: %v", err)
		return &proto.GetAllUsersResponse{
			Error: fmt.Sprintf("Failed to get users: %v", err),
		}, nil
	}

	// Convert to protobuf users
	protoUsers := make([]*proto.User, len(dbUsers))
	for i, dbUser := range dbUsers {
		u := dbUser // Create a copy to get the address
		protoUsers[i] = convertToProtoUser(&u)
	}

	return &proto.GetAllUsersResponse{
		Users: protoUsers,
	}, nil
}

// UpdateUser updates an existing user
func (s *UserService) UpdateUser(ctx context.Context, req *proto.UpdateUserRequest) (*proto.UpdateUserResponse, error) {
	log.Printf("UpdateUser called for ID: %d", req.Id)
	
	// Convert request to repository model
	dbUser := &users.User{
		ID:            int(req.Id),
		FullName:      req.FullName,
		Password:      req.Password,
		PhoneNumber:   stringToPtr(req.PhoneNumber),
		Sex:           stringToPtr(req.Sex),
		City:          stringToPtr(req.City),
		StateProvince: stringToPtr(req.StateProvince),
		PostalCode:    stringToPtr(req.PostalCode),
		CountryCode:   stringToPtr(req.CountryCode),
		Locale:        stringToPtr(req.Locale),
		Timezone:      stringToPtr(req.Timezone),
		UtcOffset:     intToPtr(int(req.UtcOffset)),
	}

	// Update user in database
	if err := s.repo.UpdateUser(dbUser); err != nil {
		log.Printf("Failed to update user: %v", err)
		return &proto.UpdateUserResponse{
			Error: fmt.Sprintf("Failed to update user: %v", err),
		}, nil
	}

	return &proto.UpdateUserResponse{
		User: convertToProtoUser(dbUser),
	}, nil
}

// DeleteUser deletes a user by ID
func (s *UserService) DeleteUser(ctx context.Context, req *proto.DeleteUserRequest) (*proto.DeleteUserResponse, error) {
	log.Printf("DeleteUser called for ID: %d", req.Id)
	
	if err := s.repo.DeleteUser(int(req.Id)); err != nil {
		log.Printf("Failed to delete user: %v", err)
		return &proto.DeleteUserResponse{
			Error: fmt.Sprintf("Failed to delete user: %v", err),
		}, nil
	}

	return &proto.DeleteUserResponse{
		Message: fmt.Sprintf("User with ID %d deleted successfully", req.Id),
	}, nil
}

// VerifyUser verifies user credentials
func (s *UserService) VerifyUser(ctx context.Context, req *proto.VerifyUserRequest) (*proto.VerifyUserResponse, error) {
	log.Printf("VerifyUser called for email: %s", req.Email)
	
	dbUser, err := s.repo.VerifyUser(req.Email, req.Password)
	if err != nil {
		log.Printf("Failed to verify user: %v", err)
		return &proto.VerifyUserResponse{
			Valid: false,
			Error: fmt.Sprintf("Invalid credentials"),
		}, nil
	}

	// Note: The repository doesn't actually verify the password,
	// so we need to do it here. In a real implementation, you'd use
	// a proper password hashing library like bcrypt
	if dbUser.Password != req.Password {
		return &proto.VerifyUserResponse{
			Valid: false,
			Error: "Invalid credentials",
		}, nil
	}

	return &proto.VerifyUserResponse{
		Valid: true,
		User:  convertToProtoUser(dbUser),
	}, nil
}

// UpsertUser creates or updates a user
func (s *UserService) UpsertUser(ctx context.Context, req *proto.UpsertUserRequest) (*proto.UpsertUserResponse, error) {
	log.Printf("UpsertUser called with email: %s", req.Email)
	
	// Convert request to repository model
	dbUser := &users.User{
		FullName:      req.FullName,
		Email:         req.Email,
		Password:      req.Password,
		PhoneNumber:   stringToPtr(req.PhoneNumber),
		Sex:           stringToPtr(req.Sex),
		City:          stringToPtr(req.City),
		StateProvince: stringToPtr(req.StateProvince),
		PostalCode:    stringToPtr(req.PostalCode),
		CountryCode:   stringToPtr(req.CountryCode),
		Locale:        stringToPtr(req.Locale),
		Timezone:      stringToPtr(req.Timezone),
		UtcOffset:     intToPtr(int(req.UtcOffset)),
	}

	// Upsert user in database
	if err := s.repo.UpsertUser(dbUser); err != nil {
		log.Printf("Failed to upsert user: %v", err)
		return &proto.UpsertUserResponse{
			Error: fmt.Sprintf("Failed to upsert user: %v", err),
		}, nil
	}

	return &proto.UpsertUserResponse{
		User: convertToProtoUser(dbUser),
	}, nil
}

// Helper function to convert database user to protobuf user
func convertToProtoUser(dbUser *users.User) *proto.User {
	protoUser := &proto.User{
		Id:            int32(dbUser.ID),
		FullName:      dbUser.FullName,
		Email:         dbUser.Email,
		PhoneNumber:   ptrToString(dbUser.PhoneNumber),
		Sex:           ptrToString(dbUser.Sex),
		City:          ptrToString(dbUser.City),
		StateProvince: ptrToString(dbUser.StateProvince),
		PostalCode:    ptrToString(dbUser.PostalCode),
		CountryCode:   ptrToString(dbUser.CountryCode),
		Locale:        ptrToString(dbUser.Locale),
		Timezone:      ptrToString(dbUser.Timezone),
		UtcOffset:     ptrToInt32(dbUser.UtcOffset),
		CreatedAt:     timestamppb.New(dbUser.CreatedAt),
		UpdatedAt:     timestamppb.New(dbUser.UpdatedAt),
	}

	if dbUser.LastActive != nil {
		protoUser.LastActive = timestamppb.New(*dbUser.LastActive)
	}

	return protoUser
}

// Helper functions for pointer conversions
func stringToPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func intToPtr(i int) *int {
	if i == 0 {
		return nil
	}
	return &i
}

func ptrToString(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func ptrToInt32(i *int) int32 {
	if i == nil {
		return 0
	}
	return int32(*i)
}
