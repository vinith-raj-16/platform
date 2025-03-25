package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
	"github.com/vinith-raj-16/user-management/internal/models"
	"github.com/vinith-raj-16/user-management/pkg/config"
)

// StorageService defines all required storage operations
type StorageHandler interface {
	ListFiles(username string, offset string, limit int) ([]models.FileMetadata, string, int, error)
	UploadFile(username, filename string, file io.Reader, size int64) (*models.FileMetadata, error)
	GetStorageUsage(username string) (*models.StorageStatus, error)
}

// LocalStorageService implements StorageService for local filesystem
type LocalStorageHandler struct {
	storagePath string
}

func NewStorageHandler(cfg *config.Config) StorageHandler {
	// Create storage directory if it doesn't exist
	os.MkdirAll(cfg.StoragePath, 0o755)
	return &LocalStorageHandler{
		storagePath: cfg.StoragePath,
	}
}

func generateFileID() string {
	b := make([]byte, 8)
	rand.Read(b)
	return hex.EncodeToString(b)
}

func Init() int64 {
	// Load the .env file
	if err := godotenv.Load("config.env"); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}
	// Get MaxUploadSize from the environment
	maxUploadSizeStr := os.Getenv("DEFAULT_USER_QUOTA")
	if maxUploadSizeStr == "" {
		panic("MAX_UPLOAD_SIZE is not set in the .env file")
	}

	// Convert MaxUploadSize to int64
	maxUploadSize, err := strconv.ParseInt(maxUploadSizeStr, 10, 64)
	if err != nil {
		panic("Invalid MAX_UPLOAD_SIZE value in .env file: " + err.Error())
	}
	return maxUploadSize
}
