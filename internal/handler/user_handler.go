package handler

import (
	"net/http"
	"strconv"

	"github.com/Skapar/backend/internal/auth"
	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
}

func NewUserHandler(cmd cqrs.Command, query cqrs.Query) *UserHandler {
	return &UserHandler{
		cmd:   cmd,
		query: query,
	}
}

// GetUserByID godoc
// @Summary Get user by ID (admin)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} entities.User
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	user, err := h.query.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateUser godoc
// @Summary Update user (admin)
// @Tags users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path int true "User ID"
// @Param body body UpdateUserRequest true "Update payload"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [put]
func (h *UserHandler) UpdateUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	var req UpdateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	user, err := h.query.GetUserByID(c, id)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Password != "" {
		hashed, hashErr := auth.HashPassword(req.Password)
		if hashErr != nil {
			c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
			return
		}
		user.Password = hashed
	}
	if req.Role != "" {
		user.Role = entities.Role(req.Role)
	}
	user.Balance = req.Balance

	if err := h.cmd.UpdateUser(c, user); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "user updated successfully"})
}

// DeleteUser godoc
// @Summary Delete user (admin)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/{id} [delete]
func (h *UserHandler) DeleteUser(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid user id"})
		return
	}

	if err := h.cmd.DeleteUser(c, id); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, MessageResponse{Message: "user deleted successfully"})
}

// GetMe godoc
// @Summary Get my profile
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} GetMeResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID := c.GetInt64("userID")

	user, err := h.query.GetUserByID(c, userID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{Error: "user not found"})
		return
	}

	c.JSON(http.StatusOK, GetMeResponse{
		Email:   user.Email,
		Balance: user.Balance,
	})
}

// GetAllUsers godoc
// @Summary Get all users (admin)
// @Tags users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} entities.User
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users/all [get]
func (h *UserHandler) GetAllUsers(c *gin.Context) {
	users, err := h.query.GetAllUsers(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}
	c.JSON(http.StatusOK, users)
}
