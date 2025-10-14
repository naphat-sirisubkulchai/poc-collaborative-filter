package main

import (
	"os"
	"os/signal"
	"syscall"

	"poc-collaborative-filter/internal/app"
	"poc-collaborative-filter/internal/di"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

func main() {
	application := fx.New(
		// Provide all dependencies
		di.Module,

		// App server setup
		app.Module,
	)

	// Handle graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-quit
		logger, _ := zap.NewDevelopment()
		logger.Info("Received shutdown signal")
	}()

	// Run the application
	application.Run()
}
