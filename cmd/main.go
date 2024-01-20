package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.App.Port),
		Handler: router,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		logger.Info("Starting server at", server.Addr, "on", time.Now())
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server:", err)
			os.Exit(1)
		}

	}()

	<-stop
	logger.Info("Received termination signal. Shutting down gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server shutdown error:", err)
	}

	logger.Info("Server gracefully stopped.")

}
