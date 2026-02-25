package handler

import (
	"time"

	"github.com/google/uuid"
	"github.com/morphlinkk/subscriptions/internal/model"
)

const dateLayout = "01-2006"

type SubscriptionResponse struct {
	ID        int64   `json:"id"`
	Service   string  `json:"service_name"`
	Price     int     `json:"price"`
	UserID    string  `json:"user_id"`
	StartDate string  `json:"start_date"` // MM-YYYY
	EndDate   *string `json:"end_date"`   // MM-YYYY
}

func ToSubscriptionResponse(s model.Subscription) SubscriptionResponse {
	start := s.StartDate.Format(dateLayout)

	var end *string
	if s.EndDate != nil {
		e := s.EndDate.Format(dateLayout)
		end = &e
	}

	return SubscriptionResponse{
		ID:        s.ID,
		Service:   s.Service,
		Price:     s.Price,
		UserID:    s.UserID.String(),
		StartDate: start,
		EndDate:   end,
	}
}

type AddSubscriptionRequest struct {
	Service   string  `json:"service_name" validate:"required"`
	Price     int     `json:"price" validate:"required,gte=0"`
	UserID    string  `json:"user_id" validate:"required,uuid"`
	StartDate string  `json:"start_date" validate:"required"` // MM-YYYY
	EndDate   *string `json:"end_date"`                       // MM-YYYY
}

func (r AddSubscriptionRequest) ToParams() (model.AddSubscriptionParams, error) {
	start, err := time.Parse(dateLayout, r.StartDate)
	if err != nil {
		return model.AddSubscriptionParams{}, err
	}

	var end *time.Time
	if r.EndDate != nil {
		e, err := time.Parse(dateLayout, *r.EndDate)
		if err != nil {
			return model.AddSubscriptionParams{}, err
		}
		end = &e
	}

	uid, err := uuid.Parse(r.UserID)
	if err != nil {
		return model.AddSubscriptionParams{}, err
	}

	return model.AddSubscriptionParams{
		Service:   r.Service,
		Price:     r.Price,
		UserID:    uid,
		StartDate: start,
		EndDate:   end,
	}, nil
}

type UpdateSubscriptionRequest struct {
	Service *string `json:"service_name"`
	Price   *int    `json:"price"`
	UserID  *string `json:"user_id"`
	EndDate *string `json:"end_date"` // MM-YYYY
}

func (r UpdateSubscriptionRequest) ToParams() (model.UpdateSubscriptionParams, error) {
	params := model.UpdateSubscriptionParams{
		Service: r.Service,
		Price:   r.Price,
	}

	if r.UserID != nil {
		uid, err := uuid.Parse(*r.UserID)
		if err != nil {
			return params, err
		}
		params.UserID = &uid
	}

	if r.EndDate != nil {
		t, err := time.Parse(dateLayout, *r.EndDate)
		if err != nil {
			return params, err
		}
		params.EndDate = &t
	}

	return params, nil
}

type ListSubscriptionsRequest struct {
	UserID *string `form:"user_id"`
	Limit  int     `form:"limit"`
	Offset int     `form:"offset"`
}

func (r ListSubscriptionsRequest) ToParams() (model.ListSubscriptionsParams, error) {
	params := model.ListSubscriptionsParams{
		Limit:  r.Limit,
		Offset: r.Offset,
	}

	if r.UserID != nil {
		uid, err := uuid.Parse(*r.UserID)
		if err != nil {
			return params, err
		}
		params.UserID = &uid
	}

	return params, nil
}

type SumOfSubscriptionPricesRequest struct {
	UserID      *string `form:"user_id"`
	Service     *string `form:"service_name"`
	PeriodStart string  `form:"period_start"` // MM-YYYY
	PeriodEnd   string  `form:"period_end"`   // MM-YYYY
}

func (r SumOfSubscriptionPricesRequest) ToParams() (model.SumOfSubscriptionPricesParams, error) {
	params := model.SumOfSubscriptionPricesParams{
		ServiceName: r.Service,
	}

	if r.UserID != nil {
		uid, err := uuid.Parse(*r.UserID)
		if err != nil {
			return params, err
		}
		params.UserID = &uid
	}

	start, err := time.Parse(dateLayout, r.PeriodStart)
	if err != nil {
		return params, err
	}
	params.PeriodStart = start

	end, err := time.Parse(dateLayout, r.PeriodEnd)
	if err != nil {
		return params, err
	}
	params.PeriodEnd = end

	return params, nil
}
