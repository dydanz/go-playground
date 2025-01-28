package middleware

import (
	"go-playground/internal/repository/postgres"
	"go-playground/internal/repository/redis"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// AuthMiddleware creates a Gin middleware function that handles authentication for protected routes.
//
// This middleware performs the following authentication flow:
// 1. Attempts to extract the authentication token from cookies first
// 2. Falls back to Bearer token in Authorization header if cookie is not present
// 3. Validates the User-ID header for additional security
// 4. Validates the token against the database
// 5. Sets user context and session cookies upon successful authentication
//
// Parameters:
//   - authRepo: Pointer to the authentication repository for token validation
//
// Returns:
//   - gin.HandlerFunc: A middleware function that can be used in Gin routes
//
// The middleware will abort the request with appropriate status codes in case of:
//   - Missing or invalid authentication token (401 Unauthorized)
//   - Missing or invalid User-ID header (401 Unauthorized)
//   - Database errors during token validation (401 Unauthorized)
//   - Expired or non-existent tokens (401 Unauthorized)
func AuthMiddleware(authRepo *postgres.AuthRepository, sessionRepo redis.SessionRepository) gin.HandlerFunc {
	return func(c *gin.Context) {
		// First try to get token from cookie
		tokenCookie, err := c.Cookie(sessionCookieName)
		log.Printf("Validating authHeader header: %s\n", tokenCookie)
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

		// Get User-ID from header first, then fallback to cookie
		userIDHeader := c.GetHeader("X-User-Id")
		if userIDHeader == "" {
			// Fallback to cookie
			userIDHeader, err = c.Cookie(userIdCookieName)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID is required"})
				c.Abort()
				return
			}
		}

		// Check session in Redis cache first
		session, err := sessionRepo.GetSession(c.Request.Context(), userIDHeader)
		if err != nil {
			log.Printf("Error getting session from Redis: %v\n", err)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session validation failed"})
			c.Abort()
			return
		}

		if session != nil {
			// Validate User-ID matches session
			if session.UserID != userIDHeader {
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID mismatch"})
				c.Abort()
				return
			}

			// Session found in cache, set user context and continue
			SetSecureCookie(c, tokenCookie)
			c.SetCookie("user_id", session.UserID, int(24*time.Hour.Seconds()), "/", "", true, false)
			c.Set("user_id", session.UserID)
			c.Next()
			log.Printf("Session found in cache %s", tokenCookie)
			return
		}

		log.Println("Continue to find Session data in database")
		// If session not found in cache, check database
		token, err := authRepo.GetTokenByHash(tokenCookie)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token validation failed"})
			log.Printf("Error getting token from database: %v\n", err)
			c.Abort()
			return
		}

		if token == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found or expired"})
			log.Printf("Token not found or expired")
			c.Abort()
			return
		}

		// Validate User-ID matches token
		if token.UserID != userIDHeader {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID mismatch"})
			c.Abort()
			return
		}

		// Store user ID in cookie and set secure cookie with session token
		SetSecureCookie(c, tokenCookie)

		// After validating the token, set the user ID in a cookie
		c.SetCookie("user_id", token.UserID, int(24*time.Hour.Seconds()), "/", "", true, false)
		c.Set("user_id", token.UserID)
		c.Next()
	}
}
