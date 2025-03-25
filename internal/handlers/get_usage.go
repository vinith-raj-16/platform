package handlers

import (
	"os"
	"path/filepath"

	"github.com/vinith-raj-16/user-management/internal/models"
	"github.com/vinith-raj-16/user-management/internal/utils"
)

func (s *LocalStorageHandler) GetStorageUsage(username string) (*models.StorageStatus, error) {
	userDir := filepath.Join(s.storagePath, username)
	var total int64

	err := filepath.Walk(userDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() {
			total += info.Size()
		}
		return nil
	})

	if os.IsNotExist(err) {
		return nil, nil
	}

	remaining := Init() - total

	return &models.StorageStatus{
		AllocatedStorage: utils.HumanizeBytes(Init()),
		UsedStorage:    utils.HumanizeBytes(total),
		Remaining:      utils.HumanizeBytes(remaining),
	}, err
}
