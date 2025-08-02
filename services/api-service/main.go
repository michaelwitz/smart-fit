package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

type LoginRequest struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string      `json:"token"`
	User  interface{} `json:"user"`
}

type UserServiceResponse struct {
	Valid bool        `json:"valid"`
	User  interface{} `json:"user"`
}

type Claims struct {
	UserID int    `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func main() {
	// Environment variables
	port := getEnv("SERVICE_PORT", "8080")
	userServiceURL := getEnv("USER_SERVICE_URL", "http://user-service:8080")
	jwtSecret := getEnv("JWT_SECRET", "your-super-secret-jwt-key-change-in-production")

	// Create Gin router
	r := gin.Default()

	// CORS middleware (for development)
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		
		c.Next()
	})

	// Health endpoint
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "healthy",
			"service": "api-service",
		})
	})

	// Authentication endpoints
	auth := r.Group("/auth")
	{
		auth.POST("/login", loginHandler(userServiceURL, jwtSecret))
	}

	// Protected routes (will add middleware later)
	api := r.Group("/api")
	{
		api.GET("/protected", authMiddleware(jwtSecret), func(c *gin.Context) {
			userID, _ := c.Get("user_id")
			email, _ := c.Get("email")
			c.JSON(200, gin.H{
				"message": "This is a protected endpoint",
				"user_id": userID,
				"email":   email,
			})
		})
	}

	log.Printf("API service starting on port %s", port)
	log.Printf("User service URL: %s", userServiceURL)
	r.Run(":" + port)
}

func loginHandler(userServiceURL, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req LoginRequest
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(400, gin.H{"error": "Invalid request payload"})
			return
		}

		// Call user-service to verify credentials
		verifyPayload := map[string]string{
			"email":    req.Email,
			"password": req.Password,
		}
		
		jsonData, err := json.Marshal(verifyPayload)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		resp, err := http.Post(userServiceURL+"/users/verify", "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			log.Printf("Error calling user service: %v", err)
			c.JSON(500, gin.H{"error": "Authentication service unavailable"})
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		var userResp UserServiceResponse
		if err := json.Unmarshal(body, &userResp); err != nil {
			c.JSON(500, gin.H{"error": "Internal server error"})
			return
		}

		if !userResp.Valid {
			c.JSON(401, gin.H{"error": "Invalid email or password"})
			return
		}

		// Extract user info for JWT claims
		userMap := userResp.User.(map[string]interface{})
		userID := int(userMap["id"].(float64))

		// Generate JWT token
		token, err := generateJWT(userID, req.Email, jwtSecret)
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to generate token"})
			return
		}

		c.JSON(200, LoginResponse{
			Token: token,
			User:  userResp.User,
		})
	}
}

func generateJWT(userID int, email, secret string) (string, error) {
	claims := Claims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func authMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(401, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenString := ""
		if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
			tokenString = authHeader[7:]
		} else {
			c.JSON(401, gin.H{"error": "Invalid authorization header format"})
			c.Abort()
			return
		}

		// Parse and validate token
		token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(jwtSecret), nil
		})

		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			c.Set("user_id", claims.UserID)
			c.Set("email", claims.Email)
			c.Next()
		} else {
			c.JSON(401, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
