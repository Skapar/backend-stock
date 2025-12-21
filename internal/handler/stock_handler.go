package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/pkg/logger"
	"github.com/gin-gonic/gin"
)

type StockHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
	log   logger.Logger
}

func NewStockHandler(cmd cqrs.Command, query cqrs.Query, log logger.Logger) *StockHandler {
	return &StockHandler{
		cmd:   cmd,
		query: query,
		log:   log,
	}
}

// CreateStock godoc
// @Summary Create stock (admin)
// @Tags stocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body entities.Stock true "Stock payload"
// @Success 201 {object} IDResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks/ [post]
func (h *StockHandler) CreateStock(c *gin.Context) {
	var input entities.Stock
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if input.Symbol == "" || input.Name == "" || input.Price <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "symbol, name and positive price are required"})
		return
	}

	input.UpdatedAt = time.Now()

	id, err := h.cmd.CreateStock(c, &input)
	if err != nil {
		h.log.Errorf("CreateStock error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to create stock"})
		return
	}

	// переиспользуем RegisterResponse как {message, user_id} не хочется.
	// Лучше вернуть стандартный MessageResponse + id:
	c.JSON(http.StatusCreated, gin.H{
		"message": "stock created successfully",
		"id":      id,
	})
}

// GetStockByID godoc
// @Summary Get stock by ID
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Param id path int true "Stock ID"
// @Success 200 {object} entities.Stock
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks/{id} [get]
func (h *StockHandler) GetStockByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid stock ID"})
		return
	}

	stock, err := h.query.GetStockByID(c, id)
	if err != nil {
		h.log.Errorf("GetStockByID error: %v", err)
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "stock not found"})
		return
	}

	c.JSON(http.StatusOK, stock)
}

// GetAllStocks godoc
// @Summary Get all stocks
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.Stock
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks/ [get]
func (h *StockHandler) GetAllStocks(c *gin.Context) {
	stocks, err := h.query.GetAllStocks(c)
	if err != nil {
		h.log.Errorf("GetAllStocks error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch stocks"})
		return
	}

	c.JSON(http.StatusOK, stocks)
}

// UpdateStock godoc
// @Summary Update stock (admin)
// @Tags stocks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "Stock ID"
// @Param body body object true "Partial update payload (symbol/name/price)"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks/{id} [put]
func (h *StockHandler) UpdateStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid stock ID"})
		return
	}

	existing, err := h.query.GetStockByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "stock not found"})
		return
	}

	var input map[string]interface{}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
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
			c.JSON(http.StatusBadRequest, ErrorResponse{Error: "price must be positive"})
			return
		}
		existing.Price = price
	}

	existing.UpdatedAt = time.Now()

	if err := h.cmd.UpdateStock(c, existing); err != nil {
		h.log.Errorf("UpdateStock error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update stock"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "stock updated successfully"})
}

// DeleteStock godoc
// @Summary Delete stock (admin)
// @Tags stocks
// @Security BearerAuth
// @Produce json
// @Param id path int true "Stock ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /stocks/{id} [delete]
func (h *StockHandler) DeleteStock(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid stock ID"})
		return
	}

	if err := h.cmd.DeleteStock(c, id); err != nil {
		h.log.Errorf("DeleteStock error: %v", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to delete stock"})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "stock deleted successfully"})
}
