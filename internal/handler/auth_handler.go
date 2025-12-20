package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/Skapar/backend/config"
	"github.com/Skapar/backend/internal/auth"
	"github.com/Skapar/backend/internal/models/entities"
	"github.com/Skapar/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	service service.Service
	cfg     *config.Config
}

func NewAuthHandler(s service.Service, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		service: s,
		cfg:     cfg,
	}
}

// POST /register
func (h *AuthHandler) Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,min=6"`
		Role     string `json:"role" binding:"required,oneof=TRADER ADMIN"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input: " + err.Error()})
		return
	}

	hashedPassword, err := auth.HashPassword(req.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	user := &entities.User{
		Email:    req.Email,
		Password: hashedPassword,
		Role:     entities.Role(strings.ToUpper(req.Role)),
		Balance:  0,
	}

	id, err := h.service.CreateUser(c, user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "user created successfully",
		"id":      id,
	})
}

// POST /login
func (h *AuthHandler) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid input"})
		return
	}

	user, err := h.service.GetUserByEmail(c, req.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	passwordOK := auth.CheckPasswordHash(user.Password, req.Password)

	if !passwordOK {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
		return
	}

	token, err := auth.GenerateToken(
		h.cfg.JWTSecret,
		h.cfg.JWTTTLMinutes,
		user.ID,
		string(user.Role),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"expiresIn": time.Now().
			Add(time.Duration(h.cfg.JWTTTLMinutes) * time.Minute).
			Unix(),
	})
}
