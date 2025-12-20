package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type OrderHandler struct {
	service service.Service
}

func NewOrderHandler(s service.Service) *OrderHandler {
	return &OrderHandler{service: s}
}

func (h *OrderHandler) CreateOrder(c *gin.Context) {
	var order entities.Order
	if err := c.ShouldBindJSON(&order); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	uid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user_id"})
		return
	}
	tokenUserID := uid.(int64)

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
		return
	}
	tokenUserRole := roleVal.(string)

	if tokenUserRole != "ADMIN" {
		order.UserID = tokenUserID
	}

	stock, err := h.service.GetStockByID(c, order.StockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stock: " + err.Error()})
		return
	}

	order.Price = stock.Price * float64(order.Quantity)
	order.Status = entities.OrderStatus("PENDING")

	id, err := h.service.CreateOrder(c, &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create order: " + err.Error()})
		return
	}

	order.ID = id
	if err := h.service.ExecuteOrder(c, &order); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute order: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"order_id": id, "message": "order executed successfully"})
}

func (h *OrderHandler) UpdateOrderStatus(c *gin.Context) {
	idStr := c.Param("id")
	orderID, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid order ID"})
		return
	}

	order, err := h.service.GetOrderByID(c, orderID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch order: " + err.Error()})
		return
	}
	if order == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "order not found"})
		return
	}

	uid, ok := c.Get("userID")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing user_id"})
		return
	}
	tokenUserID := uid.(int64)

	roleVal, ok := c.Get("role")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing role"})
		return
	}
	tokenUserRole := roleVal.(string)

	if tokenUserRole != "ADMIN" && tokenUserID != order.UserID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied"})
		return
	}

	var body struct {
		Status string `json:"status" binding:"required"`
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	if err := h.service.UpdateOrderStatus(c, orderID, entities.OrderStatus(body.Status)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update status: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "order status updated"})
}

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
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
				return
			}
		} else {
			userID = tokenUserID
		}
	} else {
		userID = tokenUserID
	}

	orders, err := h.service.GetOrdersByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get orders: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, orders)
}
