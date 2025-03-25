package main

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/vinith-raj-16/user-management/internal/handlers"
	"github.com/vinith-raj-16/user-management/internal/middleware"
	"github.com/vinith-raj-16/user-management/internal/services"
	"github.com/vinith-raj-16/user-management/pkg/config"
	"go.uber.org/zap"
)

func main() {
	// Initialize logger first
	logger, err := zap.NewProduction()
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}
	defer logger.Sync()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", zap.Error(err))
	}

	// Initialize user service
	userService, err := handlers.NewUserHandler()
	if err != nil {
		logger.Fatal("Failed to initialize user service", zap.Error(err))
	}

	// Initialize auth service
	authService := services.NewAuthService(userService, cfg, logger)

	// Initialize storage handler and storage service
	storageHandler := handlers.NewStorageHandler(cfg)                  
	storageService := services.NewStorageService(storageHandler, logger) 

	// Setup router
	r := mux.NewRouter()

	// Public routes
	r.HandleFunc("/register", authService.Register).Methods("POST")
	r.HandleFunc("/login", authService.Login).Methods("POST")

	// Protected routes
	protected := r.PathPrefix("/api").Subrouter()
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret, logger))
	protected.HandleFunc("/upload", storageService.UploadFile).Methods("POST")
	protected.HandleFunc("/storage/remaining", storageService.GetFiles).Methods("GET")
	protected.HandleFunc("/files", storageService.ListFiles).Methods("GET")

	// Start server
	serverAddr := ":" + strconv.Itoa(cfg.ServerPort)
	logger.Info("Starting server on " + serverAddr)
	if err := http.ListenAndServe(serverAddr, r); err != nil {
		logger.Fatal("Server failed to start", zap.Error(err))
	}
}
