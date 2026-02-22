package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morphlinkk/subscriptions/internal/config"
	"github.com/morphlinkk/subscriptions/internal/db"
)

type Server struct {
	store  *db.Store
	config *config.Config
}

func NewServer(conf *config.Config) (*gin.Engine, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	store, err := db.NewStore(ctx, *conf)

	slog.Info("Creating database connection pool")
	if err != nil {
		return nil, fmt.Errorf("failed to create database store: %w", err)
	}

	slog.Debug("Pinging database")
	if err := store.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	gin.SetMode(conf.EnvMode)

	r := gin.Default()
	r.SetTrustedProxies([]string{"localhost"})

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	return r, nil
}
