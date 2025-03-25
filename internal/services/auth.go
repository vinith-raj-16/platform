package services

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/vinith-raj-16/user-management/internal/handlers"
	"github.com/vinith-raj-16/user-management/internal/models"
	"github.com/vinith-raj-16/user-management/pkg/config"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	UserHandler handlers.UserHandler
	config      *config.Config
	logger      *zap.Logger
}

func NewAuthService(UserHandler handlers.UserHandler, cfg *config.Config, logger *zap.Logger) *AuthService {
	return &AuthService{
		UserHandler: UserHandler,
		config:      cfg,
		logger:      logger,
	}
}

// Helper function to send a JSON response with a specific HTTP status code
func sendJSONResponse(w http.ResponseWriter, statusCode int, message map[string]interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(message)
}

// create or Register user
func (h *AuthService) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	// Check for mandatory details
	if req.Username == "" || req.Password == "" {
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error": "All fields are required",
		})
		return
	}

	// Create user object
	user := &models.User{
		Username: req.Username,
		Password: req.Password,
	}

	if err := h.UserHandler.CreateUser(r.Context(), user); err != nil {
		h.logger.Error("Failed to create user",
			zap.String("username", req.Username),
			zap.Error(err))

		// Handle specific error if the username already exists
		if err.Error() == "user already exists" {
			sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
				"error": "Username already taken",
			})
		} else {
			// General error if user creation fails
			sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
				"error": "User creation failed",
			})
		}
		return
	}

	// Success response
	sendJSONResponse(w, http.StatusCreated, map[string]interface{}{
		"message": "User registered successfully",
	})
}

// For validating creadetial

func (h *AuthService) Login(w http.ResponseWriter, r *http.Request) {
	var creds models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		h.logger.Error("Invalid request", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid request body",
		})
		return
	}

	// Get user from database
	user, err := h.UserHandler.GetUserByUsername(r.Context(), creds.Username)
	if err != nil {
		h.logger.Error("User not found", zap.String("username", creds.Username))
		sendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid credentials",
		})
		return
	}

	// Compare passwords
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(creds.Password)); err != nil {
		h.logger.Error("Invalid password", zap.String("username", creds.Username))
		sendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"error": "Invalid credentials",
		})
		return
	}

	// Generate JWT token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": user.Username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})

	tokenString, err := token.SignedString([]byte(h.config.JWTSecret))
	if err != nil {
		h.logger.Error("Failed to generate token", zap.Error(err))
		sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Internal server error",
		})
		return
	}

	// Return token
	sendJSONResponse(w, http.StatusOK, map[string]interface{}{
		"token": tokenString,
	})
}
