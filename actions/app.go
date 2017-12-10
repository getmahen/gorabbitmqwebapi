package actions

import (
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/middleware"
	"github.com/gobuffalo/buffalo/middleware/ssl"
	"github.com/gobuffalo/envy"
	"github.com/unrolled/secure"

	"rabbitmqwebapp/messagebroker"
	"rabbitmqwebapp/processor"

	"github.com/gobuffalo/x/sessions"
)

// ENV is used to help switch settings based on where the
// application is being run. Default is "development".
var ENV = envy.Get("GO_ENV", "development")
var app *buffalo.App

// App is where all routes and middleware for buffalo
// should be defined. This is the nerve center of your
// application.
func App() *buffalo.App {
	if app == nil {
		app = buffalo.New(buffalo.Options{
			Env:          ENV,
			SessionStore: sessions.Null{},
			SessionName:  "_rabbitmqwebapp_session",
		})
		// Automatically redirect to SSL
		app.Use(ssl.ForceSSL(secure.Options{
			SSLRedirect:     ENV == "production",
			SSLProxyHeaders: map[string]string{"X-Forwarded-Proto": "https"},
		}))

		// Set the request content type to JSON
		app.Use(middleware.SetContentType("application/json"))

		if ENV == "development" {
			app.Use(middleware.ParameterLogger)
		}

		// Wraps each request in a transaction.
		//  c.Value("tx").(*pop.PopTransaction)
		// Remove to disable this.
		//app.Use(middleware.PopTransaction(models.DB))

		messageBroker := messagebroker.NewMessageBroker("amqp://guest:guest@localhost:5672")
		processor := processor.NewProcessor(70*time.Second, messageBroker)
		app.GET("/", HomeHandler)
		app.GET("/PortRequest", PortRequestHandler(processor))
		app.GET("/PortCallback", PortCallbackHandler(processor))

	}

	return app
}
