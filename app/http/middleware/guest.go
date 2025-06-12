package middleware

import (
	"github.com/goravel/framework/contracts/http"
)

func Guest() http.Middleware {
	return func(ctx http.Context) {
		if ctx.Request().Session().Get("user_id") != nil {
			ctx.Response().Redirect(http.StatusFound, "/dashboard").Render()
			return
		}
		ctx.Request().Next()
	}
}
