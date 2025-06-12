package http

import (
	"github.com/goravel/framework/contracts/http"
	sessionMiddleware "github.com/goravel/framework/session/middleware"
	"goravel/app/http/middleware"
)

type Kernel struct {
}

// The application's global HTTP middleware stack.
// These middleware are run during every request to your application.
func (kernel Kernel) Middleware() []http.Middleware {
	return []http.Middleware{
		sessionMiddleware.StartSession(),
		middleware.RememberMe(),
		middleware.GenerateCSRFToken(),
		middleware.InjectCSRFToViews(),
		middleware.AuthFunctions(),
		middleware.FormFieldsValidationCheck(),
	}
}

func (kernel Kernel) RouteMiddleware() map[string]http.Middleware {
	return map[string]http.Middleware{
		//"auth":        middleware.Auth(),
		"csrf":        middleware.CSRF(),
		"csrf.api":    middleware.CSRFForAPI(),
		"csrf.verify": middleware.VerifyCSRFToken(),
		"auth":        middleware.Authenticate(),
		"guest":       middleware.Guest(),
		//"throttle": middleware.Throttle(),
	}
}
