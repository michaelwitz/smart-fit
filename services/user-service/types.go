package main

import "time"

// User represents a user in the system
type User struct {
	ID            int        `json:"id"`
	FullName      string     `json:"fullName"`
	Email         string     `json:"email"`
	Password      string     `json:"password,omitempty"` // omitempty for security
	PhoneNumber   *string    `json:"phoneNumber"`
	Sex           *string    `json:"sex"`
	City          *string    `json:"city"`
	StateProvince *string    `json:"stateProvince"`
	PostalCode    *string    `json:"postalCode"`
	CountryCode   *string    `json:"countryCode"`
	Locale        *string    `json:"locale"`
	Timezone      *string    `json:"timezone"`
	UtcOffset     *int       `json:"utcOffset"`
	LastActive    *time.Time `json:"lastActive"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
}

// CreateUserRequest represents a request to create a new user
type CreateUserRequest struct {
	FullName      string  `json:"fullName"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	PhoneNumber   *string `json:"phoneNumber"`
	Sex           *string `json:"sex"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

// UpsertUserRequest represents a request to create or update a user
type UpsertUserRequest struct {
	FullName      string  `json:"fullName"`
	Email         string  `json:"email"`
	Password      *string `json:"password"` // Optional for updates, required for new users
	PhoneNumber   *string `json:"phoneNumber"`
	Sex           *string `json:"sex"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

// UpdateUserRequest represents a request to update an existing user
type UpdateUserRequest struct {
	FullName      *string `json:"fullName"`
	Email         *string `json:"email"`
	PhoneNumber   *string `json:"phoneNumber"`
	Sex           *string `json:"sex"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

// VerifyUserRequest represents a request to verify user credentials
type VerifyUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// VerifyUserResponse represents the response from user verification
type VerifyUserResponse struct {
	Valid bool  `json:"valid"`
	User  *User `json:"user,omitempty"`
}

// Password reset structs
type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ForgotPasswordResponse struct {
	Message string `json:"message"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"newPassword"`
}

type ResetPasswordResponse struct {
	Message string `json:"message"`
}


