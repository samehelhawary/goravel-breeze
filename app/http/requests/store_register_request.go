package requests

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/validation"
)

type StoreRegisterRequest struct {
	Name            string `form:"name" json:"name"`
	Email           string `form:"email" json:"email"`
	Password        string `form:"password" json:"password"`
	PasswordConfirm string `form:"password_confirmation" json:"password_confirmation"`
}

func (r *StoreRegisterRequest) Authorize(ctx http.Context) error {
	return nil
}

func (r *StoreRegisterRequest) Filters(ctx http.Context) map[string]string {
	return map[string]string{
		"name":                  "trim",
		"email":                 "trim",
		"password":              "trim",
		"password_confirmation": "trim",
	}
}

func (r *StoreRegisterRequest) Rules(ctx http.Context) map[string]string {
	return map[string]string{
		"name":                  "required|max_len:255",
		"email":                 "required|email|unique:users,email",
		"password":              "required|confirmed:password",
		"password_confirmation": "required",
	}
}

func (r *StoreRegisterRequest) Messages(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreRegisterRequest) Attributes(ctx http.Context) map[string]string {
	return map[string]string{}
}

func (r *StoreRegisterRequest) PrepareForValidation(ctx http.Context, data validation.Data) error {
	return nil
}
