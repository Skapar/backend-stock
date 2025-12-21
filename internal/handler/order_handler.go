package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
}

func NewOrderHandler(cmd cqrs.Command, query cqrs.Query) *OrderHandler {
	return &OrderHandler{cmd: cmd, query: query}
}

// CreateOrder godoc
// @Summary Create order
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateOrderRequest true "Order payload"
// @Success 201 {object} OrderCreatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/ [post]
func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}
	if req.Type != "BUY" && req.Type != "SELL" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "type must be BUY or SELL"})
		return
	}

	uid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing user_id"})
		return
	}
	tokenUserID := uid.(int64)

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing role"})
		return
	}
	tokenUserRole := roleVal.(string)

	order := entities.Order{
		UserID:    tokenUserID,
		StockID:   req.StockID,
		Quantity:  float64(req.Quantity),
		OrderType: entities.OrderType(req.Type),
		Status:    entities.OrderStatus("PENDING"),
	}

	// admin может создавать на другого пользователя — если вам надо, добавим позже.
	// сейчас оставим как было: не admin -> всегда свой user_id
	if tokenUserRole == "ADMIN" {
		// оставляем UserID = tokenUserID
	}

	stock, err := h.query.GetStockByID(c, order.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch stock: " + err.Error()})
		return
	}

	order.Price = stock.Price * float64(order.Quantity)

	id, err := h.cmd.CreateOrder(c, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create order: " + err.Error()})
		return
	}

	order.ID = id
	if err := h.cmd.ExecuteOrder(c, &order); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to execute order: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, OrderCreatedResponse{
		Message: "order executed successfully",
		OrderID: id,
	})
}

// UpdateOrderStatus godoc
// @Summary Update order status (admin or owner)
// @Tags orders
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Order ID"
// @Param body body UpdateOrderStatusRequest true "New status"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/{id}/status [put]
func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid order ID"})
		return
	}

	order, err := h.query.GetOrderByID(c, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch order: " + err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "order not found"})
		return
	}

	uid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing user_id"})
		return
	}
	tokenUserID := uid.(int64)

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "missing role"})
		return
	}
	tokenUserRole := roleVal.(string)

	if tokenUserRole != "ADMIN" && tokenUserID != order.UserID {
		c.JSON(http.StatusForbidden, ErrorResponse{Error: "access denied"})
		return
	}

	var body UpdateOrderStatusRequest
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}
	if body.Status == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: status is required"})
		return
	}

	if err := h.cmd.UpdateOrderStatus(c, orderID, entities.OrderStatus(body.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "order status updated"})
}

// GetOrdersByUser godoc
// @Summary Get orders (admin can pass user_id, trader gets own)
// @Tags orders
// @Security BearerAuth
// @Produce json
// @Param user_id path int false "User ID (admin only)"
// @Success 200 {array} entities.Order
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /orders/user/{user_id} [get]
// @Router /orders/me [get]
func (h *OrderHandler) GetOrdersByUser(c *gin.Context) {
	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	var userID int64
	if tokenUserRole == "ADMIN" {
		userIDStr := c.Param("user_id")
		if userIDStr != "" {
			var err error
			userID, err = strconv.ParseInt(userIDStr, 10, 64)
			if err != nil {
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
				return
			}
		} else {
			userID = tokenUserID
		}
	} else {
		userID = tokenUserID
	}

	orders, err := h.query.GetOrdersByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
