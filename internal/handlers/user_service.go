package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "modernc.org/sqlite"

	"github.com/vinith-raj-16/user-management/internal/models"
)

// UserHandler defines the interface
type UserHandler interface {
	CreateUser(ctx context.Context, user *models.User) error
	GetUserByUsername(ctx context.Context, username string) (*models.User, error)
}

// userServiceImpl is the private implementation
type UserHandlerImpl struct {
	db *sql.DB
	mu sync.Mutex
}

// NewUserHandler creates a new UserService instance
func NewUserHandler() (*UserHandlerImpl, error) { 

	dataDir := filepath.Join(".", "data")
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	// Open database
	dbPath := filepath.Join(dataDir, "users.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Create users table if not exists
	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS users (
		username TEXT PRIMARY KEY,
		password TEXT NOT NULL
	)`)
	if err != nil {
		return nil, fmt.Errorf("failed to create table: %w", err)
	}

	// Return a pointer to UserHandlerImpl, which satisfies the UserHandler interface
	return &UserHandlerImpl{db: db}, nil
}

// GetUserByUsername implements UserHandler.GetUserByUsername
func (s *UserHandlerImpl) GetUserByUsername(ctx context.Context, username string) (*models.User, error) {
	var user models.User
	err := s.db.QueryRowContext(ctx,
		"SELECT username, password FROM users WHERE username = ?",
		username).Scan(&user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found")
		}
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	return &user, nil
}
