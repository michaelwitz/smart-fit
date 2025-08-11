package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
	"time"
)

// TestConfig holds test configuration
type TestConfig struct {
	APIBaseURL     string
	UserServiceURL string
	TestUser       TestUser
}

// TestUser holds test user credentials
type TestUser struct {
	Email    string
	Password string
}

// NewTestConfig creates test configuration
func NewTestConfig() *TestConfig {
	return &TestConfig{
		APIBaseURL:     "http://localhost:8080",
		UserServiceURL: "http://localhost:8082",
		TestUser: TestUser{
			Email:    "sophia.woytowitz@gmail.com",
			Password: "password",
		},
	}
}

// TestHealthEndpoint tests the API gateway health check
func TestHealthEndpoint(t *testing.T) {
	config := NewTestConfig()
	
	resp, err := http.Get(config.APIBaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Check response content
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}

	if response["status"] != "healthy" {
		t.Errorf("Expected status 'healthy', got %v", response["status"])
	}

	if response["service"] != "api-service" {
		t.Errorf("Expected service 'api-service', got %v", response["service"])
	}
}

// TestLoginEndpoint tests the login endpoint with real test user
func TestLoginEndpoint(t *testing.T) {
	config := NewTestConfig()
	
	// Test with the real test user from seed data
	loginData := map[string]interface{}{
		"email":    config.TestUser.Email,
		"password": config.TestUser.Password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	resp, err := http.Post(config.APIBaseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", resp.StatusCode)
	}

	// Verify response structure
	var response map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	if response["token"] == "" {
		t.Error("Expected JWT token in login response")
	}

	if response["user"] == nil {
		t.Error("Expected user data in login response")
	}

	// Test that we can use the token to access protected endpoints
	token := response["token"].(string)
	req, err := http.NewRequest("GET", config.APIBaseURL+"/api/protected", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{}
	protectedResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to access protected endpoint: %v", err)
	}
	defer protectedResp.Body.Close()

	if protectedResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for protected endpoint with valid token, got %d", protectedResp.StatusCode)
	}
}

// TestProtectedEndpoint tests JWT authentication
func TestProtectedEndpoint(t *testing.T) {
	config := NewTestConfig()
	
	// Test accessing protected endpoint without token
	resp, err := http.Get(config.APIBaseURL + "/api/protected")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Should return unauthorized without JWT token
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for protected endpoint without token, got %d", resp.StatusCode)
	}
}

// TestErrorHandling tests API error responses
func TestErrorHandling(t *testing.T) {
	config := NewTestConfig()
	
	// Test invalid JSON
	resp, err := http.Post(config.APIBaseURL+"/auth/login", "application/json", bytes.NewBuffer([]byte("invalid json")))
	if err != nil {
		t.Fatalf("Failed to make invalid request: %v", err)
	}
	defer resp.Body.Close()

	// Should return bad request for invalid JSON
	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("Expected status 400 for invalid JSON, got %d", resp.StatusCode)
	}
}

// TestInvalidCredentials tests authentication with wrong password
func TestInvalidCredentials(t *testing.T) {
	config := NewTestConfig()
	
	// Test with wrong password for existing user
	loginData := map[string]interface{}{
		"email":    config.TestUser.Email,
		"password": "wrongpassword",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	resp, err := http.Post(config.APIBaseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make login request: %v", err)
	}
	defer resp.Body.Close()

	// Should return unauthorized for wrong password
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for wrong password, got %d", resp.StatusCode)
	}
}

// TestNonExistentUser tests authentication with non-existent user
func TestNonExistentUser(t *testing.T) {
	config := NewTestConfig()
	
	// Test with non-existent user
	loginData := map[string]interface{}{
		"email":    "nonexistent@example.com",
		"password": "password",
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	resp, err := http.Post(config.APIBaseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to make login request: %v", err)
	}
	defer resp.Body.Close()

	// Should return unauthorized for non-existent user
	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("Expected status 401 for non-existent user, got %d", resp.StatusCode)
	}
}

// TestCORSHeaders tests CORS middleware
func TestCORSHeaders(t *testing.T) {
	config := NewTestConfig()
	
	resp, err := http.Get(config.APIBaseURL + "/health")
	if err != nil {
		t.Fatalf("Failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Check CORS headers
	if resp.Header.Get("Access-Control-Allow-Origin") != "*" {
		t.Error("Expected CORS header Access-Control-Allow-Origin: *")
	}

	if resp.Header.Get("Access-Control-Allow-Methods") == "" {
		t.Error("Expected CORS header Access-Control-Allow-Methods")
	}
}

// TestEndToEndFlow tests the complete authentication flow
func TestEndToEndFlow(t *testing.T) {
	config := NewTestConfig()
	
	// Step 1: Login to get token
	loginData := map[string]interface{}{
		"email":    config.TestUser.Email,
		"password": config.TestUser.Password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		t.Fatalf("Failed to marshal login data: %v", err)
	}

	loginResp, err := http.Post(config.APIBaseURL+"/auth/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Failed to login: %v", err)
	}
	defer loginResp.Body.Close()

	if loginResp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed with status %d", loginResp.StatusCode)
	}

	var loginResponse map[string]interface{}
	if err := json.NewDecoder(loginResp.Body).Decode(&loginResponse); err != nil {
		t.Fatalf("Failed to decode login response: %v", err)
	}

	token := loginResponse["token"].(string)
	if token == "" {
		t.Fatal("No token received from login")
	}

	// Step 2: Use token to access protected endpoint
	req, err := http.NewRequest("GET", config.APIBaseURL+"/api/protected", nil)
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	req.Header.Set("Authorization", "Bearer "+token)

	client := &http.Client{Timeout: 10 * time.Second}
	protectedResp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to access protected endpoint: %v", err)
	}
	defer protectedResp.Body.Close()

	if protectedResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200 for protected endpoint, got %d", protectedResp.StatusCode)
	}

	// Step 3: Verify response contains user info
	var protectedResponse map[string]interface{}
	if err := json.NewDecoder(protectedResp.Body).Decode(&protectedResponse); err != nil {
		t.Fatalf("Failed to decode protected response: %v", err)
	}

	if protectedResponse["user_id"] == nil {
		t.Error("Expected user_id in protected response")
	}

	if protectedResponse["email"] != config.TestUser.Email {
		t.Errorf("Expected email %s, got %v", config.TestUser.Email, protectedResponse["email"])
	}
}
