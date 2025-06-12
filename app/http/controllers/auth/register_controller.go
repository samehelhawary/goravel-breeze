package auth

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"github.com/samehelhawary/goravel-breeze/app/http/redirect"
	"github.com/samehelhawary/goravel-breeze/app/http/requests"
	"github.com/samehelhawary/goravel-breeze/app/models"
)

type RegisterController struct {
	// Dependent services
}

func NewRegisterController() *RegisterController {
	return &RegisterController{
		// Inject services
	}
}

func (r *RegisterController) Index(ctx http.Context) http.Response {

	return ctx.Response().View().Make("auth/register", map[string]interface{}{
		"errors": ctx.Request().Session().Get("errors"),
		"old":    ctx.Request().Session().Get("_old_input"),
	})
}

func (r *RegisterController) Store(ctx http.Context) http.Response {
	var storeRegister requests.StoreRegisterRequest
	errors, err := ctx.Request().ValidateRequest(&storeRegister)
	if err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}
	if errors != nil {
		return redirect.New(ctx).Back().WithErrors(errors.All()).WithInput().Go()
	}

	password, err := facades.Hash().Make(storeRegister.Password)

	if err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}

	user := models.User{
		Name:     storeRegister.Name,
		Email:    storeRegister.Email,
		Password: password,
	}

	if err = facades.Orm().Query().Create(&user); err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}

	var loggedInUser models.User
	if err = facades.Orm().Query().Where("email", user.Email).First(&loggedInUser); err != nil {
		return ctx.Response().View().Make("error", map[string]interface{}{
			"err": err,
		})
	}

	// login user functionality
	ctx.Request().Session().Put("user_id", loggedInUser.ID)

	return redirect.New(ctx).To("/dashboard").Go()
}
