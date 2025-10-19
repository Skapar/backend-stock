package middleware

import (
	"net/http"
	"strings"

	"github.com/Skapar/backend/config"
	"github.com/Skapar/backend/internal/auth"
	"github.com/gin-gonic/gin"
)

func AuthMiddleware(cfg *config.Config, allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		if header == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
			return
		}

		tokenStr := strings.TrimPrefix(header, "Bearer ")
		claims, err := auth.ParseToken(cfg.JWTSecret, tokenStr)
		if err != nil {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			return
		}

		if len(allowedRoles) > 0 {
			ok := false
			for _, role := range allowedRoles {
				if role == claims.Role {
					ok = true
					break
				}
			}
			if !ok {
				c.AbortWithStatusJSON(http.StatusForbidden, gin.H{"error": "access denied"})
				return
			}
		}

		c.Set("userID", claims.UserID)
		c.Set("role", claims.Role)
		c.Next()
	}
}
