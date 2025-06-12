// app/http/middleware/csrf.go
package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/goravel/framework/contracts/cache"
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"time"
)

type CSRFManager struct {
	cache cache.Cache
}

func NewCSRFManager() *CSRFManager {
	return &CSRFManager{
		cache: facades.Cache(),
	}
}

// CSRF middleware that validates tokens
func CSRF() http.Middleware {
	manager := NewCSRFManager()

	return func(ctx http.Context) {
		// Skip CSRF for GET, HEAD, OPTIONS requests
		method := ctx.Request().Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			ctx.Request().Next()
			return
		}

		// Get token from request
		token := manager.getTokenFromRequest(ctx)

		// Get session ID
		sessionID := ctx.Request().Session().GetID()
		if sessionID == "" {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"error": "Session not found",
			})
			return
		}

		// Validate token
		storedToken := manager.cache.Get(manager.getCacheKey(sessionID))
		if storedToken == nil || storedToken.(string) != token {
			ctx.Request().AbortWithStatusJson(http.StatusForbidden, http.Json{
				"error": "CSRF token mismatch",
			})
			return
		}

		ctx.Request().Next()
	}
}

// GenerateCSRFToken middleware that creates tokens
func GenerateCSRFToken() http.Middleware {
	manager := NewCSRFManager()

	return func(ctx http.Context) {
		// Ensure session exists
		if ctx.Request().Session() == nil || ctx.Request().Session().GetID() == "" {
			ctx.Request().Session().Start()
		}

		// Generate or retrieve token
		token := manager.Token(ctx)

		// Store token in session for views
		ctx.Request().Session().Put("csrf_token", token)

		// Also add to shared data for views
		facades.View().Share("csrf_token", token)

		ctx.Request().Next()
	}
}

// Helper methods for CSRFManager
func (c *CSRFManager) GenerateToken(sessionID string) (string, error) {
	token := c.generateRandomToken()

	// Store in cache for 24 hours
	err := c.cache.Put(c.getCacheKey(sessionID), token, 24*time.Hour)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (c *CSRFManager) getTokenFromRequest(ctx http.Context) string {
	// Check header first
	token := ctx.Request().Header("X-CSRF-TOKEN")
	if token != "" {
		return token
	}

	// Check form data
	token = ctx.Request().Input("_token")
	if token != "" {
		return token
	}

	// Check X-XSRF-TOKEN header (common in SPAs)
	token = ctx.Request().Header("X-XSRF-TOKEN")
	if token != "" {
		return token
	}

	return ""
}

func (c *CSRFManager) generateRandomToken() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

func (c *CSRFManager) getCacheKey(sessionID string) string {
	return "csrf_token:" + sessionID
}

func (c *CSRFManager) Token(ctx http.Context) string {
	sessionID := ctx.Request().Session().GetID()
	if sessionID == "" {
		return ""
	}

	// Check if token already exists in cache
	cachedToken := c.cache.Get(c.getCacheKey(sessionID))
	if cachedToken != nil {
		return cachedToken.(string)
	}

	// Generate new token
	token, err := c.GenerateToken(sessionID)
	if err != nil {
		return ""
	}

	return token
}
