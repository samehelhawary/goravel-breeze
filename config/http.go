package config

import (
	"github.com/gin-gonic/gin/render"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/jet/v2"
	fiberfacades "github.com/goravel/fiber/facades"
	"github.com/goravel/framework/contracts/route"
	"github.com/goravel/framework/facades"
	"github.com/goravel/gin"
	ginfacades "github.com/goravel/gin/facades"
)

func init() {
	config := facades.Config()
	config.Add("http", map[string]any{
		// HTTP Driver
		"default": "fiber",
		// HTTP Drivers
		"drivers": map[string]any{
			"gin": map[string]any{
				// Optional, default is 4096 KB
				"body_limit":   4096,
				"header_limit": 4096,
				"route": func() (route.Route, error) {
					return ginfacades.Route("gin"), nil
				},
				// Optional, default is http/template
				"template": func() (render.HTMLRender, error) {
					return gin.DefaultTemplate()
				},
			},
			"fiber": map[string]any{
				// prefork mode, see https://docs.gofiber.io/api/fiber/#config
				"prefork": false,
				// Optional, default is 4096 KB
				"body_limit":                 4096,
				"header_limit":               4096,
				"proxy_header":               "",
				"enable_trusted_proxy_check": false,
				"trusted_proxies":            []string{},
				"route": func() (route.Route, error) {
					return fiberfacades.Route("fiber"), nil
				},
				// Optional, default is "html/template"
				"template": func() (fiber.Views, error) {
					jetTemplate := jet.New("./resources/views", ".jet")
					jetTemplate.Engine.AddFunc("customFunction", func() string {
						return "new customFunc"
					})
					return jetTemplate, nil
				},
			},
		},
		// HTTP URL
		"url": config.Env("APP_URL", "http://localhost"),
		// HTTP Host
		"host": config.Env("APP_HOST", "127.0.0.1"),
		// HTTP Port
		"port": config.Env("APP_PORT", "3000"),
		// HTTP Timeout, default is 3 seconds
		"request_timeout": 3,
		// HTTPS Configuration
		"tls": map[string]any{
			// HTTPS Host
			"host": config.Env("APP_HOST", "127.0.0.1"),
			// HTTPS Port
			"port": config.Env("APP_PORT", "3000"),
			// SSL Certificate, you can put the certificate in /public folder
			"ssl": map[string]any{
				// ca.pem
				"cert": "",
				// ca.key
				"key": "",
			},
		},
	})
}
