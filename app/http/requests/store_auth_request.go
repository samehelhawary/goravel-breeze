package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type StoreAuthRequest struct {
	Email    string `form:"email" json:"email"`
	Password string `form:"password" json:"password"`
	Remember string `form:"remember" json:"remember"`
}

func (r *StoreAuthRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreAuthRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "trim",
		"password": "trim",
	}
}

func (r *StoreAuthRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"email":    "required|email",
		"password": "required",
	}
}

func (r *StoreAuthRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreAuthRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreAuthRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
