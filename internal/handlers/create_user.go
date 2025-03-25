package handlers

import (
	"context"
	"fmt"

	_ "modernc.org/sqlite"

	"github.com/vinith-raj-16/user-management/internal/models"
	"golang.org/x/crypto/bcrypt"
)


/* This api used create the user by user creadential */

func (s *UserHandlerImpl) CreateUser(ctx context.Context, user *models.User) error {
	s.mu.Lock()
	defer s.mu.Unlock()

//check is the user exist
	var existingUserCount int
	err := s.db.QueryRowContext(ctx, "SELECT COUNT(*) FROM users WHERE username = ? ", user.Username).Scan(&existingUserCount)
	if err != nil {
		return fmt.Errorf("failed to check username existence: %w", err)
	}
	if existingUserCount > 0 {
		return fmt.Errorf("user already exists")
	}

//hashing the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

//create user in DB
	_, err = s.db.ExecContext(ctx,
		"INSERT INTO users (username, password) VALUES (?, ?)",
		user.Username, string(hashedPassword))
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}
