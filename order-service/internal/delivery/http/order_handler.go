package http

import (
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/cureeeeee/order-service/internal/usecase"
)

type Handler struct {
	uc *usecase.OrderUseCase
}

func NewHandler(uc *usecase.OrderUseCase) *Handler {
	return &Handler{uc: uc}
}

type createOrderRequest struct {
	Amount   float64 `json:"amount"`
	Currency string  `json:"currency"`
}

type updateStatusRequest struct {
	Status string `json:"status"`
}

func (h *Handler) RegisterRoutes(r *gin.Engine) {
	r.POST("/orders", h.createOrder)
	r.PUT("/orders/:id/status", h.updateOrderStatus)
}

func (h *Handler) createOrder(c *gin.Context) {
	var req createOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	order, transactionID, err := h.uc.CreateOrder(c.Request.Context(), req.Amount, req.Currency)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrValidation) {
			status = http.StatusBadRequest
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"order":          order,
		"transaction_id": transactionID,
	})
}

func (h *Handler) updateOrderStatus(c *gin.Context) {
	var req updateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	orderID := c.Param("id")
	order, err := h.uc.UpdateStatus(c.Request.Context(), orderID, req.Status)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, usecase.ErrValidation) {
			status = http.StatusBadRequest
		}
		if errors.Is(err, usecase.ErrNotFound) {
			status = http.StatusNotFound
		}
		c.JSON(status, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"order": order})
}
