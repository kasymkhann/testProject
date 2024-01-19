package main

import (
	"fmt"
	"net/http"
	"os"
	"testProject/internal/config"
	"testProject/internal/handlers"
	"testProject/pkg/logging"
	"testProject/repository"
	"testProject/service"
	"time"

	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	logger := logging.GetLogger()
	logger.Info("Initializing configuration...")
	cfg := config.GetConfig()
	logger.Info("Configuration initialized successfully.")

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	logger.Info("Creating repository...")
	repo, err := repository.NewRepository(dbURL, logger)
	if err != nil {
		logger.Fatalf("Failed to create repository: %v", err)
		return
	}
	logger.Info("Repository created successfully.")

	logger.Info("Creating service...")
	service := service.NewService(repo, logger)
	logger.Info("Service created successfully.")

	router := gin.Default()
	handlers.RegisterRoutes(router, service)

	port := cfg.App.Port
	addr := fmt.Sprintf(":%d", port)
	logger.Info("Starting server at", addr, "on", time.Now())
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start server:", err)
		os.Exit(1)
	}

}
