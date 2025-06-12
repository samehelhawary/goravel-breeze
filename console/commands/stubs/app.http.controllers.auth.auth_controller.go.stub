package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support/str"
	"goravel/app/http/redirect"
	"goravel/app/http/requests"
	"goravel/app/models"
)

type AuthController struct {
	// Dependent services
}

func NewAuthController() *AuthController {
	return &AuthController{
		// Inject services
	}
}

func (r *AuthController) Index(ctx http.Context) http.Response {
	return ctx.Response().View().Make("auth/login", map[string]interface{}{
		"errors": ctx.Request().Session().Get("errors"),
		"old":    ctx.Request().Session().Get("_old_input"),
	})
}

func (r *AuthController) Store(ctx http.Context) http.Response {
	var storeAuth requests.StoreAuthRequest
	//fmt.Printf("[StoreAuthRequest] storeAuth %+v", storeAuth)
	errors, err := ctx.Request().ValidateRequest(&storeAuth)
	if err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}
	if errors != nil {
		return redirect.New(ctx).Back().WithErrors(errors.All()).WithInput().With("status", "Invalid login details").Go()
	}

	var loggedInUser models.User
	if err = facades.Orm().Query().Where("email", storeAuth.Email).First(&loggedInUser); err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}

	if !facades.Hash().Check(storeAuth.Password, loggedInUser.Password) {
		return redirect.New(ctx).Back().WithInput().With("status", "Invalid login details").Go()
	}

	// --- START OF NEW REMEMBER ME LOGIC ---

	// Always log the user in for the current session first
	ctx.Request().Session().Put("user_id", loggedInUser.ID)

	// Check if the "remember" checkbox was ticked
	if ctx.Request().Input("remember") == "on" {
		// Generate a new secure token
		token := str.Random(60)

		// Save the token to the user's record in the database
		_, err := facades.Orm().Query().Model(&models.User{}).Where("id", loggedInUser.ID).Update("remember_token", token)
		if err != nil {
			facades.Log().Error("failed to save remember token: ", err)
			// Don't block login, just proceed without remember me
		} else {
			// Set the long-lived cookie on the user's browser
			rememberLifetimeSeconds := facades.Config().GetInt("session.remember_lifetime") * 60 // Convert minutes to seconds
			ctx.Response().Cookie(http.Cookie{
				Name:     "remember_me_token",
				Value:    token,
				Path:     "/",
				MaxAge:   rememberLifetimeSeconds,
				Secure:   facades.Config().GetBool("session.secure"),
				HttpOnly: true,
				SameSite: "lax",
			})
		}
	}
	// --- END OF NEW REMEMBER ME LOGIC ---

	// login user functionality
	ctx.Request().Session().Put("user_id", loggedInUser.ID)

	return redirect.New(ctx).To("/dashboard").Go()

}

func (r *AuthController) Logout(ctx http.Context) http.Response {
	userId := ctx.Request().Session().Get("user_id")

	if userId != nil {
		// Clear the remember token from the database
		_, err := facades.Orm().Query().Model(&models.User{}).Where("id", userId).Update("remember_token", nil)
		if err != nil {
			return ctx.Response().View().Make("error", map[string]interface{}{
				"err": err,
			})
		}
	}

	ctx.Request().Session().Forget("user_id")

	// Expire the remember_me cookie immediately
	ctx.Response().Cookie(http.Cookie{
		Name:   "remember_me_token",
		Value:  "",
		Path:   "/",
		MaxAge: -1, // The key to deleting a cookie
	})

	return redirect.New(ctx).To("/login").Go()
}
