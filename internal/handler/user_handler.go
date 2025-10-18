package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/auth"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.Service
}

func NewUserHandler(s service.Service) *UserHandler {
	return &UserHandler{service: s}
}

// GET /api/users/:id
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	user, err := h.service.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// PUT /api/users/:id
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	var req struct {
		Email    string  `json:"email" binding:"omitempty,email"`
		Password string  `json:"password" binding:"omitempty,min=6"`
		Role     string  `json:"role" binding:"omitempty,oneof=ADMIN TRADER"`
		Balance  float64 `json:"balance"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.service.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, _ := auth.HashPassword(req.Password)
		user.Password = hashed
	}
	if req.Role != "" {
		user.Role = entities.Role(req.Role)
	}
	user.Balance = req.Balance

	if err := h.service.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user updated successfully"})
}

// DELETE /api/users/:id
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user id"})
		return
	}

	if err := h.service.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user deleted successfully"})
}
