package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/linkeunid/go-api/internal/bootstrap"
	"go.uber.org/zap"
)

func main() {
	// Initialize the application
	app, err := bootstrap.InitializeApp()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}
	defer app.Logger.Sync()

	// Get dependencies
	logger := app.Logger
	cfg := app.Config

	// Setup HTTP server
	server := bootstrap.SetupServer(app, app.AnimalController)

	// Start the server in a goroutine
	go func() {
		bootstrap.LogServerInfo(logger, cfg.Server.Port, cfg.IsDevelopment(), cfg)
		if err := server.ListenAndServe(); err != nil {
			logger.Fatal("Server failed", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// Create a deadline to wait for
	ctx, cancel := context.WithTimeout(context.Background(), cfg.Server.ShutdownTimeout)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline
	if err := server.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	}

	logger.Info("Server exiting")
}
