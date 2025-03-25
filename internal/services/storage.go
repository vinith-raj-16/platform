package services

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/vinith-raj-16/user-management/internal/handlers"
	"go.uber.org/zap"
)

type StorageService struct {
	storageHandler handlers.StorageHandler // Using interface type
	logger         *zap.Logger
}

func NewStorageService(storageHandler handlers.StorageHandler, logger *zap.Logger) *StorageService {
	return &StorageService{
		storageHandler: storageHandler,
		logger:         logger,
	}
}

// get remaining storage
func (h *StorageService) GetFiles(w http.ResponseWriter, r *http.Request) {
	username, ok := r.Context().Value("user").(string)
	if !ok || username == "" {
		h.logger.Error("Username not found in context")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	files, err := h.storageHandler.GetStorageUsage(username) // Change here
	if err != nil {
		h.logger.Error("Failed to list files",
			zap.String("username", username),
			zap.Error(err),
		)
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(files)
}

func (h *StorageService) ListFiles(w http.ResponseWriter, r *http.Request) {
	// Extract username from context
	username, ok := r.Context().Value("user").(string)
	if !ok || username == "" {
		h.logger.Error("username not found in context")
		sendJSONResponse(w, http.StatusUnauthorized, map[string]interface{}{
			"error": "unauthorized",
		})
		return
	}

	// Parse pagination parameters
	paginationID := r.URL.Query().Get("paginationID")
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit <= 0 {
		limit = 10 // Default limit
	}

	// Fetch files and total count
	files, nextPaginationID, total, err := h.storageHandler.ListFiles(username, paginationID, limit)
	if err != nil {
		h.logger.Error("failed to list files", zap.Error(err))
		sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "failed to list files",
		})
		return
	}

	// Build response
	response := map[string]interface{}{
		"files": files,
		"total": total, // Shows the total count of files, not just the paginated results
	}

	if nextPaginationID != "" {
		response["nextPaginationID"] = nextPaginationID
	}

	// Send success response with files data
	sendJSONResponse(w, http.StatusOK, response)
}

// upload file
func (h *StorageService) UploadFile(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	size := handlers.Init()
	if err := r.ParseMultipartForm(size); err != nil {
		h.logger.Error("Failed to parse form", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error": "Invalid form data",
		})
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		h.logger.Error("File parameter missing", zap.Error(err))
		sendJSONResponse(w, http.StatusBadRequest, map[string]interface{}{
			"error": "File upload required",
		})
		return
	}
	defer file.Close()

	// Username from JWT context
	username := r.Context().Value("user").(string)

	// Save file via service
	metadata, err := h.storageHandler.UploadFile(username, header.Filename, file, header.Size) // Change here
	if err != nil {
		h.logger.Error("Upload failed",
			zap.String("username", username),
			zap.Error(err))
		sendJSONResponse(w, http.StatusInternalServerError, map[string]interface{}{
			"error": "Upload failed",
		})
		return
	}

	// Return complete metadata
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metadata)
}
