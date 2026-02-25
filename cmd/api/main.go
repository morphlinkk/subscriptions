// @title Subscriptions API
// @version 1.0
// @description API for managing user subscriptions
// @host localhost:3000
// @BasePath /

package main

import (
	"log/slog"

	"github.com/morphlinkk/subscriptions/internal/config"
	"github.com/morphlinkk/subscriptions/internal/logger"
	"github.com/morphlinkk/subscriptions/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		logger.Fatal("Failed to load configuration", "error", err)
	}

	logger.SetLogLevel(cfg)

	srv, err := server.NewServer(cfg)
	if err != nil {
		logger.Fatal("Failed to create server", "error", err)
	}

	slog.Info("Starting server on", "port", cfg.ServerPort)
	if err := srv.Run(":" + cfg.ServerPort); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
