package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

func RememberMe() http.Middleware {
	return func(ctx http.Context) {
		// 1. If user is already logged in via session, do nothing.
		if ctx.Request().Session().Get("user_id") != nil {
			ctx.Request().Next()
			return
		}

		// 2. Check for the remember_me_token cookie.
		rememberToken := ctx.Request().Cookie("remember_me_token")
		if rememberToken == "" {
			// No cookie found, proceed as a guest.
			ctx.Request().Next()
			return
		}

		// 3. Find a user with this token in the database.
		var user models.User
		err := facades.Orm().Query().Where("remember_token", rememberToken).First(&user)
		if err != nil {
			// Token is invalid or user doesn't exist.
			// It's good practice to delete the invalid cookie from the user's browser.
			ctx.Response().Cookie(http.Cookie{Name: "remember_me_token", MaxAge: -1})
			ctx.Request().Next()
			return
		}

		// 4. Log the user in for this request.
		ctx.Request().Session().Put("user_id", user.ID)
		facades.Log().Infof("User %d logged in via Remember Me token.", user.ID)

		ctx.Request().Next()
	}
}
