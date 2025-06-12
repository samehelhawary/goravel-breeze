package middleware

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/facades"
)

func FormFieldsValidationCheck() http.Middleware {
	return func(ctx http.Context) {
		// Get CSRF token from session
		errors := ctx.Request().Session().Get("errors")

		// Also add helper functions
		facades.View().Share("hasError", func(field string) bool {
			errs, ok := ctx.Request().Session().Get("errors").(map[string]interface{})
			if ok {
				_, exists := errs[field]
				return exists
			}
			return false
		})

		facades.View().Share("firstError", func(field string) string {
			if errs, ok := errors.(map[string]interface{}); ok {
				if fieldErrors, exists := errs[field]; exists {
					if fieldErrorMap, ok := fieldErrors.(map[string]interface{}); ok {
						for _, message := range fieldErrorMap {
							return message.(string)
						}
					}
				}
			}
			return ""
		})

		ctx.Request().Next()
	}
}
