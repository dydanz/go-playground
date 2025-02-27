package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	csrfTokenLength = 32
	csrfCookieName  = "csrf_token"
	csrfHeaderName  = "X-CSRF-Token"
)

// GenerateCSRFToken creates a cryptographically secure random token for CSRF protection.
//
// The function generates a random token of length specified by csrfTokenLength constant
// and encodes it using URL-safe base64 encoding.
//
// Returns:
//   - string: A URL-safe base64 encoded random token
func GenerateCSRFToken() string {
	b := make([]byte, csrfTokenLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// CSRFMiddleware creates a Gin middleware function that provides CSRF protection for routes.
//
// This middleware implements the following security measures:
// 1. Skips CSRF validation for GET requests and authentication endpoints (login/register)
// 2. Validates the presence of CSRF token in both cookie and request header
// 3. Ensures the tokens match to prevent cross-site request forgery attacks
//
// Returns:
//   - gin.HandlerFunc: A middleware function that can be used in Gin routes
//
// The middleware will abort the request with appropriate status codes in case of:
//   - Missing CSRF cookie (403 Forbidden)
//   - Missing CSRF header token (403 Forbidden)
//   - Token mismatch between cookie and header (403 Forbidden)
func CSRFMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Skip CSRF check for login/register
		if c.Request.Method == "GET" ||
			c.Request.URL.Path == "/api/auth/login" ||
			c.Request.URL.Path == "/api/auth/register" {
			c.Next()
			return
		}

		// Get token from cookie
		cookie, err := c.Cookie(csrfCookieName)
		if err != nil {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF cookie not found"})
			c.Abort()
			return
		}

		// Get token from header
		header := c.GetHeader(csrfHeaderName)
		if header == "" {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token missing in headers"})
			c.Abort()
			return
		}

		// Compare cookie and header tokens
		if cookie != header {
			c.JSON(http.StatusForbidden, gin.H{"error": "CSRF token mismatch"})
			c.Abort()
			return
		}

		c.Next()
	}
}
