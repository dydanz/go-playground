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

func GenerateCSRFToken() string {
	b := make([]byte, csrfTokenLength)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

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
