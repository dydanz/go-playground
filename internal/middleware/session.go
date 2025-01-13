package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
)

const (
	sessionCookieName = "session_token"
)

func SetSecureCookie(c *gin.Context, token string) {
	c.SetCookie(
		sessionCookieName,
		token,
		int(24*time.Hour.Seconds()), // 24 hours
		"/",
		"",   // domain
		true, // secure
		true, // httpOnly
	)

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

func ClearSecureCookie(c *gin.Context) {
	c.SetCookie(sessionCookieName, "", -1, "/", "", true, true)
	c.SetCookie(csrfCookieName, "", -1, "/", "", true, false)
}
