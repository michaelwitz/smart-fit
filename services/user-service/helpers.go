package main

import (
	"context"
	"time"

	"user-service/proto"
)

// Helper function to convert proto User to JSON User struct
func protoToUser(protoUser *proto.User) *User {
	if protoUser == nil {
		return nil
	}

	user := &User{
		ID:       int(protoUser.Id),
		FullName: protoUser.FullName,
		Email:    protoUser.Email,
	}

	// Convert optional fields
	if protoUser.PhoneNumber != "" {
		user.PhoneNumber = &protoUser.PhoneNumber
	}
	if protoUser.Sex != "" {
		user.Sex = &protoUser.Sex
	}
	if protoUser.City != "" {
		user.City = &protoUser.City
	}
	if protoUser.StateProvince != "" {
		user.StateProvince = &protoUser.StateProvince
	}
	if protoUser.PostalCode != "" {
		user.PostalCode = &protoUser.PostalCode
	}
	if protoUser.CountryCode != "" {
		user.CountryCode = &protoUser.CountryCode
	}
	if protoUser.Locale != "" {
		user.Locale = &protoUser.Locale
	}
	if protoUser.Timezone != "" {
		user.Timezone = &protoUser.Timezone
	}
	if protoUser.UtcOffset != 0 {
		utcOffset := int(protoUser.UtcOffset)
		user.UtcOffset = &utcOffset
	}

	// Convert timestamps
	if protoUser.LastActive != nil {
		t := protoUser.LastActive.AsTime()
		user.LastActive = &t
	}
	if protoUser.CreatedAt != nil {
		user.CreatedAt = protoUser.CreatedAt.AsTime()
	}
	if protoUser.UpdatedAt != nil {
		user.UpdatedAt = protoUser.UpdatedAt.AsTime()
	}

	return user
}

// Helper function to safely dereference string pointers
func derefString(s *string) string {
	if s != nil {
		return *s
	}
	return ""
}

// Helper function to safely dereference int pointers
func derefInt(i *int) int32 {
	if i != nil {
		return int32(*i)
	}
	return 0
}

// Circuit breaker wrapper for gRPC calls
func (s *Server) callWithCircuitBreaker(ctx context.Context, fn func() (interface{}, error)) (interface{}, error) {
	result, err := s.cb.Execute(func() (interface{}, error) {
		return fn()
	})
	return result, err
}

// Create a context with timeout for gRPC calls
func createContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 5*time.Second)
}
