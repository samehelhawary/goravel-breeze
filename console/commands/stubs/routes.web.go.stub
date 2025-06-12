package routes

import (
	"github.com/goravel/framework/contracts/http"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/framework/support"
	"goravel/app/http/controllers"
	"goravel/app/http/controllers/auth"
	"goravel/app/http/middleware"
)

func Web() {

	// Create middleware instances
	//generateToken := middleware.NewGenerateCSRFToken()
	//csrfProtection := middleware.NewCSRFFiber()

	// Apply to all routes that need CSRF
	//app.Use(generateToken.Handle())

	facades.Route().Get("/", func(ctx http.Context) http.Response {
		return ctx.Response().View().Make("home", map[string]any{
			"version": support.Version,
		})
	})

	dashboardController := controllers.NewDashboardController()
	facades.Route().Middleware(middleware.Authenticate()).Get("/dashboard", dashboardController.Index)

	registerController := auth.NewRegisterController()
	authController := auth.NewAuthController()

	facades.Route().Middleware(middleware.Guest()).Group(func(router route.Router) {
		router.Get("/register", registerController.Index)
		router.Get("/login", authController.Index)
	})

	facades.Route().Middleware(middleware.CSRF()).Group(func(router route.Router) {
		router.Post("/register", registerController.Store)
		router.Post("/login", authController.Store)
		router.Post("/logout", authController.Logout)
	})

	//facades.Route().Get("/debug-db", func(ctx http.Context) http.Response {
	//	var count int64
	//	err := facades.Orm().Query().Table("users").Count(&count)
	//	if err != nil {
	//		return ctx.Response().String(500, "DATABASE CONNECTION FAILED: "+err.Error())
	//	}
	//	return ctx.Response().String(200, fmt.Sprintf("Successfully connected to DB. Found %d users.", count))
	//})
}
