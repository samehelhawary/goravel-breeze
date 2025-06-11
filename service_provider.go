package breeze

import (
	"github.com/goravel/fiber"
	"github.com/goravel/framework/contracts/console"
	"github.com/goravel/framework/contracts/foundation"
	"github.com/samehelhawary/goravel-breeze/console/commands"
)

const Binding = "breeze"

var App foundation.Application

type ServiceProvider struct {
	goravelFiberProvider *fiber.ServiceProvider
}

func (receiver *ServiceProvider) Register(app foundation.Application) {
	App = app

	// fmt.Println(app.BasePath("app/console"))

	app.Bind(Binding, func(app foundation.Application) (any, error) {
		return nil, nil
	})

	receiver.goravelFiberProvider = &fiber.ServiceProvider{}
	receiver.goravelFiberProvider.Register(app)

	app.Publishes("github.com/samehelhawary/goravel-breeze", map[string]string{
		"routes":           app.BasePath("routes"),
		"resources":        app.BasePath("resources"),
		"config/http.go":   app.ConfigPath("http.go"),
		"config/breeze.go": app.ConfigPath("breeze.go"),
	})

	app.Commands([]console.Command{
		&commands.Install{},
	})
}

func (receiver *ServiceProvider) Boot(app foundation.Application) {

	//if facades.Config().GetBool("app.running_in_console") {
	//	routes.Web()
	//}
	// Manually boot sub-providers
	if receiver.goravelFiberProvider != nil {
		receiver.goravelFiberProvider.Boot(app)
	}
}
