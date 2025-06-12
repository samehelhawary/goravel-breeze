package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func InjectCSRFToViews() http.Middleware {
	return func(ctx http.Context) {
		// Get CSRF token from session
		token := ctx.Request().Session().Get("csrf_token", "")

		// Share with all views
		facades.View().Share("csrf_token", token)

		// Also add helper functions
		facades.View().Share("csrf_field", func() string {
			return `<input type="hidden" name="_token" value="` + token.(string) + `">`
		})

		facades.View().Share("csrf_meta", func() string {
			return `<meta name="csrf-token" content="` + token.(string) + `">`
		})

		ctx.Request().Next()
	}
}
