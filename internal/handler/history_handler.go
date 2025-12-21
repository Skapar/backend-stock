package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/gin-gonic/gin"
)

type HistoryHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
}

func NewHistoryHandler(cmd cqrs.Command, query cqrs.Query) *HistoryHandler {
	return &HistoryHandler{cmd: cmd, query: query}
}

// AddHistory godoc
// @Summary Add history record
// @Tags history
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param body body entities.History true "History payload"
// @Success 201 {object} HistoryCreatedResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /history/ [post]
func (h *HistoryHandler) AddHistory(c *gin.Context) {
	var rec entities.History
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid JSON: " + err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	if tokenUserRole != "ADMIN" {
		rec.UserID = tokenUserID
	} else {
		// admin если не указал user_id — записываем на себя
		if rec.UserID == 0 {
			rec.UserID = tokenUserID
		}
	}

	id, err := h.cmd.AddHistoryRecord(c, &rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to add history record: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, HistoryCreatedResponse{HistoryID: id})
}

// GetHistoryByUser godoc
// @Summary Get history (admin can pass user_id, trader gets own)
// @Tags history
// @Security BearerAuth
// @Produce json
// @Param user_id path int false "User ID (admin only)"
// @Success 200 {array} entities.History
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /history/user/{user_id} [get]
// @Router /history/me [get]
func (h *HistoryHandler) GetHistoryByUser(c *gin.Context) {
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

	history, err := h.query.GetHistoryByUserID(c, userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to get history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
