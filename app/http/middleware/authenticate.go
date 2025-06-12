package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

func Authenticate() http.Middleware {
	return func(ctx http.Context) {
		if ctx.Request().Session().Get("user_id") == nil {
			ctx.Response().Redirect(http.StatusFound, "/login").Render()
			return
		}
		ctx.Request().Next()
	}
}
