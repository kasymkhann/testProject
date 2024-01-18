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

	_ "github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)

	cfg := config.GetConfig()

	logger := logging.GetLogger()

	dbURL := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DB.Host, cfg.DB.User, cfg.DB.Password, cfg.DB.Name, cfg.DB.Port)

	repo, err := repository.NewRepository(dbURL, logger)
	if err != nil {
		logger.Fatalf("Failed to create repository: %v", err)
		return
	}
	service := service.NewService(repo, logger)

	router := gin.Default()

	handlers.RegisterRoutes(router, service)

	port := cfg.App.Port
	addr := fmt.Sprintf(":%d", port)

	logger.Info("Starting server at", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		logger.Fatal("Failed to start server:", err)
		os.Exit(1)
	}

}
