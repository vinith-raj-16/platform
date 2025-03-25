package models

import "time"

type User struct {
	ID        string    `json:"id"`
	Username  string    `json:"username"`
	Password  string    `json:"-"`
	CreatedAt time.Time `json:"createdAt"`
}

type RegisterRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FileMetadata struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	Filename   string    `json:"filename"`
	Path       string    `json:"path"`
	Size       int64     `json:"size"`
	SizeHuman  string    `json:"sizeHuman"`
	UploadedAt time.Time `json:"uploadedAt"`
	MimeType   string    `json:"mimeType,omitempty"`
}

type StorageStatus struct {
	AllocatedStorage string `json:"allocatedStorage"`
	UsedStorage    string `json:"usedStorage"`
	Remaining      string `json:"remaining"`
}
