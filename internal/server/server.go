package server

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/morphlinkk/subscriptions/internal/config"
	"github.com/morphlinkk/subscriptions/internal/db"
	"github.com/morphlinkk/subscriptions/internal/repository"
	"github.com/morphlinkk/subscriptions/internal/server/handler"
	"github.com/morphlinkk/subscriptions/internal/server/service"
)

type Server struct {
	store  *db.Store
	config *config.Config
}

type Repositories struct {
	Subscription repository.SubscriptionRepository
}

type Services struct {
	Subscription service.SubscriptionService
}

type Handlers struct {
	Subscription handler.SubscriptionHandler
}

func initRepositories(store *db.Store) *Repositories {
	return &Repositories{
		Subscription: repository.NewSubscriptionRepository(store),
	}
}

func initServices(repositories *Repositories) *Services {
	return &Services{
		Subscription: service.NewSubscriptionService(repositories.Subscription),
	}
}

func initHandlers(services *Services) *Handlers {
	return &Handlers{
		Subscription: handler.NewSubscriptionHandler(services.Subscription),
	}
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

	repositories := initRepositories(store)
	services := initServices(repositories)
	handlers := initHandlers(services)

	r := gin.Default()
	r.SetTrustedProxies([]string{"localhost"})

	subs := r.Group("/subscriptions")
	{
		subs.POST("/", handlers.Subscription.AddSubscription)
		subs.GET("/:id", handlers.Subscription.GetSubscriptionByID)
		subs.PATCH("/:id", handlers.Subscription.UpdateSubscription)
		subs.GET("/", handlers.Subscription.ListSubscriptions)
		subs.GET("/sum", handlers.Subscription.GetSumOfSubscriptionPrices)
	}

	return r, nil
}
