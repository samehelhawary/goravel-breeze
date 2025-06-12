// app/http/middleware/verify_csrf_token.go
package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

// VerifyCSRFToken with excluded paths
func VerifyCSRFToken(excludedPaths ...string) http.Middleware {
	manager := NewCSRFManager()

	return func(ctx http.Context) {
		// Skip for excluded paths
		currentPath := ctx.Request().Path()
		for _, path := range excludedPaths {
			if currentPath == path {
				ctx.Request().Next()
				return
			}
		}

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

// CSRFForAPI for API routes with different error format
func CSRFForAPI() http.Middleware {
	manager := NewCSRFManager()

	return func(ctx http.Context) {
		// Skip CSRF for safe methods
		method := ctx.Request().Method()
		if method == "GET" || method == "HEAD" || method == "OPTIONS" {
			ctx.Request().Next()
			return
		}

		token := manager.getTokenFromRequest(ctx)
		sessionID := ctx.Request().Session().GetID()

		if sessionID == "" || token == "" {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"success": false,
				"message": "CSRF token required",
				"code":    "CSRF_TOKEN_MISSING",
			})
			return
		}

		storedToken := manager.cache.Get(manager.getCacheKey(sessionID))
		if storedToken == nil || storedToken.(string) != token {
			ctx.Request().AbortWithStatusJson(http.StatusUnauthorized, http.Json{
				"success": false,
				"message": "Invalid CSRF token",
				"code":    "CSRF_TOKEN_INVALID",
			})
			return
		}

		ctx.Request().Next()
	}
}
