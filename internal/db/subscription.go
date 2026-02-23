package db

import (
	"context"

	"github.com/morphlinkk/subscriptions/internal/model"
)

const getSubscriptionByIdQuery = `
	SELECT id,service_name,price,user_id,start_date,end_date
	FROM subscriptions 
	WHERE id = $1
`

func (q *Queries) GetSubscriptionById(ctx context.Context, id int64) (model.Subscription, error) {
	row := q.db.QueryRow(ctx, getSubscriptionByIdQuery, id)
	var s model.Subscription
	err := row.Scan(
		&s.ID,
		&s.Service,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
	)
	return s, err
}

const addSubscriptionQuery = `
	INSERT INTO subscriptions (
			service_name,
			price,
			user_id,
			start_date,
			end_date
	)
	VALUES ($1,$2,$3,$4,$5)
	RETURNING id, service_name, price, user_id, start_date, end_date
`

func (q *Queries) AddSubscription(ctx context.Context, sub model.AddSubscriptionParams) (model.Subscription, error) {
	row := q.db.QueryRow(ctx, addSubscriptionQuery,
		sub.Service,
		sub.Price,
		sub.UserID,
		sub.StartDate,
		sub.EndDate,
	)
	var s model.Subscription
	err := row.Scan(
		&s.ID,
		&s.Service,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
	)
	return s, err
}

const updateSubscriptionQuery = `
	UPDATE subscriptions
	SET service_name = $1,
	    price = $2,
	    start_date = $3,
	    end_date = $4
	WHERE id = $5
	RETURNING id, service_name, price, user_id, start_date, end_date
`

func (q *Queries) UpdateSubscription(ctx context.Context, sub model.Subscription) (model.Subscription, error) {
	row := q.db.QueryRow(ctx, updateSubscriptionQuery,
		sub.Service,
		sub.Price,
		sub.StartDate,
		sub.EndDate,
		sub.ID,
	)

	var s model.Subscription
	err := row.Scan(
		&s.ID,
		&s.Service,
		&s.Price,
		&s.UserID,
		&s.StartDate,
		&s.EndDate,
	)

	return s, err
}

const listSubscriptionsPaginatedQuery = `
	SELECT id, service_name, price, user_id, start_date, end_date
	FROM subscriptions
	WHERE ($1::uuid IS NULL OR user_id = $1)
	ORDER BY start_date DESC
	LIMIT $2 OFFSET $3
`

func (q *Queries) ListSubscriptions(ctx context.Context, params model.ListSubscriptionsParams) ([]model.Subscription, error) {
	rows, err := q.db.Query(ctx, listSubscriptionsPaginatedQuery,
		params.UserID,
		params.Limit,
		params.Offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subs []model.Subscription

	for rows.Next() {
		var s model.Subscription
		if err := rows.Scan(
			&s.ID,
			&s.Service,
			&s.Price,
			&s.UserID,
			&s.StartDate,
			&s.EndDate,
		); err != nil {
			return nil, err
		}
		subs = append(subs, s)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return subs, nil
}

const SumOfSubscriptionPricesQuery = `
	SELECT 
			COALESCE(SUM(price), 0) AS total_price,
			COUNT(*) AS total_subscriptions
	FROM subscriptions
	WHERE ($1::uuid IS NULL OR user_id = $1)
		AND ($2::text IS NULL OR service_name = $2)
		AND start_date >= $3
		AND (end_date IS NULL OR end_date <= $4)
`

func (q *Queries) SumOfSubscriptionPrices(ctx context.Context, params model.SumOfSubscriptionPricesParams) (int64, error) {
	row := q.db.QueryRow(ctx, SumOfSubscriptionPricesQuery,
		params.UserID,
		params.ServiceName,
		params.PeriodStart,
		params.PeriodEnd,
	)

	var totalPrice int64
	err := row.Scan(&totalPrice)
	return totalPrice, err
}
