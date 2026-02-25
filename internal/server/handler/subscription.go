package handler

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/morphlinkk/subscriptions/internal/server/service"
)

type SubscriptionHandler interface {
	AddSubscription(c *gin.Context)
	GetSubscriptionByID(c *gin.Context)
	UpdateSubscription(c *gin.Context)
	ListSubscriptions(c *gin.Context)
	GetSumOfSubscriptionPrices(c *gin.Context)
}

type subscriptionHandler struct {
	subscriptionService service.SubscriptionService
}

func NewSubscriptionHandler(service service.SubscriptionService) SubscriptionHandler {
	return &subscriptionHandler{
		subscriptionService: service,
	}
}

// AddSubscription godoc
// @Summary Add a new subscription
// @Description Add a subscription for a user
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param subscription body handler.AddSubscriptionRequest true "Subscription info"
// @Success 201 {object} handler.SubscriptionResponse "Created"
// @Failure 400 {object} Response "Invalid request"
// @Failure 500 {object} Response "Internal server error"
// @Router /subscriptions [post]
func (h *subscriptionHandler) AddSubscription(c *gin.Context) {
	var req AddSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Debug("invalid request body", "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid request body")
		return
	}

	params, err := req.ToParams()
	if err != nil {
		slog.Debug("failed to parse AddSubscriptionRequest", "error", err, "body", req)
		JSONErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.subscriptionService.AddSubscription(c.Request.Context(), params)
	if err != nil {
		slog.Error("failed to add subscription", "error", err, "user_id", req.UserID)
		JSONError(c, http.StatusInternalServerError, err)
		return
	}

	slog.Info("subscription created", "id", sub.ID, "user_id", sub.UserID)
	JSONSuccess(c, http.StatusCreated, ToSubscriptionResponse(*sub))
}

// GetSubscriptionByID godoc
// @Summary Get subscription by ID
// @Description Retrieve a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Success 200 {object} Response{data=SubscriptionResponse} "OK"
// @Failure 400 {object} Response "Invalid ID"
// @Failure 404 {object} Response "Not found"
// @Failure 500 {object} Response "Internal server error"
// @Router /subscriptions/{id} [get]
func (h *subscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Debug("invalid subscription id param", "param", idStr, "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	sub, err := h.subscriptionService.GetByID(c.Request.Context(), id)
	if err != nil {
		slog.Error("failed to get subscription by id", "id", id, "error", err)
		JSONError(c, http.StatusInternalServerError, err)
		return
	}
	if sub == nil {
		slog.Debug("subscription not found", "id", id)
		JSONErrorMessage(c, http.StatusNotFound, "subscription not found")
		return
	}

	JSONSuccess(c, http.StatusOK, ToSubscriptionResponse(*sub))
}

// UpdateSubscription godoc
// @Summary Update subscription
// @Description Update a subscription by its ID
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param id path int true "Subscription ID"
// @Param subscription body UpdateSubscriptionRequest true "Subscription update info"
// @Success 200 {object} Response{data=SubscriptionResponse} "Updated"
// @Failure 400 {object} Response "Invalid request or ID"
// @Failure 500 {object} Response "Internal server error"
// @Router /subscriptions/{id} [patch]
func (h *subscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		slog.Debug("invalid subscription id param", "param", idStr, "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid subscription id")
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		slog.Debug("invalid request body for update", "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid request body")
		return
	}

	params, err := req.ToParams()
	if err != nil {
		slog.Debug("failed to parse UpdateSubscriptionRequest", "error", err, "body", req)
		JSONErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	sub, err := h.subscriptionService.UpdateSubscription(c.Request.Context(), id, params)
	if err != nil {
		slog.Error("failed to update subscription", "id", id, "error", err)
		JSONError(c, http.StatusInternalServerError, err)
		return
	}

	slog.Info("subscription updated", "id", sub.ID)
	JSONSuccess(c, http.StatusOK, ToSubscriptionResponse(*sub))
}

// ListSubscriptions godoc
// @Summary List subscriptions
// @Description Get a paginated list of subscriptions, optionally filtered by user_id
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param user_id query string false "Filter by User ID"
// @Param limit query int false "Pagination limit"
// @Param offset query int false "Pagination offset"
// @Success 200 {array} Response{data=SubscriptionResponse} "OK"
// @Failure 400 {object} Response "Invalid query parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /subscriptions [get]
func (h *subscriptionHandler) ListSubscriptions(c *gin.Context) {
	var req ListSubscriptionsRequest
	if err := c.BindQuery(&req); err != nil {
		slog.Debug("invalid query params for ListSubscriptions", "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	params, err := req.ToParams()
	if err != nil {
		slog.Debug("failed to parse ListSubscriptionsRequest", "error", err, "query", req)
		JSONErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	subs, err := h.subscriptionService.ListSubscriptions(c.Request.Context(), params)
	if err != nil {
		slog.Error("failed to list subscriptions", "error", err)
		JSONError(c, http.StatusInternalServerError, err)
		return
	}

	responses := make([]SubscriptionResponse, len(subs))
	for i, s := range subs {
		responses[i] = ToSubscriptionResponse(s)
	}

	JSONSuccess(c, http.StatusOK, responses)
}

// GetSumOfSubscriptionPrices godoc
// @Summary Get sum of subscription prices
// @Description Get total subscription prices for a period, optionally filtered by user or service
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param period_start query string true "Period start MM-YYYY"
// @Param period_end query string true "Period end MM-YYYY"
// @Param user_id query string false "Filter by User ID"
// @Param service_name query string false "Filter by Service name"
// @Success 200 {object} Response{data=map[string]int64} "OK"
// @Failure 400 {object} Response "Invalid query parameters"
// @Failure 500 {object} Response "Internal server error"
// @Router /subscriptions/sum [get]
func (h *subscriptionHandler) GetSumOfSubscriptionPrices(c *gin.Context) {
	var req SumOfSubscriptionPricesRequest
	if err := c.BindQuery(&req); err != nil {
		slog.Debug("invalid query params for SumOfSubscriptionPrices", "error", err)
		JSONErrorMessage(c, http.StatusBadRequest, "invalid query parameters")
		return
	}

	params, err := req.ToParams()
	if err != nil {
		slog.Debug("failed to parse SumOfSubscriptionPricesRequest", "error", err, "query", req)
		JSONErrorMessage(c, http.StatusBadRequest, err.Error())
		return
	}

	sum, err := h.subscriptionService.GetSumOfSubscriptionPrices(c.Request.Context(), params)
	if err != nil {
		slog.Error("failed to calculate sum of subscription prices", "error", err)
		JSONError(c, http.StatusInternalServerError, err)
		return
	}

	JSONSuccess(c, http.StatusOK, map[string]int64{"total_price": sum})
}
