package middleware

import (
	"fmt"
	"go-cursor/internal/repository/postgres"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authRepo *postgres.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing authorization header"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
			c.Abort()
			return
		}

		tokenHash := parts[1]
		log.Printf("Verifying token: %s", tokenHash)

		token, err := authRepo.GetTokenByHash(tokenHash)
		if err != nil {
			log.Printf("Token verification error: %v", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("token verification error: %v", err)})
			c.Abort()
			return
		}
		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found"})
			c.Abort()
			return
		}

		log.Printf("Token verified for user: %s", token.UserID)
		c.Set("user_id", token.UserID)
		c.Next()
	}
}
