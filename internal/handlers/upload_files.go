package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/vinith-raj-16/user-management/internal/models"
	"github.com/vinith-raj-16/user-management/internal/utils"
)

func (s *LocalStorageHandler) UploadFile(username, filename string, file io.Reader, size int64) (*models.FileMetadata, error) {
	userDir := filepath.Join(s.storagePath, username)
	if err := os.MkdirAll(userDir, 0o755); err != nil {
		return nil, err
	}

	// Check if file exceeds max size
	if size > Init() {
		return nil, fmt.Errorf("file size exceeds the limit")
	}

	fileID := generateFileID()
	storagePath := filepath.Join(username, fileID+"_"+filename)
	fullPath := filepath.Join(s.storagePath, storagePath)

	// Save file
	dst, err := os.Create(fullPath)
	if err != nil {
		return nil, err
	}
	defer dst.Close()

	if _, err := io.Copy(dst, file); err != nil {
		return nil, err
	}

	// Get file info
	info, err := os.Stat(fullPath)
	if err != nil {
		return nil, err
	}
	// Check if file exceeds max size
	if size > info.Size() {
		return nil, fmt.Errorf("file size exceeds the limit")
	}

	// Detect MIME type
	if seeker, ok := file.(io.Seeker); ok {
		seeker.Seek(0, io.SeekStart)
	} else {
		return nil, fmt.Errorf("file does not support seeking")
	}
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}
	mimeType := http.DetectContentType(buffer)

	return &models.FileMetadata{
		ID:         fileID,
		UserID:     username,
		Filename:   filename,
		Path:       storagePath,
		Size:       info.Size(),
		SizeHuman:  utils.HumanizeBytes(info.Size()),
		UploadedAt: info.ModTime(),
		MimeType:   mimeType,
	}, nil
}
