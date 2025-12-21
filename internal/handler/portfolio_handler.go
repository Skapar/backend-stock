package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/gin-gonic/gin"
)

type PortfolioHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
}

func NewPortfolioHandler(cmd cqrs.Command, query cqrs.Query) *PortfolioHandler {
	return &PortfolioHandler{cmd: cmd, query: query}
}

// GetPortfolio godoc
// @Summary Get portfolio record (admin can specify user_id, trader gets own)
// @Tags portfolio
// @Security BearerAuth
// @Produce json
// @Param user_id path int false "User ID (admin only)"
// @Param stock_id path int true "Stock ID"
// @Success 200 {object} entities.Portfolio
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portfolio/{user_id}/{stock_id} [get]
func (h *PortfolioHandler) GetPortfolio(c *gin.Context) {
	stockIDStr := c.Param("stock_id")
	stockID, err := strconv.ParseInt(stockIDStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid stock ID"})
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
				c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user ID"})
				return
			}
		} else {
			userID = tokenUserID
		}
	} else {
		userID = tokenUserID
	}

	p, err := h.query.GetPortfolio(c, userID, stockID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get portfolio: " + err.Error()})
		return
	}

	if p == nil {
		// оставляем поведение как в твоём коде (200), но делаем модель понятной:
		c.JSON(http.StatusOK, MessageResponse{Message: "no portfolio record found"})
		return
	}

	c.JSON(http.StatusOK, p)
}

// CreateOrUpdatePortfolio godoc
// @Summary Create or update portfolio
// @Tags portfolio
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body CreateOrUpdatePortfolioRequest true "Payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portfolio/ [post]
func (h *PortfolioHandler) CreateOrUpdatePortfolio(c *gin.Context) {
	var body CreateOrUpdatePortfolioRequest

	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	if body.StockID == 0 || body.Quantity == 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: stock_id and quantity are required"})
		return
	}

	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	if tokenUserRole != "ADMIN" {
		body.UserID = tokenUserID
	} else {
		// admin должен явно указать user_id, иначе пусть по умолчанию будет он сам
		if body.UserID == 0 {
			body.UserID = tokenUserID
		}
	}

	if err := h.cmd.CreateOrUpdatePortfolio(c, &entities.Portfolio{
		UserID:   body.UserID,
		StockID:  body.StockID,
		Quantity: body.Quantity,
	}); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to update portfolio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "portfolio updated successfully"})
}

// GetMyPortfolio godoc
// @Summary Get my portfolio
// @Tags portfolio
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.Portfolio
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /portfolio/me [get]
func (h *PortfolioHandler) GetMyPortfolio(c *gin.Context) {
	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)

	portfolios, err := h.query.GetPortfoliosByUserID(c, tokenUserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to fetch portfolio: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, portfolios)
}
