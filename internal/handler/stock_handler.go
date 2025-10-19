package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
	"github.com/Skapar/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	service service.Service
	log     logger.Logger
}

func NewStockHandler(service service.Service, log logger.Logger) *StockHandler {
	return &StockHandler{
		service: service,
		log:     log,
	}
}

func (h *StockHandler) CreateStock(c *gin.Context) {
	var input entities.Stock
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	if input.Symbol == "" || input.Name == "" || input.Price <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "symbol, name and positive price are required"})
		return
	}

	input.UpdatedAt = time.Now()

	id, err := h.service.CreateStock(c, &input)
	if err != nil {
		h.log.Errorf("CreateStock error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create stock"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "stock created successfully",
		"id":      id,
	})
}

func (h *StockHandler) GetStockByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock ID"})
		return
	}

	stock, err := h.service.GetStockByID(c, id)
	if err != nil {
		h.log.Errorf("GetStockByID error: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

func (h *StockHandler) GetAllStocks(c *gin.Context) {
	stocks, err := h.service.GetAllStocks(c)
	if err != nil {
		h.log.Errorf("GetAllStocks error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch stocks"})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

func (h *StockHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock ID"})
		return
	}

	existing, err := h.service.GetStockByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "stock not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	if symbol, ok := input["symbol"].(string); ok {
		existing.Symbol = symbol
	}
	if name, ok := input["name"].(string); ok {
		existing.Name = name
	}
	if price, ok := input["price"].(float64); ok {
		if price <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "price must be positive"})
			return
		}
		existing.Price = price
	}

	existing.UpdatedAt = time.Now()

	if err := h.service.UpdateStock(c, existing); err != nil {
		h.log.Errorf("UpdateStock error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock updated successfully"})
}

func (h *StockHandler) DeleteStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid stock ID"})
		return
	}

	if err := h.service.DeleteStock(c, id); err != nil {
		h.log.Errorf("DeleteStock error: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete stock"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "stock deleted successfully"})
}
