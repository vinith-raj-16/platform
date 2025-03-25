package handlers

import (
	"fmt"
	"mime"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/vinith-raj-16/user-management/internal/models"
	"github.com/vinith-raj-16/user-management/internal/utils"
)

func (s *LocalStorageHandler) ListFiles(username string, paginationID string, limit int) ([]models.FileMetadata, string, int, error) {
	userDir := filepath.Join(s.storagePath, username)

	var allFiles []models.FileMetadata
	err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() {
			relPath, err := filepath.Rel(s.storagePath, path)
			if err != nil {
				return err
			}

			base := filepath.Base(path)
			id := strings.Split(base, "_")[0]
			filename := strings.Join(strings.Split(base, "_")[1:], "_")

			mimeType := mime.TypeByExtension(filepath.Ext(path))
			if mimeType == "" {
				mimeType = "application/octet-stream"
			}

			allFiles = append(allFiles, models.FileMetadata{
				ID:         id,
				UserID:     username,
				Filename:   filename,
				Path:       relPath,
				Size:       info.Size(),
				SizeHuman:  utils.HumanizeBytes(info.Size()),
				UploadedAt: info.ModTime(),
				MimeType:   mimeType,
			})
		}
		return nil
	})
	if err != nil {
		return nil, "", 0, fmt.Errorf("failed to list files: %w", err)
	}

	// Sort by UploadedAt (newest first)
	sort.Slice(allFiles, func(i, j int) bool {
		return allFiles[i].UploadedAt.After(allFiles[j].UploadedAt)
	})

	// Total count (all files, ignoring pagination)
	total := len(allFiles)

	// Find starting index for pagination
	startIndex := 0
	if paginationID != "" {
		for i, file := range allFiles {
			if file.ID == paginationID {
				startIndex = i + 1
				break
			}
		}
	}

	// Apply pagination
	endIndex := startIndex + limit
	if endIndex > len(allFiles) {
		endIndex = len(allFiles)
	}
	paginatedFiles := allFiles[startIndex:endIndex]

	// Set nextPaginationID if more files exist
	var nextPaginationID string
	if endIndex < len(allFiles) {
		nextPaginationID = allFiles[endIndex].ID
	}

	return paginatedFiles, nextPaginationID, total, nil
}
