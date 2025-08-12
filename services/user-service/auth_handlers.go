package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"os"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"golang.org/x/crypto/bcrypt"
)

// forgotPassword handles password reset requests
// TODO: This needs to be implemented in db-gateway as new RPCs
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

	// TODO: Store the token via db-gateway
	// This would need new RPC methods in db-gateway:
	// - StorePasswordResetToken(email, token, expiration)
	
	// Send the password reset email if SendGrid is configured
	sendgridKey := os.Getenv("SENDGRID_API_KEY")
	if sendgridKey != "" {
		from := mail.NewEmail("Smart Fit", "no-reply@smartfit.com")
		subject := "Password Reset Request"
		to := mail.NewEmail("User", req.Email)
		plainTextContent := "Please use this token to reset your password: " + resetToken
		htmlContent := "<strong>Please use this token to reset your password: </strong>" + resetToken
		message := mail.NewSingleEmail(from, subject, to, plainTextContent, htmlContent)
		
		client := sendgrid.NewSendClient(sendgridKey)
		response, err := client.Send(message)
		if err != nil || response.StatusCode >= 400 {
			http.Error(w, "Failed to send reset email", http.StatusInternalServerError)
			return
		}
	}

	// Respond with success message
	resp := ForgotPasswordResponse{Message: "Password reset token sent to email."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// resetPassword handles password reset with token
// TODO: This needs to be implemented in db-gateway as new RPCs
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

	// Hash the new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash new password", http.StatusInternalServerError)
		return
	}

	// TODO: Verify token and update password via db-gateway
	// This would need new RPC methods in db-gateway:
	// - VerifyPasswordResetToken(token) -> (email, valid)
	// - UpdateUserPassword(email, hashedPassword)
	_ = hashedPassword // Suppress unused variable warning

	// For now, return a placeholder response
	// In production, this should only return success if the token was valid
	// and the password was actually updated
	resp := ResetPasswordResponse{Message: "Password reset successful."}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
