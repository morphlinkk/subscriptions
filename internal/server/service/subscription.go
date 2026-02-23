package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/morphlinkk/subscriptions/internal/model"
	"github.com/morphlinkk/subscriptions/internal/repository"
)

type SubscriptionService interface {
	GetByID(ctx context.Context, id int64) (*model.Subscription, error)
	AddSubscription(ctx context.Context, params model.AddSubscriptionParams) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id int64, sub model.UpdateSubscriptionParams) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context, params model.ListSubscriptionsParams) ([]model.Subscription, error)
	GetSumOfSubscriptionPrices(ctx context.Context, params model.SumOfSubscriptionPricesParams) (int64, error)
}

type subscriptionService struct {
	repo repository.SubscriptionRepository
}

func NewSubscriptionService(repo repository.SubscriptionRepository) SubscriptionService {
	return &subscriptionService{
		repo,
	}
}

func (s *subscriptionService) GetByID(ctx context.Context, id int64) (*model.Subscription, error) {
	if id <= 0 {
		return nil, errors.New("invalid subscription id")
	}
	return s.repo.GetById(ctx, id)
}

func (s *subscriptionService) AddSubscription(ctx context.Context, params model.AddSubscriptionParams) (*model.Subscription, error) {
	if params.Price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if params.Service == "" {
		return nil, errors.New("service name is required")
	}
	if params.UserID == uuid.Nil {
		return nil, errors.New("user_id is required")
	}
	return s.repo.AddSubscription(ctx, &params)
}

func (s *subscriptionService) UpdateSubscription(ctx context.Context, id int64, params model.UpdateSubscriptionParams) (*model.Subscription, error) {
	if id <= 0 {
		return nil, errors.New("invalid subscription id")
	}
	if params.Price != nil && *params.Price <= 0 {
		return nil, errors.New("price must be positive")
	}
	if params.Service != nil && *params.Service == "" {
		return nil, errors.New("service name is required")
	}
	return s.repo.UpdateSubscription(ctx, id, &params)
}

func (s *subscriptionService) ListSubscriptions(ctx context.Context, params model.ListSubscriptionsParams) ([]model.Subscription, error) {
	if params.Limit <= 0 {
		params.Limit = 20
	}
	if params.Offset < 0 {
		params.Offset = 0
	}
	return s.repo.ListSubscriptions(ctx, &params)
}

func (s *subscriptionService) GetSumOfSubscriptionPrices(ctx context.Context, params model.SumOfSubscriptionPricesParams) (int64, error) {
	if params.PeriodStart.IsZero() {
		return 0, errors.New("period start is required")
	}
	return s.repo.GetSumOfSubscriptionPrices(ctx, &params)
}
