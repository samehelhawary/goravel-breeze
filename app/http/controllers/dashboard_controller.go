package controllers

import (
	"github.com/goravel/framework/contracts/http"
)

type DashboardController struct {
	// Dependent services
}

func NewDashboardController() *DashboardController {
	return &DashboardController{
		// Inject services
	}
}

func (r *DashboardController) Index(ctx http.Context) http.Response {
	return ctx.Response().View().Make("dashboard")
}
