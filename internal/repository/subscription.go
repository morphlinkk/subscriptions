package repository

import (
	"context"

	"github.com/morphlinkk/subscriptions/internal/db"
	"github.com/morphlinkk/subscriptions/internal/model"
)

type SubscriptionRepository interface {
	GetById(ctx context.Context, id int64) (*model.Subscription, error)
	AddSubscription(ctx context.Context, params *model.AddSubscriptionParams) (*model.Subscription, error)
	UpdateSubscription(ctx context.Context, id int64, params *model.UpdateSubscriptionParams) (*model.Subscription, error)
	ListSubscriptions(ctx context.Context, params *model.ListSubscriptionsParams) ([]model.Subscription, error)
	GetSumOfSubscriptionPrices(ctx context.Context, params *model.SumOfSubscriptionPricesParams) (int64, error)
}

type subscriptionRepository struct {
	store *db.Store
}

func NewSubscriptionRepository(store *db.Store) SubscriptionRepository {
	return &subscriptionRepository{
		store,
	}
}

func (r *subscriptionRepository) GetById(ctx context.Context, id int64) (*model.Subscription, error) {
	s, err := r.store.GetSubscriptionById(ctx, id)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) AddSubscription(ctx context.Context, params *model.AddSubscriptionParams) (*model.Subscription, error) {
	s, err := r.store.AddSubscription(ctx, *params)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) UpdateSubscription(ctx context.Context, id int64, params *model.UpdateSubscriptionParams) (*model.Subscription, error) {
	s, err := r.store.UpdateSubscription(ctx, id, *params)
	if err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *subscriptionRepository) ListSubscriptions(ctx context.Context, params *model.ListSubscriptionsParams) ([]model.Subscription, error) {
	s, err := r.store.ListSubscriptions(ctx, *params)
	if err != nil {
		return nil, err
	}
	return s, nil
}

func (r *subscriptionRepository) GetSumOfSubscriptionPrices(ctx context.Context, params *model.SumOfSubscriptionPricesParams) (int64, error) {
	s, err := r.store.SumOfSubscriptionPrices(ctx, *params)
	if err != nil {
		return 0, err
	}
	return s, err
}
