package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/Skapar/backend/config"
	"github.com/Skapar/backend/internal/auth"
	"github.com/Skapar/backend/internal/cqrs"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	cmd   cqrs.Command
	query cqrs.Query
	cfg   *config.Config
}

func NewAuthHandler(cmd cqrs.Command, query cqrs.Query, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		cmd:   cmd,
		query: query,
		cfg:   cfg,
	}
}

// Register godoc
// @Summary Register new user
// @Tags auth
// @Accept json
// @Produce json
// @Param body body RegisterRequest true "Register payload"
// @Success 201 {object} RegisterResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req RegisterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	// Валидация (так как RegisterRequest не содержит binding теги)
	if req.Email == "" || req.Password == "" || req.Role == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: email, password, role are required"})
		return
	}
	if len(req.Password) < 6 {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: password must be at least 6 characters"})
		return
	}
	roleUpper := strings.ToUpper(req.Role)
	if roleUpper != "TRADER" && roleUpper != "ADMIN" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: role must be TRADER or ADMIN"})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to hash password"})
		return
	}

	user := &entities.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entities.Role(roleUpper),
		Balance:  0,
	}

	id, err := h.cmd.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, RegisterResponse{
		Message: "user created successfully",
		UserID:  id,
	})
}

// Login godoc
// @Summary Login
// @Tags auth
// @Accept json
// @Produce json
// @Param body body LoginRequest true "Login payload"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req LoginRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: " + err.Error()})
		return
	}

	if req.Email == "" || req.Password == "" {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: "invalid input: email and password are required"})
		return
	}

	user, err := h.query.GetUserByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	if !auth.CheckPasswordHash(user.Password, req.Password) {
		c.JSON(http.StatusUnauthorized, ErrorResponse{Error: "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(
		h.cfg.JWTSecret,
		h.cfg.JWTTTLMinutes,
		user.ID,
		string(user.Role),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "failed to generate token"})
		return
	}

	// Возвращаем токен + expiresIn (unix time)
	c.JSON(http.StatusOK, gin.H{
		"token":     token,
		"expiresIn": time.Now().Add(time.Duration(h.cfg.JWTTTLMinutes) * time.Minute).Unix(),
	})
}
