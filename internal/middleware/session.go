package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	sessionCookieName  = "session_token"
	userIdCookieName   = "user_id"
	userNameCookieName = "user_name"
)

// SetSecureCookie sets secure HTTP-only cookies for session management and CSRF protection.
//
// This function performs two main operations:
// 1. Sets a secure, HTTP-only session cookie with the provided token
// 2. Generates and sets a CSRF token in both cookie and response header
//
// Parameters:
//   - c: The Gin context for the current request
//   - token: The session token to be stored in the cookie
//
// The function sets the following cookies:
//   - Session cookie: HTTP-only, secure cookie containing the session token
//   - CSRF cookie: Secure (not HTTP-only) cookie containing the CSRF token
func SetSecureCookie(c *gin.Context, token string, userID string, userName string) {
	c.SetCookie(
		sessionCookieName,
		token,
		int(24*time.Hour.Seconds()), // 24 hours
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)

	// Set user ID cookie
	c.SetCookie(
		userIdCookieName,
		userID,
		int(24*time.Hour.Seconds()),
		"/",
		"",
		true,
		false, // not httpOnly so JS can read it
	)

	// Set user name cookie
	if userName != "" {
		c.SetCookie(
			userNameCookieName,
			userName,
			int(24*time.Hour.Seconds()),
			"/",
			"",
			true,
			false, // not httpOnly so JS can read it
		)
	}

	// Set CSRF token
	csrfToken := GenerateCSRFToken()
	c.SetCookie(
		csrfCookieName,
		csrfToken,
		int(24*time.Hour.Seconds()),
		"/",
		"",
		true,
		false, // not httpOnly so JS can read it
	)

	// Send CSRF token in response header
	c.Header(csrfHeaderName, csrfToken)
}

// ClearSecureCookie removes all session-related cookies by setting them to expire immediately.
//
// This function is typically called during logout to clear both the session and CSRF cookies.
//
// Parameters:
//   - c: The Gin context for the current request
func ClearSecureCookie(c *gin.Context) {
	c.SetCookie(sessionCookieName, "", -1, "/", "", true, true)
	c.SetCookie(csrfCookieName, "", -1, "/", "", true, false)
}
