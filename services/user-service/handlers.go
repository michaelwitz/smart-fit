package main

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID            int        `json:"id" db:"id"`
	FullName      string     `json:"fullName" db:"full_name"`
	Email         string     `json:"email" db:"email"`
	Password      string     `json:"password,omitempty" db:"password"` // omitempty for security
	PhoneNumber   *string    `json:"phoneNumber" db:"phone_number"`
	IdentifyAs    *string    `json:"identifyAs" db:"identify_as"`
	City          *string    `json:"city" db:"city"`
	StateProvince *string    `json:"stateProvince" db:"state_province"`
	PostalCode    *string    `json:"postalCode" db:"postal_code"`
	CountryCode   *string    `json:"countryCode" db:"country_code"`
	Locale        *string    `json:"locale" db:"locale"`
	Timezone      *string    `json:"timezone" db:"timezone"`
	UtcOffset     *int       `json:"utcOffset" db:"utc_offset"`
	LastActive    *time.Time `json:"lastActive" db:"last_active"`
	CreatedAt     time.Time  `json:"createdAt" db:"created_at"`
	UpdatedAt     time.Time  `json:"updatedAt" db:"updated_at"`
}

type CreateUserRequest struct {
	FullName      string  `json:"fullName"`
	Email         string  `json:"email"`
	Password      string  `json:"password"`
	PhoneNumber   *string `json:"phoneNumber"`
	IdentifyAs    *string `json:"identifyAs"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

type UpsertUserRequest struct {
	FullName      string  `json:"fullName"`
	Email         string  `json:"email"`
	Password      *string `json:"password"` // Optional for updates, required for new users
	PhoneNumber   *string `json:"phoneNumber"`
	IdentifyAs    *string `json:"identifyAs"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

type UpdateUserRequest struct {
	FullName      *string `json:"fullName"`
	Email         *string `json:"email"`
	PhoneNumber   *string `json:"phoneNumber"`
	IdentifyAs    *string `json:"identifyAs"`
	City          *string `json:"city"`
	StateProvince *string `json:"stateProvince"`
	PostalCode    *string `json:"postalCode"`
	CountryCode   *string `json:"countryCode"`
	Locale        *string `json:"locale"`
	Timezone      *string `json:"timezone"`
	UtcOffset     *int    `json:"utcOffset"`
}

type VerifyUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type VerifyUserResponse struct {
	Valid bool `json:"valid"`
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

// Survey-related structs
type Goal struct {
	ID          int       `json:"id" db:"id"`
	Category    string    `json:"category" db:"category"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description" db:"description"`
	CreatedAt   time.Time `json:"createdAt" db:"created_at"`
	UpdatedAt   time.Time `json:"updatedAt" db:"updated_at"`
}

// Goal with selection state for UI
type GoalWithSelection struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Selected    bool   `json:"selected"`
}

// Goals grouped by category for UI
type GoalsByCategory struct {
	Weight     []GoalWithSelection `json:"weight"`
	Appearance []GoalWithSelection `json:"appearance"`
	Strength   []GoalWithSelection `json:"strength"`
	Endurance  []GoalWithSelection `json:"endurance"`
}

type Survey struct {
	ID            int       `json:"id" db:"id"`
	UserID        int       `json:"userId" db:"user_id"`
	CurrentWeight float64   `json:"currentWeight" db:"current_weight"`
	TargetWeight  float64   `json:"targetWeight" db:"target_weight"`
	ActivityLevel int       `json:"activityLevel" db:"activity_level"`
	CreatedAt     time.Time `json:"createdAt" db:"created_at"`
}

// Optimized survey response with goals grouped by category
type SurveyWithGoals struct {
	Survey
	Goals GoalsByCategory `json:"goals"`
}

// Response for getting all available goals (for survey creation form)
type AllGoalsResponse struct {
	Goals GoalsByCategory `json:"goals"`
}

type CreateSurveyRequest struct {
	CurrentWeight float64 `json:"currentWeight"`
	TargetWeight  float64 `json:"targetWeight"`
	ActivityLevel int     `json:"activityLevel"`
	GoalIDs       []int   `json:"goalIds"`
}

func (s *Server) forgotPassword(w http.ResponseWriter, r *http.Request) {
	var req ForgotPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate email
	if req.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Generate a password reset token
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		http.Error(w, "Failed to generate reset token", http.StatusInternalServerError)
		return
	}
	resetToken := hex.EncodeToString(b)

	// Store the token in the database with an expiration date
	expiration := time.Now().Add(1 * time.Hour)
	query := "INSERT INTO PASSWORD_RESETS (email, token, expires_at) VALUES ($1, $2, $3)"
	_, err := s.db.Exec(query, req.Email, resetToken, expiration)
	if err != nil {
		http.Error(w, "Error saving reset token", http.StatusInternalServerError)
		return
	}

	// Send the password reset email
	from := mail.NewEmail("Example App", "no-reply@example.com")
	subject := "Password Reset Request"
	to := mail.NewEmail("User", req.Email)
	plainTextContent := "Please use this token to reset your password: " + resetToken
	htmlContent := "<strong>Please use this token to reset your password: </strong>" + resetToken
	message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
	client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
	response, err := client.Send(message)
	if err != nil || response.StatusCode >= 400 {
		http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	resp := ForgotPasswordResponse{Message: "Password reset token sent to email."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) resetPassword(w http.ResponseWriter, r *http.Request) {
	var req ResetPasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate input
	if req.Token == "" || req.NewPassword == "" {
		http.Error(w, "Token and new password are required", http.StatusBadRequest)
		return
	}

	// Verify the token and get the associated email
	var email string
	query := "SELECT email FROM PASSWORD_RESETS WHERE token = $1 AND expires_at > $2"
	if err := s.db.QueryRow(query, req.Token, time.Now()).Scan(&email); err != nil {
		http.Error(w, "Invalid or expired token", http.StatusBadRequest)
		return
	}

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
		return
	}

	// Update the user's password in the database
	updateQuery := "UPDATE USERS SET password = $2 WHERE email = $1"
	_, err = s.db.Exec(updateQuery, email, string(hashedPassword))
	if err != nil {
		http.Error(w, "Error updating password", http.StatusInternalServerError)
		return
	}

	// Respond with success message
	resp := ResetPasswordResponse{Message: "Password reset successful."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) getAllUsers(w http.ResponseWriter, r *http.Request) {
	query := `
		SELECT id, full_name, email, phone_number, identify_as, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		ORDER BY created_at DESC`

	var users []User
	err := s.db.Select(&users, query)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

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

	query := `
		SELECT id, full_name, email, phone_number, identify_as, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		WHERE id = $1`

	var user User
	err = s.db.Get(&user, query, id)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

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

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		return
	}

	query := `
		INSERT INTO USERS (full_name, email, password, phone_number, identify_as, 
		                  city, state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
		RETURNING id, created_at, updated_at`

	var user User
	err = s.db.QueryRow(
		query, req.FullName, req.Email, string(hashedPassword),
		req.PhoneNumber, req.IdentifyAs, req.City, req.StateProvince,
		req.PostalCode, req.CountryCode, req.Locale, req.Timezone, req.UtcOffset,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// Fill in the response user object (excluding password)
	user.FullName = req.FullName
	user.Email = req.Email
	user.PhoneNumber = req.PhoneNumber
	user.IdentifyAs = req.IdentifyAs
	user.City = req.City
	user.StateProvince = req.StateProvince
	user.PostalCode = req.PostalCode
	user.CountryCode = req.CountryCode
	user.Locale = req.Locale
	user.Timezone = req.Timezone
	user.UtcOffset = req.UtcOffset

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

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

	// Check if user exists
	var existingUserID int
	var existingPassword string
	err := s.db.QueryRow("SELECT id, password FROM USERS WHERE email = $1", req.Email).Scan(&existingUserID, &existingPassword)
	
	if err == sql.ErrNoRows {
		// User doesn't exist - INSERT (password is required for new users)
		if req.Password == nil || *req.Password == "" {
			http.Error(w, "password is required for new users", http.StatusBadRequest)
			return
		}

		// Hash password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "Failed to hash password", http.StatusInternalServerError)
			return
		}

		query := `
			INSERT INTO USERS (full_name, email, password, phone_number, identify_as, 
			                  city, state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
			RETURNING id, created_at, updated_at`

		var user User
		err = s.db.QueryRow(
			query, req.FullName, req.Email, string(hashedPassword),
			req.PhoneNumber, req.IdentifyAs, req.City, req.StateProvince,
			req.PostalCode, req.CountryCode, req.Locale, req.Timezone, req.UtcOffset,
		).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		// Fill in the response user object (excluding password)
		user.FullName = req.FullName
		user.Email = req.Email
		user.PhoneNumber = req.PhoneNumber
		user.IdentifyAs = req.IdentifyAs
		user.City = req.City
		user.StateProvince = req.StateProvince
		user.PostalCode = req.PostalCode
		user.CountryCode = req.CountryCode
		user.Locale = req.Locale
		user.Timezone = req.Timezone
		user.UtcOffset = req.UtcOffset

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)

	} else if err != nil {
		// Database error
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	} else {
		// User exists - UPDATE (password update is optional)
		passwordToUse := existingPassword // Keep existing password by default
		
		// If new password is provided, hash it
		if req.Password != nil && *req.Password != "" {
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(*req.Password), bcrypt.DefaultCost)
			if err != nil {
				http.Error(w, "Failed to hash password", http.StatusInternalServerError)
				return
			}
			passwordToUse = string(hashedPassword)
		}

		query := `
			UPDATE USERS 
			SET full_name = $1, password = $2, phone_number = $3, identify_as = $4, 
			    city = $5, state_province = $6, postal_code = $7, country_code = $8, locale = $9, timezone = $10, 
			    utc_offset = $11, updated_at = CURRENT_TIMESTAMP
			WHERE id = $12
			RETURNING id, full_name, email, phone_number, identify_as, city, 
			          state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at`

		var user User
		err = s.db.QueryRow(
			query, req.FullName, passwordToUse, req.PhoneNumber, req.IdentifyAs,
			req.City, req.StateProvince, req.PostalCode, req.CountryCode, req.Locale, req.Timezone, 
			req.UtcOffset, existingUserID,
		).Scan(
			&user.ID, &user.FullName, &user.Email, &user.PhoneNumber,
			&user.IdentifyAs, &user.City, &user.StateProvince, &user.PostalCode,
			&user.CountryCode, &user.Locale, &user.Timezone, &user.UtcOffset, 
			&user.CreatedAt, &user.UpdatedAt,
		)

		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(user)
	}
}

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

	// Build dynamic update query
	setParts := []string{}
	args := []interface{}{}
	argIndex := 1

	if req.FullName != nil {
		setParts = append(setParts, fmt.Sprintf("full_name = $%d", argIndex))
		args = append(args, *req.FullName)
		argIndex++
	}
	if req.Email != nil {
		setParts = append(setParts, fmt.Sprintf("email = $%d", argIndex))
		args = append(args, *req.Email)
		argIndex++
	}
	if req.PhoneNumber != nil {
		setParts = append(setParts, fmt.Sprintf("phone_number = $%d", argIndex))
		args = append(args, *req.PhoneNumber)
		argIndex++
	}
	if req.IdentifyAs != nil {
		setParts = append(setParts, fmt.Sprintf("identify_as = $%d", argIndex))
		args = append(args, *req.IdentifyAs)
		argIndex++
	}
	if req.City != nil {
		setParts = append(setParts, fmt.Sprintf("city = $%d", argIndex))
		args = append(args, *req.City)
		argIndex++
	}
	if req.StateProvince != nil {
		setParts = append(setParts, fmt.Sprintf("state_province = $%d", argIndex))
		args = append(args, *req.StateProvince)
		argIndex++
	}
	if req.PostalCode != nil {
		setParts = append(setParts, fmt.Sprintf("postal_code = $%d", argIndex))
		args = append(args, *req.PostalCode)
		argIndex++
	}
	if req.CountryCode != nil {
		setParts = append(setParts, fmt.Sprintf("country_code = $%d", argIndex))
		args = append(args, *req.CountryCode)
		argIndex++
	}
	if req.Locale != nil {
		setParts = append(setParts, fmt.Sprintf("locale = $%d", argIndex))
		args = append(args, *req.Locale)
		argIndex++
	}
	if req.Timezone != nil {
		setParts = append(setParts, fmt.Sprintf("timezone = $%d", argIndex))
		args = append(args, *req.Timezone)
		argIndex++
	}
	if req.UtcOffset != nil {
		setParts = append(setParts, fmt.Sprintf("utc_offset = $%d", argIndex))
		args = append(args, *req.UtcOffset)
		argIndex++
	}

	if len(setParts) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Always update the updated_at field
	setParts = append(setParts, "updated_at = CURRENT_TIMESTAMP")
	args = append(args, id)

	query := fmt.Sprintf(`
		UPDATE USERS 
		SET %s 
		WHERE id = $%d
		RETURNING id, full_name, email, phone_number, identify_as, city, 
		          state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at`,
		strings.Join(setParts, ", "), argIndex)

	var user User
	err = s.db.QueryRow(query, args...).Scan(
		&user.ID, &user.FullName, &user.Email, &user.PhoneNumber,
		&user.IdentifyAs, &user.City, &user.StateProvince, &user.PostalCode,
		&user.CountryCode, &user.Locale, &user.Timezone, &user.UtcOffset, 
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	if err != nil {
		if strings.Contains(err.Error(), "duplicate key") {
			http.Error(w, "User with this email already exists", http.StatusConflict)
			return
		}
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

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

	query := `DELETE FROM USERS WHERE id = $1`
	result, err := s.db.Exec(query, id)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, `{"message": "User deleted successfully"}`)
}

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

	// Get user from database
	query := `
		SELECT id, full_name, email, password, phone_number, identify_as, city, 
		       state_province, postal_code, country_code, locale, timezone, utc_offset, created_at, updated_at 
		FROM USERS 
		WHERE email = $1`

	var user User
	err := s.db.Get(&user, query, req.Email)

	if err == sql.ErrNoRows {
		// User not found - return invalid without revealing this
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyUserResponse{Valid: false})
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// Verify password
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		// Password doesn't match - return invalid
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(VerifyUserResponse{Valid: false})
		return
	}

	// Password matches - update user activity and return valid with user data (excluding password)
	s.updateUserActivity(user.ID) // Track user login activity
	user.Password = "" // Clear password field for security
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(VerifyUserResponse{Valid: true, User: &user})
}

// Survey handlers
func (s *Server) createSurvey(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 3 || pathParts[2] != "survey" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	var req CreateSurveyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if req.CurrentWeight <= 0 || req.TargetWeight <= 0 {
		http.Error(w, "currentWeight and targetWeight must be positive", http.StatusBadRequest)
		return
	}
	if req.ActivityLevel < 0 || req.ActivityLevel > 10 {
		http.Error(w, "activityLevel must be between 0 and 10", http.StatusBadRequest)
		return
	}
	if len(req.GoalIDs) == 0 {
		http.Error(w, "at least one goal must be selected", http.StatusBadRequest)
		return
	}

	// Verify user exists
	var userExists bool
	err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM USERS WHERE id = $1)", userID).Scan(&userExists)
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}
	if !userExists {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	// Verify all goals exist
	for _, goalID := range req.GoalIDs {
		var goalExists bool
		err = s.db.QueryRow("SELECT EXISTS(SELECT 1 FROM GOALS WHERE id = $1)", goalID).Scan(&goalExists)
		if err != nil {
			http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
			return
		}
		if !goalExists {
			http.Error(w, fmt.Sprintf("Goal with ID %d not found", goalID), http.StatusBadRequest)
			return
		}
	}

	// Begin transaction
	tx, err := s.db.Begin()
	if err != nil {
		http.Error(w, fmt.Sprintf("Transaction error: %v", err), http.StatusInternalServerError)
		return
	}
	defer tx.Rollback()

	// Insert survey
	var survey Survey
	surveyQuery := `
		INSERT INTO SURVEYS (user_id, current_weight, target_weight, activity_level, created_at)
		VALUES ($1, $2, $3, $4, CURRENT_TIMESTAMP)
		RETURNING id, user_id, current_weight, target_weight, activity_level, created_at`

	err = tx.QueryRow(surveyQuery, userID, req.CurrentWeight, req.TargetWeight, req.ActivityLevel).Scan(
		&survey.ID, &survey.UserID, &survey.CurrentWeight, &survey.TargetWeight, &survey.ActivityLevel, &survey.CreatedAt,
	)
	if err != nil {
		http.Error(w, fmt.Sprintf("Survey creation error: %v", err), http.StatusInternalServerError)
		return
	}

	// Insert survey-goal associations
	for _, goalID := range req.GoalIDs {
		_, err = tx.Exec(
			"INSERT INTO USER_SURVEY_GOALS (survey_id, goal_id, created_at) VALUES ($1, $2, CURRENT_TIMESTAMP)",
			survey.ID, goalID,
		)
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key") {
				http.Error(w, "Duplicate goal selected", http.StatusBadRequest)
				return
			}
			http.Error(w, fmt.Sprintf("Goal association error: %v", err), http.StatusInternalServerError)
			return
		}
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		http.Error(w, fmt.Sprintf("Transaction commit error: %v", err), http.StatusInternalServerError)
		return
	}

	// Fetch the complete survey with goals for response
	surveyWithGoals := s.getSurveyWithGoals(survey.ID)
	if surveyWithGoals == nil {
		http.Error(w, "Failed to fetch created survey", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(surveyWithGoals)
}

func (s *Server) getLatestSurvey(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL path
	pathParts := strings.Split(strings.Trim(r.URL.Path, "/"), "/")
	if len(pathParts) != 4 || pathParts[2] != "survey" || pathParts[3] != "latest" {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(pathParts[1])
	if err != nil {
		http.Error(w, "Invalid user ID format", http.StatusBadRequest)
		return
	}

	// Get the latest survey for the user
	query := `
		SELECT id, user_id, current_weight, target_weight, activity_level, created_at
		FROM SURVEYS
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 1`

	var survey Survey
	err = s.db.QueryRow(query, userID).Scan(
		&survey.ID, &survey.UserID, &survey.CurrentWeight, &survey.TargetWeight, &survey.ActivityLevel, &survey.CreatedAt,
	)

	if err == sql.ErrNoRows {
		http.Error(w, "No surveys found for user", http.StatusNotFound)
		return
	}
	if err != nil {
		http.Error(w, fmt.Sprintf("Database error: %v", err), http.StatusInternalServerError)
		return
	}

	// Get the survey with goals
	surveyWithGoals := s.getSurveyWithGoals(survey.ID)
	if surveyWithGoals == nil {
		http.Error(w, "Failed to fetch survey details", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(surveyWithGoals)
}

func (s *Server) getAllGoals(w http.ResponseWriter, r *http.Request) {
	// Get all goals grouped by category
	goalsByCategory := s.getGoalsByCategory(nil) // nil means no survey ID, so all goals are unselected
	if goalsByCategory == nil {
		http.Error(w, "Failed to fetch goals", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(AllGoalsResponse{Goals: *goalsByCategory})
}

// Helper function to get survey with associated goals grouped by category
func (s *Server) getSurveyWithGoals(surveyID int) *SurveyWithGoals {
	// Get survey details
	var survey Survey
	squery := `
		SELECT id, user_id, current_weight, target_weight, activity_level, created_at
		FROM SURVEYS
		WHERE id = $1`

	err := s.db.Get(&survey, squery, surveyID)
	if err != nil {
		return nil
	}

	// Get goals grouped by category with selection status
	goalsByCategory := s.getGoalsByCategory(&surveyID)
	if goalsByCategory == nil {
		return nil
	}

	return &SurveyWithGoals{
		Survey: survey,
		Goals:  *goalsByCategory,
	}
}

// Helper function to get all goals grouped by category with optional selection status
func (s *Server) getGoalsByCategory(surveyID *int) *GoalsByCategory {
	// Get all goals
	query := `
		SELECT id, category, name, description
		FROM GOALS
		ORDER BY category, name`

	rows, err := s.db.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	// Get selected goal IDs if survey ID is provided
	selectedGoals := make(map[int]bool)
	if surveyID != nil {
		selectedQuery := `
			SELECT goal_id
			FROM USER_SURVEY_GOALS
			WHERE survey_id = $1`

		selectedRows, err := s.db.Query(selectedQuery, *surveyID)
		if err != nil {
			return nil
		}
		defer selectedRows.Close()

		for selectedRows.Next() {
			var goalID int
			err := selectedRows.Scan(&goalID)
			if err != nil {
				return nil
			}
			selectedGoals[goalID] = true
		}
	}

	// Group goals by category
	goalsByCategory := GoalsByCategory{
		Weight:     []GoalWithSelection{},
		Appearance: []GoalWithSelection{},
		Strength:   []GoalWithSelection{},
		Endurance:  []GoalWithSelection{},
	}

	for rows.Next() {
		var goal struct {
			ID          int    `db:"id"`
			Category    string `db:"category"`
			Name        string `db:"name"`
			Description string `db:"description"`
		}

		err := rows.Scan(&goal.ID, &goal.Category, &goal.Name, &goal.Description)
		if err != nil {
			return nil
		}

		goalWithSelection := GoalWithSelection{
			ID:          goal.ID,
			Name:        goal.Name,
			Description: goal.Description,
			Selected:    selectedGoals[goal.ID],
		}

		// Add to appropriate category
		switch strings.ToLower(goal.Category) {
		case "weight":
			goalsByCategory.Weight = append(goalsByCategory.Weight, goalWithSelection)
		case "appearance":
			goalsByCategory.Appearance = append(goalsByCategory.Appearance, goalWithSelection)
		case "strength":
			goalsByCategory.Strength = append(goalsByCategory.Strength, goalWithSelection)
		case "endurance":
			goalsByCategory.Endurance = append(goalsByCategory.Endurance, goalWithSelection)
		}
	}

	return &goalsByCategory
}

// updateUserActivity updates the user's last_active field with 1-hour debouncing
// Only updates if the current last_active is more than 1 hour old or null
func (s *Server) updateUserActivity(userID int) {
	query := `
		UPDATE USERS 
		SET last_active = CURRENT_TIMESTAMP
		WHERE id = $1 
		  AND (last_active IS NULL OR last_active < (CURRENT_TIMESTAMP - INTERVAL '1 hour'))`
	
	// Execute the update - we don't need to check if rows were affected
	// as the debouncing logic is handled by the WHERE clause
	s.db.Exec(query, userID)
}
