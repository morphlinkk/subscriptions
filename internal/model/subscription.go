package model

import (
	"time"

	"github.com/google/uuid"
)

type Subscription struct {
	ID        int64      `json:"id"`
	Service   string     `json:"service_name"`
	Price     int        `json:"price"`
	UserID    uuid.UUID  `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type AddSubscriptionParams struct {
	Service   string     `json:"service_name"`
	Price     int        `json:"price"`
	UserID    uuid.UUID  `json:"user_id"`
	StartDate time.Time  `json:"start_date"`
	EndDate   *time.Time `json:"end_date"`
}

type ListSubscriptionsParams struct {
	UserID *uuid.UUID
	Limit  int32
	Offset int32
}

type SumOfSubscriptionPricesParams struct {
	UserID      *uuid.UUID
	ServiceName *string
	PeriodStart time.Time
	PeriodEnd   *time.Time
}
