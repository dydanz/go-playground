package middleware

import (
	"go-playground/pkg/logging"
	"go-playground/server/repository/postgres"
	"go-playground/server/repository/redis"
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
	logger := logging.GetLogger()

	return func(c *gin.Context) {
		// First try to get token from cookie
		tokenCookie, err := c.Cookie(sessionCookieName)
		if err != nil {
			// Fallback to Authorization header
			authHeader := c.GetHeader("Authorization")

			if authHeader == "" {
				logger.Error().
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.RequestURI()).
					Msg("No authorization header")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
				c.Abort()
				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				logger.Error().
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.RequestURI()).
					Str("auth_header", authHeader).
					Msg("Invalid authorization format")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization format"})
				c.Abort()
				return
			}

			tokenCookie = parts[1]
		}

		// Get User-ID from header first, then fallback to cookie
		userIDHeader := c.GetHeader("X-User-Id")
		userIDCookie, _ := c.Cookie(userIdCookieName)

		// Use either header or cookie value for user ID
		userID := userIDHeader
		if userID == "" {
			userID = userIDCookie
		}

		if userID == "" {
			logger.Error().
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Msg("User-ID is required")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID is required"})
			c.Abort()
			return
		}

		// Check session in Redis cache first
		session, err := sessionRepo.GetSession(c.Request.Context(), userID)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Str("user_id", userID).
				Msg("Error getting session from Redis")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "session validation failed"})
			c.Abort()
			return
		}

		if session != nil {
			// Validate User-ID matches session
			if session.UserID != userID {
				logger.Error().
					Str("method", c.Request.Method).
					Str("url", c.Request.URL.RequestURI()).
					Str("session_user_id", session.UserID).
					Str("request_user_id", userID).
					Msg("User-ID mismatch")
				c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID mismatch"})
				c.Abort()
				return
			}

			// Session found in cache, set user context and continue
			SetSecureCookie(c, tokenCookie, session.UserID, "")
			c.SetCookie(userIdCookieName, session.UserID, int(24*time.Hour.Seconds()), "/", "", true, false)

			c.Set("user_id", session.UserID)
			c.Next()
			return
		}

		// If session not found in cache, check database
		token, err := authRepo.GetTokenByHash(c.Request.Context(), tokenCookie)
		if err != nil {
			logger.Error().
				Err(err).
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Str("user_id", userID).
				Msg("Error getting token from database")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token validation failed"})
			c.Abort()
			return
		}

		if token == nil {
			logger.Error().
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Str("user_id", userID).
				Msg("Token not found or expired")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "token not found or expired"})
			c.Abort()
			return
		}

		// Validate User-ID matches token
		if token.UserID != userID {
			logger.Error().
				Str("method", c.Request.Method).
				Str("url", c.Request.URL.RequestURI()).
				Str("token_user_id", token.UserID).
				Str("request_user_id", userID).
				Msg("User-ID mismatch")
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User-ID mismatch"})
			c.Abort()
			return
		}

		// Store user ID in cookie and set secure cookie with session token
		SetSecureCookie(c, tokenCookie, token.UserID, "")
		c.SetCookie(userIdCookieName, token.UserID, int(24*time.Hour.Seconds()), "/", "", true, false)
		c.Set(userIdCookieName, token.UserID)
		c.Next()
	}
}
