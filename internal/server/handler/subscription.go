package handler

import (
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
		service,
	}
}

func (h *subscriptionHandler) AddSubscription(c *gin.Context) {
	var req AddSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	params, err := req.ToParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subscriptionService.AddSubscription(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, ToSubscriptionResponse(*sub))
}

func (h *subscriptionHandler) GetSubscriptionByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	sub, err := h.subscriptionService.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "subscription not found"})
		return
	}

	c.JSON(http.StatusOK, ToSubscriptionResponse(*sub))
}

func (h *subscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var req UpdateSubscriptionRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	params, err := req.ToParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sub, err := h.subscriptionService.UpdateSubscription(c.Request.Context(), id, params)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, ToSubscriptionResponse(*sub))
}

func (h *subscriptionHandler) ListSubscriptions(c *gin.Context) {
	var req ListSubscriptionsRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	params, err := req.ToParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	subs, err := h.subscriptionService.ListSubscriptions(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	responses := make([]SubscriptionResponse, len(subs))
	for i, s := range subs {
		responses[i] = ToSubscriptionResponse(s)
	}

	c.JSON(http.StatusOK, responses)
}

func (h *subscriptionHandler) GetSumOfSubscriptionPrices(c *gin.Context) {
	var req SumOfSubscriptionPricesRequest
	if err := c.BindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query parameters"})
		return
	}

	params, err := req.ToParams()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	sum, err := h.subscriptionService.GetSumOfSubscriptionPrices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"total_price": sum})
}
