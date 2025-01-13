package middleware

import (
	"go-cursor/internal/repository/postgres"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(authRepo *postgres.AuthRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First try to get token from cookie
		tokenCookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")
			if authHeader == "" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
				c.Abort()
				return
			}

			tokenCookie = parts[1]
		}

		token, err := authRepo.GetTokenByHash(tokenCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
			c.Abort()
			return
		}

		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found or expired"})
			c.Abort()
			return
		}

		// Store user ID in cookie and set secure cookie with session token
		SetSecureCookie(c, tokenCookie)

		// After validating the token, set the user ID in a cookie
		c.SetCookie("user_id", token.UserID, int(24*time.Hour.Seconds()), "/", "", true, true)
		c.Set("user_id", token.UserID)
		c.Next()
	}
}
