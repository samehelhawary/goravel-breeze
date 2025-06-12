package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
	"goravel/app/models"
)

type Auth struct {
	id any
}

type UserData struct {
	ID    uint   `gorm:"primaryKey" json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// use in views: {{ auth().Check() }}
func (a *Auth) Check() any {
	return a.id
}

func (a *Auth) GetUser() *UserData {
	var user models.User
	err := facades.Orm().Query().Where("id", a.id).First(&user)
	if err != nil {
		return nil
	}

	return &UserData{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
}

func NewAuth(id any) *Auth {
	return &Auth{id: id}
}

func AuthFunctions() http.Middleware {
	return func(ctx http.Context) {
		// Get CSRF token from session
		userId := ctx.Request().Session().Get("user_id", "")

		facades.View().Share("auth", func() *Auth {
			return NewAuth(userId)
		})

		facades.View().Share("session", func(field string) any {
			return ctx.Request().Session().Get(field, nil)
		})

		ctx.Request().Next()
	}
}
