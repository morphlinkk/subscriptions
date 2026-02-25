package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        int64
	Service   string
	Price     int
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   *time.Time
}

type AddSubscriptionParams struct {
	Service   string
	Price     int
	UserID    uuid.UUID
	StartDate time.Time
	EndDate   *time.Time
}

type UpdateSubscriptionParams struct {
	Service *string
	Price   *int
	UserID  *uuid.UUID
	EndDate *time.Time
}

type ListSubscriptionsParams struct {
	UserID *uuid.UUID
	Limit  int
	Offset int
}

type SumOfSubscriptionPricesParams struct {
	UserID      *uuid.UUID
	ServiceName *string
	PeriodStart *time.Time
	PeriodEnd   *time.Time
}
