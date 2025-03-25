// pkg/config/config.go
package config

import (
	"os"
	"path/filepath"
	"strconv"
)

type Config struct {
	JWTSecret     string
	StoragePath   string
	ServerPort    int
	MaxUploadSize int64  // Maximum file upload size in bytes
	AllowedTypes  string // Comma-separated list of allowed file types
	DatabasePath  string // For SQLite persistence
	TokenExpiry   int    // Token expiry in hours
}

func Load() (*Config, error) {
	// Ensure storage directory exists
	storagePath := getEnv("STORAGE_PATH", "./storage")
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return nil, err
	}

	// Parse server port
	port, err := strconv.Atoi(getEnv("SERVER_PORT", "8080"))
	if err != nil {
		return nil, err
	}

	// Parse max upload size (default 10MB)
	maxUpload, err := strconv.ParseInt(getEnv("MAX_UPLOAD_SIZE", "10485760"), 10, 64)
	if err != nil {
		return nil, err
	}

	// Parse token expiry (default 24 hours)
	tokenExpiry, err := strconv.Atoi(getEnv("TOKEN_EXPIRY", "24"))
	if err != nil {
		return nil, err
	}

	return &Config{
		JWTSecret:     getEnv("JWT_SECRET", "default-secret-key"),
		StoragePath:   storagePath,
		ServerPort:    port,
		MaxUploadSize: maxUpload,
		AllowedTypes:  getEnv("ALLOWED_TYPES", "pdf,docx,ppt,jpg,png"),
		DatabasePath:  filepath.Join(getEnv("DATA_PATH", "./data"), "users.db"),
		TokenExpiry:   tokenExpiry,
	}, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
