package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	service service.Service
}

func NewPortfolioHandler(s service.Service) *PortfolioHandler {
	return &PortfolioHandler{service: s}
}

// ADMIN GET /api/portfolio/:user_id/:stock_id
// USER GET /api/portfolio/:stock_id
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	stockIDStr := c.Param("stock_id")
	stockID, err := strconv.ParseInt(stockIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock ID"})
		return
	}

	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	var userID int64
	if tokenUserRole == "ADMIN" {
		userIDStr := c.Param("user_id")
		if userIDStr != "" {
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

	p, err := h.service.GetPortfolio(c, userID, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get portfolio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, p)
}

func (h *PortfolioHandler) CreateOrUpdatePortfolio(c *gin.Context) {
	var body struct {
		UserID   int64   `json:"user_id"`
		StockID  int64   `json:"stock_id" binding:"required"`
		Quantity float64 `json:"quantity" binding:"required"`
	}

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	if tokenUserRole != "ADMIN" {
		body.UserID = tokenUserID
	}

	if err := h.service.CreateOrUpdatePortfolio(c, &entities.Portfolio{
		UserID:   body.UserID,
		StockID:  body.StockID,
		Quantity: body.Quantity,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update portfolio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "portfolio updated successfully"})
}
