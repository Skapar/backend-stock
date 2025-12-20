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

// POST /api/history
func (h *HistoryHandler) AddHistory(c *gin.Context) {
	var rec entities.History
	if err := c.ShouldBindJSON(&rec); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	uid, _ := c.Get("userID")
	tokenUserID := uid.(int64)
	role, _ := c.Get("role")
	tokenUserRole := role.(string)

	if tokenUserRole != "ADMIN" {
		rec.UserID = tokenUserID
	}

	id, err := h.cmd.AddHistoryRecord(c, &rec)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to add history record: " + err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"history_id": id})
}

// GET /api/history/user/:user_id
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
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
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
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get history: " + err.Error()})
		return
	}

	c.JSON(http.StatusOK, history)
}
