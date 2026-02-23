package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/morphlinkk/subscriptions/internal/model"
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
	var req model.AddSubscriptionParams
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	sub, err := h.subscriptionService.AddSubscription(c.Request.Context(), req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, sub)
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

	c.JSON(http.StatusOK, sub)
}

func (h *subscriptionHandler) UpdateSubscription(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid subscription id"})
		return
	}

	var req model.UpdateSubscriptionParams
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request body"})
		return
	}

	sub, err := h.subscriptionService.UpdateSubscription(c.Request.Context(), id, req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, sub)
}

func (h *subscriptionHandler) ListSubscriptions(c *gin.Context) {
	var params model.ListSubscriptionsParams

	if limitStr := c.Query("limit"); limitStr != "" {
		limit, err := strconv.Atoi(limitStr)
		if err == nil {
			params.Limit = limit
		}
	}
	if offsetStr := c.Query("offset"); offsetStr != "" {
		offset, err := strconv.Atoi(offsetStr)
		if err == nil {
			params.Offset = offset
		}
	}

	uid, err := uuid.Parse(c.Query("user_id"))
	if err == nil {
		params.UserID = &uid
	}

	subs, err := h.subscriptionService.ListSubscriptions(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, subs)
}

func (h *subscriptionHandler) GetSumOfSubscriptionPrices(c *gin.Context) {
	var params model.SumOfSubscriptionPricesParams

	if userID := c.Query("user_id"); userID != "" {
		uid, err := uuid.Parse(userID)
		if err == nil {
			params.UserID = &uid
		}
	}
	if serviceName := c.Query("service_name"); serviceName != "" {
		params.ServiceName = &serviceName
	}

	startStr := c.Query("period_start")
	if startStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "period_start is required"})
		return
	}
	start, err := time.Parse(time.RFC3339, startStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_start"})
		return
	}
	params.PeriodStart = start

	endStr := c.Query("period_end")
	if endStr != "" {
		end, err := time.Parse(time.RFC3339, endStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid period_end"})
			return
		}
		params.PeriodEnd = &end
	}

	sum, err := h.subscriptionService.GetSumOfSubscriptionPrices(c.Request.Context(), params)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"total_price": sum,
	})
}
